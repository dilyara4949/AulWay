package user

import (
	"aulway/internal/domain"
	auth "aulway/internal/handler/auth/model"
	"aulway/internal/handler/pagination"
	"aulway/internal/handler/user/model"
	"aulway/internal/utils/errs"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service interface {
	CreateUser(ctx context.Context, request auth.SignupRequest) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByFbUid(ctx context.Context, uid string) (*domain.User, error)
	GetUserById(ctx context.Context, uid string) (*domain.User, error)
	UpdateUser(ctx context.Context, req model.UpdateUserRequest, id string) error
	ResetPassword(ctx context.Context, password auth.ResetPassword) error
	GetUsers(ctx context.Context, page, pageSize int) ([]domain.User, error)
	ValidateUser(ctx context.Context, signin auth.SigninRequest) (*domain.User, error)
}

// UpdateUserHandler updates user information
// @Summary Update user details
// @Description Updates user information based on the given user ID
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param request body model.UpdateUserRequest true "User update request body"
// @Security BearerAuth
// @Success 200 {object} nil "User updated successfully"
// @Failure 400 {object} errs.Err "Invalid request body"
// @Failure 500 {object} errs.Err "Internal server error"
// @Router /api/users/{userId} [put]
func UpdateUserHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userId")

		req := model.UpdateUserRequest{}

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "Binding request body failed", ErrDesc: err.Error()})
		}

		err := service.UpdateUser(c.Request().Context(), req, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Update user failed", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, nil)
	}
}

// GetUserByIdHandler retrieves user information by ID
// @Summary Get user details
// @Description Fetches user details based on the provided user ID
// @Tags users
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} domain.User "User details retrieved successfully"
// @Failure 500 {object} errs.Err "Internal server error"
// @Router /api/users/{userId} [get]
func GetUserByIdHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Param("userId")

		user, err := service.GetUserById(c.Request().Context(), userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Get user failed", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, user)
	}
}

// GetUsersList retrieves a paginated list of users.
// @Summary Get list of users
// @Description Retrieves a paginated list of users.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1) minimum(1)
// @Param page_size query int false "Number of users per page" default(10) minimum(1) maximum(100)
// @Success 200 {array} domain.User "List of users"
// @Failure 500 {object} errs.Err "Failed to retrieve users"
// @Router /api/users [get]
func GetUsersList(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, pageSize := pagination.GetPageInfo(c)

		users, err := service.GetUsers(c.Request().Context(), page, pageSize)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Get users failed", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, users)
	}
}
