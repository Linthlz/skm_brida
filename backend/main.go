package main

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/controllers"

	// "github.com/itsXanadium/singa-riset-backend/database/seeders"
	"github.com/itsXanadium/singa-riset-backend/infrastructure/code"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/routes"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

func errorHandler(c *fiber.Ctx, err error) error {
	var fe *fiber.Error
	if errors.As(err, &fe) {
		switch fe.Code {
		case fiber.StatusNotFound:
			return utils.NotFound(c, "Resource not found", fe.Message)
		case fiber.StatusRequestEntityTooLarge:
			return c.Status(fe.Code).JSON(fiber.Map{"success": false, "message": "Request too large", "error": fe.Message})
		case fiber.StatusMethodNotAllowed, fiber.StatusBadRequest:
			return c.Status(fe.Code).JSON(fiber.Map{"success": false, "message": "Bad request", "error": fe.Message})
		}
	}
	return utils.InternalServerError(c, "Something went wrong", err.Error())
}

func main() {
	config.LoadEnv()
	config.DBConnect()
	smtpConfig := config.GetConfig()
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
		BodyLimit:    12 * 1024 * 1024,
	})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.CORSOrigins,
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Content-Disposition",
		MaxAge:           600,
	}))
	//Repo and Service Piles
	userRepo := repositories.NewUserRepository()
	smtpService := utils.NewSMTPService(*smtpConfig)
	codeGen := code.NewCodeGenerator()
	userService := services.NewUserService(userRepo, smtpService, codeGen, models.UpdateUser{})
	userController := controllers.NewUserController(userService)
	researchRepo := repositories.NewResearchRepository()
	notifRepo := repositories.NewNotificationRepository()
	notifService := services.NewNotificationService(notifRepo, config.DB)
	researchService := services.NewResearchService(userRepo, researchRepo, notifRepo, notifService)
	researchGroupRepo := repositories.NewResearchGroupRepository()
	researchGroupService := services.NewResearchGroupService(researchGroupRepo, notifService, config.DB, userRepo)
	submissionCommentRepo := repositories.NewSubmissionCommentRepository()
	SubmissionCommentService := services.NewSubmissionCommentService(submissionCommentRepo, researchRepo, userRepo, notifService, config.DB, researchGroupService)
	GroupDocumentationRepo := repositories.NewGroupDocumentationRepository()
	GroupDocumentationService := services.NewGroupDocumentationService(GroupDocumentationRepo, researchGroupRepo, config.DB)
	catRepo := repositories.NewNewsCategoryRepository()
	catService := services.NewNewsCategoryService(catRepo, config.DB)
	newsRepo := repositories.NewNewsRepository()
	newsService := services.NewNewsService(newsRepo, catRepo, config.DB)
	deliverableRepo := repositories.NewResearchDeliverableRepository()
	deliverableService := services.NewResearchDeliverableService(deliverableRepo, researchGroupRepo, config.DB)
	opdRepo := repositories.NewOpdRepository()
	monevRepo := repositories.NewMonevRepository()
	monevService := services.NewMonevService(monevRepo, userRepo, opdRepo, config.DB)
	//Controllers Piles
	monevController := controllers.NewMonevController(monevService)
	researchController := controllers.NewResearchController(researchService, researchGroupService)
	notifController := controllers.NewNotificationController(notifService)
	ReviewResearchController := controllers.NewResearchReviewController(SubmissionCommentService, researchService)
	groupDocumentationController := controllers.NewGroupDocumentationController(GroupDocumentationService)
	catController := controllers.NewNewsCategoryController(catService)
	newsController := controllers.NewNewsController(newsService)
	deliverableController := controllers.NewResearchDeliverableController(deliverableService)
	skmRepo := repositories.NewSkmRepository()
	skmService := services.NewSkmService(skmRepo)
	skmController := controllers.NewSkmController(skmService)

	app.Get("/", func(c *fiber.Ctx) error {
		return utils.Success(c, "API is running", fiber.Map{"API Status": "Running"})
	})
	app.Get("/api/health", func(c *fiber.Ctx) error {
		sqlDB, err := config.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"success": false, "message": "Database unavailable", "error": "database unavailable"})
		}
		return utils.Success(c, "OK", fiber.Map{"database": "up"})
	})

	// Riset yang disetujui sebelum fitur kelompok tersedia belum memiliki kelompok.
	if n, err := researchGroupService.BackfillMissingGroups(); err != nil {
		log.Printf("[groups] backfill gagal: %v", err)
	} else if n > 0 {
		log.Printf("[groups] %d kelompok riset dibuat untuk riset yang sudah disetujui", n)
	}

	routes.Setup(app, userController, researchController, notifController, ReviewResearchController, groupDocumentationController, catController, newsController, deliverableController, monevController)
	routes.SetupSkm(app, skmController, userService)

	app.Use(func(c *fiber.Ctx) error {
		return utils.NotFound(c, "Endpoint not found", c.Method()+" "+c.Path())
	})

	port := config.AppConfig.AppPort
	log.Println("Server is running on port: ", port)
	log.Fatal(app.Listen(":" + port))
}
