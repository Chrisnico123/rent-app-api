package middleware

import (
	"rent-application/shared/constants"
	"rent-application/shared/helper"
	"rent-application/shared/web"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// PermissionConfig menyimpan mapping permission ke role IDs
type PermissionConfig map[string][]int

// DefaultPermissions konfigurasi standar
var DefaultPermissions = PermissionConfig{
	"public":    {constants.CustomerRole, constants.SellerRole, constants.AdminRole},
	"customer":  {constants.CustomerRole},
	"seller":    {constants.SellerRole},
	"admin":     {constants.AdminRole},
	"moderator": {constants.SellerRole, constants.AdminRole},
}

func RoleBasedAuth(permissions string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ekstrak token
		authHeader := c.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		if token == "" {
			return c.Status(401).JSON(web.ErrUnauthorized("Token required"))
		}

		// Parse JWT
		roleID, userID, err := helper.ParseJwt(token)
		if err != nil {
			return c.Status(401).JSON(web.ErrUnauthorized("Invalid token"))
		}

		intRoleID, _ := strconv.Atoi(roleID)

		// Cek permission
		if !hasPermission(intRoleID, permissions) {
			return c.Status(403).JSON(web.ErrForbidden(
				"Role " + constants.GetRoleName(intRoleID) + " cannot access this endpoint",
			))
		}

		// Simpan data user di context
		c.Locals("userID", userID)
		c.Locals("roleID", intRoleID)

		return c.Next()
	}
}

// Helper function untuk permission check
func hasPermission(userRole int, requiredRoles string) bool {
	for _, role := range DefaultPermissions[requiredRoles] {
		if userRole == role {
			return true
		}
	}
	return false
}
