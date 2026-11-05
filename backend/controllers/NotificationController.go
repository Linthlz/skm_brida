package controllers

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/itsXanadium/singa-riset-backend/models"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
)

type NotificationController struct {
	service services.NotificationService
}

func NewNotificationController(service services.NotificationService) *NotificationController {
	return &NotificationController{service: service}
}

func (c *NotificationController) FetchNotifications(ctx *fiber.Ctx) error {
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userInternalID := int64(claims["user_id"].(float64))
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "5"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 5
	}
	offset := (page - 1) * limit
	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")
	notifications, total, err := c.service.NotifPiles(userInternalID, filter, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to fetch Data", err.Error())
	}
	response := make([]models.NotificationOutput, 0, len(notifications))
	for _, n := range notifications {
		response = append(response, models.NotificationOutput{
			PublicID:          n.PublicID,
			Subject:           n.Subject,
			Message:           n.Message,
			ReferenceType:     n.ReferenceType,
			ReferencePublicID: n.ReferencePublicID,
			CreatedAt:         n.CreatedAt,
		})
	}
	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}
	return utils.SuccessPaginated(ctx, "Notifications Fetched", response, meta)
}
