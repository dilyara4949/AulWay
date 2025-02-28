package middleware

import (
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func AccessCheckMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get(RoleContextKey).(string)
			if !ok {
				log.Println("Role not found in context")
				return echo.NewHTTPError(http.StatusForbidden, "access denied")
			}

			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					return next(c)
				}
			}

			log.Printf("Access denied for role: %s", role)
			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
	}
}
