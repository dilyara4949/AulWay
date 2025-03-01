package middleware

import (
	"aulway/internal/utils/errs"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

const (
	UserIDKey   = "user_id"
	UserRoleKey = "user_role"
)

func JWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, errs.Err{Err: "authorization failed", ErrDesc: "authorization header required"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, errs.Err{Err: "authorization failed", ErrDesc: "Bearer token required"})
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, errs.Err{Err: "authorization failed", ErrDesc: "invalid token"})
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				claimedUID, ok := claims["user_id"].(string)
				if !ok {
					slog.Error("authorization", "error", "no user property in claims")
					return c.JSON(http.StatusBadRequest, errs.Err{Err: "authorization failed", ErrDesc: "invalid token"})
				}

				claimedRole, ok := claims[UserRoleKey].(string)
				if !ok {
					slog.Error("authorization", "error", "no role property in claims")
					return c.JSON(http.StatusBadRequest, errs.Err{Err: "authorization failed", ErrDesc: "invalid token"})
				}

				c.Set(UserIDKey, claimedUID)
				c.Set(UserRoleKey, claimedRole)

				return next(c)
			}
			return c.JSON(http.StatusUnauthorized, errs.Err{Err: "authorization failed", ErrDesc: "invalid token"})
		}
	}
}
