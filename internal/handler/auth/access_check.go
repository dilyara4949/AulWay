package auth

import "github.com/labstack/echo/v4"

const (
	adminRole   = "admin"
	userRoleKey = "user_role"
)

func AccessCheck(c echo.Context, expectedContextID, expectedIDKey string) bool {
	role := c.Get(userRoleKey)
	if role == nil {
		return false
	}

	userRole, ok := role.(string)
	if !ok {
		return false
	}

	userID := c.Param(expectedIDKey)
	if userRole == adminRole || (expectedContextID == userID && expectedContextID != "") {
		return true
	}

	return false
}
