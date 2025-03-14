package user

import (
	"aulway/internal/domain"
	"aulway/internal/handler/access"
	authModel "aulway/internal/handler/auth/model"
	"aulway/internal/handler/pagination"
	"aulway/internal/handler/user/model"
	rerrs "aulway/internal/repository/errs"
	"aulway/internal/utils/errs"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service interface {
	CreateUser(ctx context.Context, request authModel.SignupRequest) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByFbUid(ctx context.Context, uid string) (*domain.User, error)
	GetUserById(ctx context.Context, uid string) (*domain.User, error)
	UpdateUser(ctx context.Context, req model.UpdateUserRequest, id string) (*domain.User, error)
	ResetPassword(ctx context.Context, password model.ResetPasswordRequest, requirePasswordReset bool) error
	GetUsers(ctx context.Context, page, pageSize int) ([]domain.User, error)
	ValidateUser(ctx context.Context, signin authModel.SigninRequest) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
}

// ChangePasswordHandler change user password
// @Summary      Change password
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userId   path      string  true  "User ID"
// @Success 200 {string} string "password change was successful"
// @Failure      400      {object}  errs.Err
// @Failure      500      {object}  errs.Err
// @Router       /api/users/{userId}/change-password [put]
func ChangePasswordHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !access.Check(c, c.Get("user_id"), "userId") {
			return c.JSON(http.StatusForbidden, errs.Err{Err: "Update user failed", ErrDesc: "access denied"})
		}

		var req model.ResetPasswordRequest

		err := c.Bind(&req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "error at binding request body", ErrDesc: err.Error()})
		}

		if req.NewPassword == "" || req.OldPassword == "" || req.Email == "" {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "request body is incorrect"})
		}

		err = service.ResetPassword(c.Request().Context(), req, false)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "reset password failed", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, "password change was successful")
	}
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
// @Success 200 {object} model.UserResponse "User updated successfully"
// @Failure 400 {object} errs.Err "Invalid request body"
// @Failure 500 {object} errs.Err "Internal server error"
// @Router /api/users/{userId} [put]
func UpdateUserHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !access.Check(c, c.Get("user_id"), "userId") {
			return c.JSON(http.StatusForbidden, errs.Err{Err: "Update user failed", ErrDesc: "access denied"})
		}

		userId := c.Param("userId")

		req := model.UpdateUserRequest{}

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "Binding request body failed", ErrDesc: err.Error()})
		}

		usr, err := service.UpdateUser(c.Request().Context(), req, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Update user failed", ErrDesc: err.Error()})
		}

		response := domainUserToResponse(*usr)

		return c.JSON(http.StatusOK, response)
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
		if !access.Check(c, c.Get("user_id"), "userId") {
			return c.JSON(http.StatusForbidden, errs.Err{Err: "Update user failed", ErrDesc: "access denied"})
		}

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

// DeleteUserHandler deletes a user by ID.
// @Summary      Delete a user
// @Description  Deletes a user by their ID. Only admin or the user themselves can delete.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userId   path      string  true  "User ID"
// @Success 200 {string} string "User successfully deleted"
// @Failure      400      {object}  errs.Err "Invalid user ID"
// @Failure      403      {object}  errs.Err "Unauthorized"
// @Failure      404      {object}  errs.Err "User not found"
// @Failure      500      {object}  errs.Err "Failed to delete user"
// @Router       /api/users/{userId} [delete]
func DeleteUserHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !access.Check(c, c.Get("user_id"), "userId") {
			return c.JSON(http.StatusForbidden, errs.Err{Err: "Delete user failed", ErrDesc: "access denied"})
		}

		userID := c.Param("userId")

		err := service.DeleteUser(c.Request().Context(), userID)
		if errors.Is(err, rerrs.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, errs.Err{Err: "error", ErrDesc: "User not found"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "error", ErrDesc: "Failed to delete user"})
		}

		return c.JSON(http.StatusOK, "user deleted")
	}
}

func domainUserToResponse(user domain.User) model.UserResponse {
	return model.UserResponse{
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
