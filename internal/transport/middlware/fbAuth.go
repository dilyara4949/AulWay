package middleware

import (
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
)

const (
	RoleContextKey     = "role"
	FbUserIdContextKey = "user_fid"
)

func FirebaseAuthMiddleware(authClient *auth.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				log.Println("Missing Authorization header")
				return echo.NewHTTPError(http.StatusUnauthorized, "missing auth token")
			}

			// Extract token from "Bearer <token>"
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader {
				log.Println("Invalid Authorization header format")
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid auth header")
			}

			// Verify token with Firebase
			token, err := authClient.VerifyIDToken(c.Request().Context(), tokenStr)
			if err != nil {
				log.Printf("Firebase token verification failed: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Extract role from Firebase claims
			role, ok := token.Claims[RoleContextKey].(string)
			if !ok {
				role = "user"
			}

			// Store user ID & role in context
			c.Set(FbUserIdContextKey, token.UID)
			c.Set(RoleContextKey, role)

			log.Printf("User authenticated: UID=%s, Role=%s", token.UID, role)

			return next(c)
		}
	}
}
