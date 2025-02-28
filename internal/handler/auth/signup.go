package auth

import (
	"aulway/internal/domain"
	"aulway/internal/handler/auth/model"
	"aulway/internal/handler/user"
	usermodel "aulway/internal/handler/user/model"
	"aulway/internal/utils/config"
	uerrs "aulway/internal/utils/errs"
	"context"
	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"
	"log"
	"log/slog"
	"net/http"
)

const (
	AdminRole = "admin"
	UserRole  = "user"
)

type Service interface {
	VerifyFirebaseToken(client *auth.Client, idToken string) (*auth.Token, error)
	CreateAccessToken(ctx context.Context, user domain.User, jwtSecret string, expiry int) (string, error)
}

// SignupHandler
// @Summary User Signup
// @Description Register a new user and return an access token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.SignupRequest true "Signup Request Body"
// @Success 200 {object} model.SignupResponse
// @Failure 400 {object} errs.Err "Bad Request - Invalid request body"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /auth/signup [post]
func SignupHandler(authService Service, userService user.Service, cfg config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.SignupRequest

		err := c.Bind(&req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "incorrect req body", ErrDesc: err.Error()})
		}

		if req.Password == "" || req.Email == "" {
			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "signup failed", ErrDesc: "fields cannot be empty"})
		}

		usr, err := userService.CreateUser(c.Request().Context(), req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "signup failed", ErrDesc: err.Error()})
		}

		token, err := authService.CreateAccessToken(c.Request().Context(), *usr, cfg.JWTTokenSecret, cfg.AccessTokenExpire)
		if err != nil {
			slog.Error("signup: error at creating access token,", "error", err.Error())

			return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "signup failed", ErrDesc: err.Error()})
		}

		resp := model.SignupResponse{
			AccessToken: token,
			User:        domainUserToResponse(*usr),
		}

		return c.JSON(http.StatusOK, resp)
	}
}

func domainUserToResponse(user domain.User) usermodel.UserResponse {
	return usermodel.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

//func domainUsersToResponse(users []domain.User) []response.User {
//	res := make([]response.User, 0)
//
//	for _, user := range users {
//		res = append(res, response.User{
//			ID:        user.ID,
//			Email:     user.Email,
//			Phone:     user.Phone,
//			CreatedAt: user.CreatedAt,
//			UpdatedAt: user.UpdatedAt,
//		})
//	}
//	return res
//}

// FirebaseSignIn handles Firebase authentication.
// @Summary Authenticate user via Firebase ID token
// @Description Verifies the Firebase ID token, retrieves or creates the user, and assigns roles if needed.
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Firebase ID Token prefixed with 'Bearer '"
// @Success 200 {object} domain.User "User authenticated successfully"
// @Failure 400 {object} errs.Err "Bad Request - Missing or invalid email"
// @Failure 401 {object} errs.Err "Unauthorized - Missing or invalid Authorization header"
// @Failure 403 {object} errs.Err "Forbidden - Email must be verified"
// @Failure 500 {object} errs.Err "Internal Server Error - Error processing request"
// @Router /auth/firebase-signin [post]
//func FirebaseSignIn(userService user.Service, firebaseClient *auth.Client) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		authHeader := c.Request().Header.Get("Authorization")
//		if !strings.HasPrefix(authHeader, "Bearer ") {
//			return c.JSON(http.StatusUnauthorized, uerrs.Err{Err: "Unauthorized", ErrDesc: "Missing or invalid Authorization header"})
//		}
//
//		fbToken := strings.TrimPrefix(authHeader, "Bearer ")
//		if fbToken == "" {
//			return c.JSON(http.StatusUnauthorized, uerrs.Err{Err: "Unauthorized", ErrDesc: "Firebase token is missing"})
//		}
//
//		token, err := firebaseClient.VerifyIDToken(c.Request().Context(), fbToken)
//		if err != nil {
//			return c.JSON(http.StatusUnauthorized, uerrs.Err{Err: "Unauthorized", ErrDesc: "Invalid Firebase token"})
//		}
//
//		fbUid := token.UID
//		email, ok := token.Claims["email"].(string)
//		if !ok || email == "" {
//			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "Bad Request", ErrDesc: "Email is required"})
//		}
//
//		usr, err := userService.GetUserByFbUid(c.Request().Context(), fbUid)
//		if err != nil {
//			if errors.Is(err, errs.ErrRecordNotFound) {
//				//if emailVerified, ok := token.Claims["email_verified"].(bool); !ok || !emailVerified {
//				//	return c.JSON(http.StatusForbidden, uerrs.Err{Err: "Forbidden", ErrDesc: "Email must be verified"})
//				//}
//
//				if err := AssignRole(firebaseClient, fbUid, UserRole); err != nil {
//					return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "Internal Server Error", ErrDesc: "Failed to assign role"})
//				}
//
//				usr, err = userService.CreateUser(c.Request().Context(), email, fbUid)
//				if err != nil {
//					return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "Internal Server Error", ErrDesc: "Error creating user"})
//				}
//			} else {
//				return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "Internal Server Error", ErrDesc: err.Error()})
//			}
//		}
//
//		return c.JSON(http.StatusOK, usr)
//	}
//}

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

//func ResetPasswordHandler(userService user.Service) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		var req model.ResetPassword
//
//		err := c.Bind(&req)
//		if err != nil {
//			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: "error at binding request body", ErrDesc: err.Error()})
//		}
//
//		if req.NewPassword == "" || req.OldPassword == "" || req.Email == "" {
//			return c.JSON(http.StatusBadRequest, uerrs.Err{Err: errs.ErrInvalidEmailPassword})
//		}
//
//		err = userService.ResetPassword(c.Request().Context(), req)
//		if err != nil {
//			return c.JSON(http.StatusInternalServerError, uerrs.Err{Err: "Failed to reset password", ErrDesc: err.Error()})
//		}
//
//		return c.JSON(http.StatusOK, "reset password succeeded")
//	}
//}
