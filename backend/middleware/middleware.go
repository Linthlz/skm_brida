package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/itsXanadium/singa-riset-backend/repositories"
	"github.com/itsXanadium/singa-riset-backend/services"
	"github.com/itsXanadium/singa-riset-backend/utils"
	"gorm.io/gorm"
)

// Fix, Unimplemented somehow
func Verified(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken, ok := c.Locals("user").(*jwt.Token)
		if !ok || userToken == nil {
			return utils.UnauthorizedAccess(c, "Error", "Unauthorized Access")
		}
		claims, ok := userToken.Claims.(jwt.MapClaims)
		if !ok {
			return utils.UnauthorizedAccess(c, "Error", "Token invalid!")
		}
		publicID, ok := claims["pub_id"].(string)
		if !ok {
			return utils.UnauthorizedAccess(c, "Unauthorized Access!", "")
		}
		user, err := s.FetchUserByPublicID(publicID)
		if err != nil {
			return utils.UnauthorizedAccess(c, "Unauthorized Access!", "user not found")
		}
		fmt.Println("{+}: User found")
		if !user.Verified {
			return utils.Forbidden(c, "Email not verified", "verify your email first")
		}
		return c.Next()
	}
}
func UserProfile(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken, ok := c.Locals("user").(*jwt.Token)
		if !ok || userToken == nil {
			return utils.UnauthorizedAccess(c, "Error", "Unauthorized Access")
		}
		claims, ok := userToken.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("{-}: Token Invalid")
			return utils.UnauthorizedAccess(c, "Error", "Token invalid!")
		}
		publicID, ok := claims["pub_id"].(string)
		ownerID := c.Params("id")
		if !ok {
			fmt.Println("{-}: Unathorized Access")
			return utils.UnauthorizedAccess(c, "Unauthorized Access!", "")
		}
		if publicID != ownerID {
			fmt.Println("{-}: Unathorized Access")
			return utils.UnauthorizedAccess(c, "Not the owner", "Unauthorized Access")
		}
		return c.Next()
	}

}
func OwnerShip(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken, ok := c.Locals("user").(*jwt.Token)
		if !ok || userToken == nil {
			return utils.UnauthorizedAccess(c, "Error", "Unauthorized Access")
		}
		claims, ok := userToken.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("{-}: Token Invalid")
			return utils.UnauthorizedAccess(c, "Error", "Token invalid!")
		}
		userInternalID := int64(claims["user_id"].(float64))
		userRole, _ := claims["role"].(string)
		publicID := c.Params("id")
		targetUser, err := s.FetchUserByPublicID(publicID)
		if err != nil || targetUser == nil {
			return utils.NotFound(c, "Not Found", "User not found")
		}
		isOwner := targetUser.InternalID == userInternalID
		isPrivileged := userRole == "admin" || userRole == "reviewer"

		if !isOwner && !isPrivileged {
			return utils.UnauthorizedAccess(c, "Forbidden", "You can only modify your own account")
		}
		c.Locals("user", targetUser)
		return c.Next()
	}

}
func AdminOnly(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// fmt.Println("=================================")
		// fmt.Println("{+}: Profile owner middleware reached")
		userToken, ok := c.Locals("user").(*jwt.Token)
		if !ok || userToken == nil {
			return utils.UnauthorizedAccess(c, "Error", "Unauthorized Access")
		}
		claims, ok := userToken.Claims.(jwt.MapClaims)
		// fmt.Println("{+}: Token Exist")
		if !ok {
			// fmt.Println("{-}: Token Invalid")
			return utils.UnauthorizedAccess(c, "Error", "Token invalid!")
		}
		roles, ok := claims["role"].(string)
		if !ok {
			return utils.UnauthorizedAccess(c, "Unauthorized Access", "")
		}
		if roles != "admin" {
			fmt.Println("{-}: Unathorized Access")
			return utils.Forbidden(c, "Forbidden", "Admin Only!")
		}
		// fmt.Println("==========End admin endpoint===============")
		return c.Next()
	}
}
func OPDOnly(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken, ok := c.Locals("user").(*jwt.Token)
		if !ok || userToken == nil {
			return utils.UnauthorizedAccess(c, "Access Forbidden", "Missing token")
		}
		claims, ok := userToken.Claims.(jwt.MapClaims)
		if !ok {
			return utils.UnauthorizedAccess(c, "Access Forbidden", "Invalid token")
		}
		role, _ := claims["role"].(string)
		if role != "opd" && role != "admin" {
			fmt.Println("{-}: Unathorized Access")
			return utils.Forbidden(c, "Forbidden", "OPD access required")
		}
		// Optional: validasi user punya opd_internal_id
		if role == "opd" {
			FuserID, ok := claims["user_id"].(float64)
			if !ok {
				return utils.UnauthorizedAccess(c, "Access Forbidden", "Invalid user_id claim")
			}
			user, err := s.FetchUserByID(uint(int64(FuserID)))
			if err != nil || user == nil {
				return utils.UnauthorizedAccess(c, "Access Forbidden", "User not found")
			}
			if user.OpdInternalID == nil {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error":   "Forbidden",
					"message": "Akun Anda belum terhubung ke OPD",
				})
			}
		}
		return c.Next()
	}
}

//	func OPDOnly() fiber.Handler {
//		return func(c *fiber.Ctx) error {
//			userToken, ok := c.Locals("user").(*jwt.Token)
//			if !ok || userToken == nil {
//				return utils.UnauthorizedAccess(c, "Access Forbidden", "Missing token")
//			}
//			claims, ok := userToken.Claims.(jwt.MapClaims)
//			if !ok {
//				return utils.UnauthorizedAccess(c, "Access Forbidden", "Invalid token")
//			}
//			role, _ := claims["role"].(string)
//			if role != "opd" && role != "admin" {
//				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
//					"error":   "Forbidden",
//					"message": "OPD access required",
//				})
//			}
//			return c.Next()
//		}
//	}
func BridaEmployee(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// fmt.Println("=================================")
		// fmt.Println("{+}: Profile owner middleware reached")
		userToken, ok := c.Locals("user").(*jwt.Token)
		if !ok || userToken == nil {
			return utils.UnauthorizedAccess(c, "Error", "Unauthorized Access")
		}
		claims, ok := userToken.Claims.(jwt.MapClaims)
		// fmt.Println("{+}: Token Exist")
		if !ok {
			// fmt.Println("{-}: Token Invalid")
			return utils.UnauthorizedAccess(c, "Error", "Token invalid!")
		}
		roles, ok := claims["role"].(string)
		if !ok {
			return utils.UnauthorizedAccess(c, "Unauthorized Access", "")
		}
		if roles != "pegawai-brida" {
			fmt.Println("{-}: Unathorized Access")
			return utils.Forbidden(c, "Forbidden", "Employee Only!")
		}
		// fmt.Println("==========End Brida Employee endpoint===============")
		return c.Next()
	}
}
func MediaEmployee(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// fmt.Println("=================================")
		// fmt.Println("{+}: Profile owner middleware reached")
		userToken, ok := c.Locals("user").(*jwt.Token)
		if !ok || userToken == nil {
			return utils.UnauthorizedAccess(c, "Error", "Unauthorized Access")
		}
		claims, ok := userToken.Claims.(jwt.MapClaims)
		// fmt.Println("{+}: Token Exist")
		if !ok {
			// fmt.Println("{-}: Token Invalid")
			return utils.UnauthorizedAccess(c, "Error", "Token invalid!")
		}
		roles, ok := claims["role"].(string)
		if !ok {
			return utils.UnauthorizedAccess(c, "Unauthorized Access", "")
		}
		if roles != "media-brida" && roles != "admin" {
			fmt.Println("{-}: Unathorized Access")
			return utils.Forbidden(c, "Forbidden", "Employee Only!")
		}
		// fmt.Println("==========End Brida Employee endpoint===============")
		return c.Next()
	}
}
func ReviewerRole(s services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// fmt.Println("=================================")
		// fmt.Println("{+}: Profile owner middleware reached")
		userToken, ok := c.Locals("user").(*jwt.Token)
		if !ok || userToken == nil {
			return utils.UnauthorizedAccess(c, "Error", "Unauthorized Access")
		}
		claims, ok := userToken.Claims.(jwt.MapClaims)
		// fmt.Println("{+}: Token Exist")
		if !ok {
			fmt.Println("{-}: Token Invalid")
			return utils.UnauthorizedAccess(c, "Error", "Token invalid!")
		}
		roles, ok := claims["role"].(string)
		if !ok {
			return utils.UnauthorizedAccess(c, "Unauthorized Access", "")
		}
		if roles != "reviewer" && roles != "admin" {
			fmt.Println("{-}: Unathorized Access")
			return utils.Forbidden(c, "Forbidden", "reviewer Only!")
		}
		// fmt.Println("==========End Reviewer endpoint===============")
		return c.Next()
	}
}

func GroupMemberOnly(groupRepo repositories.ResearchGroupRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken, ok := c.Locals("user").(*jwt.Token)
		if !ok || userToken == nil {
			return utils.UnauthorizedAccess(c, "Access Forbidden", "Missing token")
		}
		claims, ok := userToken.Claims.(jwt.MapClaims)
		if !ok {
			return utils.UnauthorizedAccess(c, "Access Forbidden", "Invalid token")
		}
		fUserID, ok := claims["user_id"].(float64)
		if !ok {
			return utils.UnauthorizedAccess(c, "Access Forbidden", "Invalid user_id claim")
		}
		userID := int64(fUserID)
		globalRole, _ := claims["role"].(string)
		groupPublicID := c.Params("id")
		if _, err := uuid.Parse(groupPublicID); err != nil {
			return utils.BadReq(c, "Invalid group ID", "id must be a valid UUID")
		}
		group, err := groupRepo.FetchByPublicID(groupPublicID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NotFound(c, "Group Not Found", "group not found")
		}
		if err != nil {
			return utils.InternalServerError(c, "Failed to fetch group", err.Error())
		}
		if group == nil {
			return utils.NotFound(c, "Group Not Found", "")
		}
		allowed := false
		groupRole := ""
		switch {
		case strings.EqualFold(globalRole, "admin"):
			allowed, groupRole = true, "admin"
		case group.LeaderInternalID == userID:
			allowed, groupRole = true, "leader"
		default:
			isMember, err := groupRepo.IsMember(group.InternalID, userID)
			if err != nil {
				return utils.InternalServerError(c, "Failed to check membership", err.Error())
			}
			if isMember {
				allowed, groupRole = true, "member"
			} else {
				isSupervisor, err := groupRepo.IsSupervisor(group.InternalID, userID)
				if err != nil {
					return utils.InternalServerError(c, "Failed to check membership", err.Error())
				}
				if isSupervisor {
					allowed, groupRole = true, "pembimbing"
				}
			}
		}
		if !allowed {
			return utils.Forbidden(c, "Forbidden", "You are not a member of this group")
		}
		c.Locals("group", group)
		c.Locals("user_id", userID)
		c.Locals("group_role", groupRole)
		return c.Next()
	}
}
