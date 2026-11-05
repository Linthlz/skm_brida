package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func currentUserID(ctx *fiber.Ctx) (int64, bool) {
	token, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || token == nil {
		return 0, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}
	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, false
	}
	return int64(id), true
}

func isUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
