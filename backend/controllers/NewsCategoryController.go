package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

type NewsCategoryController struct {
	service services.NewsCategoryService
}

func NewNewsCategoryController(service services.NewsCategoryService) *NewsCategoryController {
	return &NewsCategoryController{service: service}
}

func (c *NewsCategoryController) List(ctx *fiber.Ctx) error {
	activeOnly := ctx.Query("active_only", "true") == "true"
	cats, err := c.service.List(activeOnly)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch categories", err.Error())
	}
	return utils.Success(ctx, "Categories fetched", cats)
}

func (c *NewsCategoryController) Create(ctx *fiber.Ctx) error {
	var body struct {
		Name         string `json:"name"`
		DisplayOrder int    `json:"display_order"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	cat, err := c.service.Create(models.NewsCategory{
		Name:         body.Name,
		DisplayOrder: body.DisplayOrder,
		IsActive:     true,
	})
	if err != nil {
		return utils.HandleError(ctx, "Cannot create category", err)
	}
	return utils.CreatedWithData(ctx, "Category created", cat)
}

func (c *NewsCategoryController) Update(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if !isUUID(publicID) {
		return utils.BadReq(ctx, "Invalid ID", "id must be a valid UUID")
	}
	var body struct {
		Name         string `json:"name"`
		DisplayOrder int    `json:"display_order"`
		IsActive     bool   `json:"is_active"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadReq(ctx, "Invalid request body", "failed to parse request body")
	}
	cat, err := c.service.Update(publicID, models.NewsCategory{
		Name:         body.Name,
		DisplayOrder: body.DisplayOrder,
		IsActive:     body.IsActive,
	})
	if err != nil {
		return utils.HandleError(ctx, "Cannot update category", err)
	}
	return utils.Success(ctx, "Category updated", cat)
}

func (c *NewsCategoryController) Delete(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if !isUUID(publicID) {
		return utils.BadReq(ctx, "Invalid ID", "id must be a valid UUID")
	}
	if err := c.service.Delete(publicID); err != nil {
		return utils.HandleError(ctx, "Cannot delete category", err)
	}
	return utils.Success(ctx, "Category deleted", nil)
}
