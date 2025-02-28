package auth

import (
	"aulway/internal/handler/auth/model"
	"aulway/internal/handler/user"
	"aulway/internal/repository/errs"
	uerrs "aulway/internal/utils/errs"
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strings"
)

const (
	AdminRole = "admin"
	UserRole  = "user"
)

type Service interface {
	VerifyFirebaseToken(client *auth.Client, idToken string) (*auth.Token, error)
}

//type UserService interface {
//	CreateUser(ctx context.Context, email string, uid string) (*domain.User, error)
//	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
//	GetUserByFbUid(ctx context.Context, uid string) (*domain.User, error)
//	ResetPassword(ctx context.Context, password model.ResetPassword) error
//}

func FirebaseSignIn(userService user.Service, firebaseClient *auth.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, uerrs.Err{Err: "Unauthorized", ErrDesc: "Missing or invalid Authorization header"})
		}

		fbToken := strings.TrimPrefix(authHeader, "Bearer ")
		if fbToken == "" {
			return c.JSON(http.StatusUnauthorized, uerrs.Err{Err: "Unauthorized", ErrDesc: "Firebase token is missing"})
		}

		token, err := firebaseClient.VerifyIDToken(c.Request().Context(), fbToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, uerrs.Err{Err: "Unauthorized", ErrDesc: "Invalid Firebase token"})
		}

		fbUid := token.UID
		email, ok := token.Claims["email"].(string)
		if !ok || email == "" {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "Bad Request", ErrDesc: "Email is required"})
		}

		user, err := userService.GetUserByFbUid(c.Request().Context(), fbUid)
		if err != nil {
			if errors.Is(err, errs.ErrRecordNotFound) {
				if emailVerified, ok := token.Claims["email_verified"].(bool); !ok || !emailVerified {
					return c.JSON(http.StatusForbidden, uerrs.Err{Err: "Forbidden", ErrDesc: "Email must be verified"})
				}

				if err := AssignRole(firebaseClient, fbUid, UserRole); err != nil {
					return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "Internal Server Error", ErrDesc: "Failed to assign role"})
				}

				user, err = userService.CreateUser(c.Request().Context(), email, fbUid)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "Internal Server Error", ErrDesc: "Error creating user"})
				}
			} else {
				return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "Internal Server Error", ErrDesc: err.Error()})
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Sign-in successful",
			"user":    user,
		})
	}
}

func AssignRole(authClient *auth.Client, uid string, role string) error {
	claims := map[string]interface{}{
		"role": role,
	}
	err := authClient.SetCustomUserClaims(context.Background(), uid, claims)
	if err != nil {
		log.Printf("Error setting role for user %s: %v", uid, err)
		return err
	}
	log.Printf("Role '%s' assigned to user %s", role, uid)
	return nil
}

func ResetPasswordHandler(userService user.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.ResetPassword

		err := c.Bind(&req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "error at binding request body", ErrDesc: err.Error()})
		}

		if req.NewPassword == "" || req.OldPassword == "" || req.Email == "" {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: errs.ErrInvalidEmailPassword})
		}

		err = userService.ResetPassword(c.Request().Context(), req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "Failed to reset password", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, "reset password succeeded")
	}
}
