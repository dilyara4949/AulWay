package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
)

func FirebaseAuthMiddleware(authClient *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing auth token")
			}

			// Extract token from "Bearer <token>"
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth header")
			}

			// Verify token with Firebase
			token, err := authClient.VerifyIDToken(context.Background(), tokenStr)
			if err != nil {
				log.Println("Firebase token verification failed:", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Extract role from Firebase claims
			role, ok := token.Claims["role"].(string)
			if !ok {
				role = "user" // Default role if not set
			}

			// Store user ID & role in context
			c.Set("user_id", token.UID)
			c.Set("role", role)

			return next(c)
		}
	}
}
