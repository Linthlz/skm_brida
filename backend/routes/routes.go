package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/itsXanadium/singa-riset-backend/api/handler"
	"github.com/itsXanadium/singa-riset-backend/application"
	"github.com/itsXanadium/singa-riset-backend/config"
	"github.com/itsXanadium/singa-riset-backend/controllers"
	"github.com/itsXanadium/singa-riset-backend/infrastructure/code"
	"github.com/itsXanadium/singa-riset-backend/middleware"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
	// "github.com/itsXanadium/singa-riset-backend/email"
)

func Setup(app *fiber.App, uc *controllers.UserController, rc *controllers.ResearchController, nc *controllers.NotificationController, rvr *controllers.ResearchReviewController, gdc *controllers.GroupDocumentationController, nct *controllers.NewsCategoryController, nw *controllers.NewsController, rdv *controllers.ResearchDeliverableController, mc *controllers.MonevController) {
	cfg := config.GetConfig()
	smtpConfig := config.GetConfig()
	verificationRepo := repositories.NewVerificationRepository()
	emailService := utils.NewSMTPService(*cfg)
	codeGen := code.NewCodeGenerator()
	userRepo := repositories.NewUserRepository()
	// userService := services.NewUserService(userRepo)
	smtpService := utils.NewSMTPService(*smtpConfig)
	userService := services.NewUserService(userRepo, smtpService, codeGen, models.UpdateUser{})
	groupRepo := repositories.NewResearchGroupRepository()
	verificationService := application.NewVerificationService(verificationRepo, emailService, codeGen, cfg)
	verificationHandler := handler.NewVerificationHandler(verificationService, userService)
	// repos := repositories.NewVerificationRepository()
	// repo := domain.EmailService
	//============ Public Endpoints ============
	//V1 Requires Auth
	//V2 Public
	app.Static("/storage", "./public/storage")
	cat := controllers.NewResearchCatalogController()
	app.Get("/api/v2/researches", cat.PublicList)
	app.Get("/api/v2/researches/facets", cat.PublicFacets)
	app.Get("/api/v2/researches/:slug", cat.PublicDetail)
	app.Get("/api/v2/documentations", cat.PublicDocumentations)
	content := controllers.NewContentController()
	app.Get("/api/v2/funding", content.ListFunding)
	app.Get("/api/v2/settings", content.GetSettings)
	app.Get("/api/v2/stats", content.Stats)
	app.Get("/api/v2/news-categories", nct.List)
	app.Get("/api/v2/news", nw.List)
	app.Get("/api/v2/news/:slug", nw.GetBySlug)
	app.Get("/api/v2/group/list", rc.ListPublic)
	app.Get("/api/v2/group/:slug", rc.GetPublicBySlug)
	app.Post("/api/v1/auth/register", uc.Register)
	app.Post("/api/v1/auth/login", uc.Login)
	app.Post("/api/v1/auth/verify/send/:email", limiter.New(limiter.Config{Max: 5, Expiration: 1 * time.Minute}), verificationHandler.SendVerification)
	app.Post("/api/v1/auth/verify", verificationHandler.CheckCode)
	//============ Authorization req ============
	api := app.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey: []byte(config.AppConfig.JWTSecret),
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.UnauthorizedAccess(c, "Unauthorized", "missing, invalid or expired token")
		},
	}))

	//============ Profile Endpoints ============
	userGroup := api.Group("/users")
	userGroup.Patch("/:id", middleware.Verified(userService), middleware.UserProfile(userService), uc.UpdateUser)
	userGroup.Get("/:id", middleware.Verified(userService), middleware.UserProfile(userService), uc.FetchProfileData)

	//============ My Research Group Endpoints============
	//============ Research Endpoints============
	research := api.Group("/research")
	research.Post("/submit", middleware.Verified(userService), rc.ResearchSubmission)
	research.Get("/submission", middleware.Verified(userService), rc.MyResearchSubmission)
	research.Patch("/revision/:id", middleware.Verified(userService), rc.Revision)
	research.Get("/:slug", middleware.Verified(userService), cat.Detail)

	//============ Research Group Endpoints============
	grp := controllers.NewGroupController()
	member := middleware.GroupMemberOnly(groupRepo)
	researchGroup := api.Group("/research/group")
	researchGroup.Post("/join", middleware.Verified(userService), rc.JoinGroup)
	researchGroup.Get("/list", rc.List)
	researchGroup.Get("/:slug", rc.GetBySlug)
	researchGroup.Patch("/:id/research", middleware.Verified(userService), member, grp.UpdateResearch)
	researchGroup.Delete("/:id/members/:memberID", middleware.Verified(userService), member, grp.RemoveMember)
	researchGroup.Post("/:id/leave", middleware.Verified(userService), member, grp.Leave)
	researchGroup.Post("/:id/join-code", middleware.Verified(userService), member, grp.RegenerateCode)
	researchGroup.Patch("/:id/supervisor", middleware.Verified(userService), member, grp.SetSupervisor)
	researchGroup.Get("/:id/documentations", middleware.Verified(userService), member, gdc.GetGroupDocumentations)
	researchGroup.Post("/:id/documentations", middleware.Verified(userService), member, gdc.CreateDocumentation)
	researchGroup.Get("/:id/documentations/:docID", middleware.Verified(userService), member, gdc.GetDocumentation)
	researchGroup.Patch("/:id/documentations/:docID", middleware.Verified(userService), member, gdc.UpdateDocumentation)
	researchGroup.Delete("/:id/documentations/:docID", middleware.Verified(userService), member, gdc.DeleteDocumentation)

	//============ Research Group Deliverables============
	researchGroup.Get("/:id/deliverables", middleware.Verified(userService), member, rdv.List)
	researchGroup.Post("/:id/deliverables", middleware.Verified(userService), member, rdv.Create)
	researchGroup.Patch("/:id/deliverables/:dID", middleware.Verified(userService), member, rdv.Update)
	researchGroup.Post("/:id/deliverables/:dID/done", middleware.Verified(userService), member, rdv.MarkDone)
	researchGroup.Post("/:id/deliverables/:dID/pending", middleware.Verified(userService), member, rdv.MarkPending)
	researchGroup.Delete("/:id/deliverables", middleware.Verified(userService), member, rdv.Delete)
	// Alias sesuai koleksi Postman ("My Research Groups").
	userGroup.Get("/research/group", middleware.Verified(userService), grp.Mine)
	userGroup.Get("/research/group/:id", middleware.Verified(userService), member, grp.Detail)

	//============ Admin Endpoints ===============
	adminGroup := api.Group("/admin")
	adminGroup.Post("/create", middleware.Verified(userService), middleware.AdminOnly(userService), uc.CreateUser)
	adminGroup.Post("/user", middleware.Verified(userService), middleware.AdminOnly(userService), uc.CreateUser)
	adminGroup.Get("/user", middleware.Verified(userService), middleware.AdminOnly(userService), uc.Users)
	adminGroup.Patch("/user/:id", middleware.Verified(userService), middleware.AdminOnly(userService), uc.UpdateUserByAdmin)
	adminGroup.Delete("/user/:id", middleware.Verified(userService), middleware.AdminOnly(userService), uc.DeleteUserByAdmin)
	adminGroup.Put("/settings", middleware.Verified(userService), middleware.AdminOnly(userService), content.SaveSettings)
	// Kelompok riset untuk panel admin (admin melihat seluruh kelompok).
	adminGroup.Get("/research/groups", middleware.Verified(userService), middleware.AdminOnly(userService), grp.Mine)
	adminGroup.Post("/research/groups/sync", middleware.Verified(userService), middleware.AdminOnly(userService), rc.SyncGroups)
	adminGroup.Get("/research/groups/:slug", middleware.Verified(userService), middleware.AdminOnly(userService), grp.BySlug)
	api.Post("/funding", middleware.Verified(userService), middleware.AdminOnly(userService), content.CreateFunding)
	api.Patch("/funding/:id", middleware.Verified(userService), middleware.AdminOnly(userService), content.UpdateFunding)
	api.Delete("/funding/:id", middleware.Verified(userService), middleware.AdminOnly(userService), content.DeleteFunding)

	//============ Reviewer Endpoints ===============
	reviewerGroup := api.Group("/reviewer")
	reviewerGroup.Get("/researches", middleware.ReviewerRole(userService), cat.StaffList)
	reviewerGroup.Get("/notifications", middleware.ReviewerRole(userService), nc.FetchNotifications)
	reviewerGroup.Post("/assign/:id", middleware.ReviewerRole(userService), rvr.AssignReviewer)
	reviewerGroup.Post("/approve/:id", middleware.ReviewerRole(userService), rvr.ApproveSubmission)
	reviewerGroup.Post("/comment/:id", middleware.ReviewerRole(userService), rvr.AddComment)

	//============ News Endpoints ===============
	api.Post("/news-categories", middleware.Verified(userService), middleware.MediaEmployee(userService), nct.Create)
	api.Patch("/news-categories/:id", middleware.Verified(userService), middleware.MediaEmployee(userService), nct.Update)
	api.Delete("/news-categories/:id", middleware.Verified(userService), middleware.MediaEmployee(userService), nct.Delete)

	api.Get("/news", middleware.Verified(userService), middleware.MediaEmployee(userService), nw.ListAll)
	api.Get("/news/:id", middleware.Verified(userService), middleware.MediaEmployee(userService), nw.GetByID)
	api.Post("/news", middleware.Verified(userService), middleware.MediaEmployee(userService), nw.CreateNews)
	api.Patch("/news/:id", middleware.Verified(userService), middleware.MediaEmployee(userService), nw.Update)
	api.Delete("/news/:id", middleware.Verified(userService), middleware.MediaEmployee(userService), nw.Delete)

	monevGroup := api.Group("/monev", middleware.Verified(userService))

	// Admin only — manage batch
	monevGroup.Post("/batches", middleware.AdminOnly(userService), mc.CreateBatch)
	monevGroup.Get("/batches", mc.ListBatches) // semua role bisa lihat
	monevGroup.Get("/batches/:id", mc.GetBatch)
	monevGroup.Patch("/batches/:id", middleware.AdminOnly(userService), mc.UpdateBatch)
	monevGroup.Delete("/batches/:id", middleware.AdminOnly(userService), mc.DeleteBatch)

	// Admin — manage entries
	monevGroup.Post("/batches/:id/opds", middleware.AdminOnly(userService), mc.AddOpd)
	monevGroup.Get("/batches/:id/entries", mc.ListBatchEntries)
	monevGroup.Delete("/entries/:entryID", middleware.AdminOnly(userService), mc.RemoveEntry)
	monevGroup.Post("/entries/:entryID/verify", middleware.AdminOnly(userService), mc.VerifyEntry)

	// OPD + Admin — entry detail
	monevGroup.Get("/entries/:entryID", mc.GetEntry)
	monevGroup.Get("/my-entries", middleware.OPDOnly(userService), mc.ListMyEntries)
	monevGroup.Patch("/entries/:entryID/representative", middleware.OPDOnly(userService), mc.FillRepresentative)
	monevGroup.Post("/entries/:entryID/submit", middleware.OPDOnly(userService), mc.SubmitEntry)

	// Kajian
	monevGroup.Post("/entries/:entryID/kajians", mc.AddKajian)
	monevGroup.Patch("/kajians/:kajianID", mc.UpdateKajian)
	monevGroup.Delete("/kajians/:kajianID", mc.DeleteKajian)

	// Points
	monevGroup.Post("/kajians/:kajianID/points", middleware.AdminOnly(userService), mc.AddPoint)
	monevGroup.Patch("/points/:pointID", mc.UpdatePoint)
	monevGroup.Delete("/points/:pointID", middleware.AdminOnly(userService), mc.DeletePoint)
}
