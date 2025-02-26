package user

import (
	"aulway/internal/domain"
	"aulway/internal/handler/user/model"
	"aulway/internal/utils/errs"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service interface {
	CreateUser(ctx context.Context, email string, uid string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByFbUid(ctx context.Context, uid string) (*domain.User, error)
	UpdateUser(ctx context.Context, req model.UpdateUserRequest, id string) error
	//UpdateUserRole(ctx context.Context, uid string, role string) error
}

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
