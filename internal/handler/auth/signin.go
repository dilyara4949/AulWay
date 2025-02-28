package auth

import (
	"aulway/internal/handler/auth/model"
	"aulway/internal/handler/user"
	"aulway/internal/utils/config"
	"aulway/internal/utils/errs"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

// SigninHandler
// @Summary User Signin
// @Description Authenticate a user and return an access token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.SigninRequest true "Signin Request Body"
// @Success 200 {object} model.SigninResponse
// @Failure 400 {object} errs.Err "Bad Request - Invalid request body"
// @Failure 403 {object} errs.Err "Forbidden - Password reset required"
// @Failure 404 {object} errs.Err "Not Found - User not found or incorrect credentials"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /auth/signin [post]
func SigninHandler(authService Service, userService user.Service, cfg config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.SigninRequest

		err := c.Bind(&req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "failed to signin", ErrDesc: err.Error()})
		}

		if req.Email == "" || req.Password == "" {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "failed to signin", ErrDesc: "fields cannot be empty"})
		}

		usr, err := userService.ValidateUser(c.Request().Context(), req)
		if err != nil {
			return c.JSON(http.StatusNotFound, errs.Err{Err: "failed to signin", ErrDesc: err.Error()})
		}

		if usr.RequirePasswordReset {
			return c.JSON(http.StatusForbidden, errs.Err{Err: "access denied", ErrDesc: "reset password required"})
		}

		token, err := authService.CreateAccessToken(c.Request().Context(), *usr, cfg.JWTTokenSecret, cfg.AccessTokenExpire)
		if err != nil {
			slog.Error("signin: error at creating access token,", "error", err.Error())
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "create access token error", ErrDesc: err.Error()})
		}

		resp := model.SigninResponse{
			AccessToken: token,
			User:        domainUserToResponse(*usr),
		}
		return c.JSON(http.StatusOK, resp)
	}
}
