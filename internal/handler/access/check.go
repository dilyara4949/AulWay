package access

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

const (
	adminRole   = "admin"
	userRoleKey = "user_role"
)

func Check(c echo.Context, expectedContextID interface{}, expectedIDKey string) bool {
	contextIDStr := fmt.Sprintf("%v", expectedContextID)

	role := c.Get(userRoleKey)
	if role == nil {
		return false
	}

	userRole, ok := role.(string)
	if !ok {
		return false
	}

	userID := c.Param(expectedIDKey)
	if userRole == adminRole || (contextIDStr == userID && contextIDStr != "") {
		return true
	}

	return false
}
