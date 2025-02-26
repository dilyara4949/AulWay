package bus

import (
	"aulway/internal/domain"
	"aulway/internal/handler/bus/model"
	"aulway/internal/utils/config"
	"aulway/internal/utils/errs"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service interface {
	CreateBus(ctx context.Context, request model.CreateRequest) (*domain.Bus, error)
	Get(ctx context.Context, id string) (*domain.Bus, error)
	GetByNumber(ctx context.Context, number string) (*domain.Bus, error)
}

// CreateBusHandler
// @Description Create Bus
// @Tags bus
// @Accept json
// @Produce json
// @Param requestBody body model.CreateRequest true
// @Success 200 {object} domain.Bus  "Success"
// @Failure 400 {object} errs.Err
// @Failure 500 {object} errs.Err
// @Router /bus [post]
func CreateBusHandler(busService Service, cfg config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		//log
		var request model.CreateRequest

		if err := c.Bind(&request); err != nil {
			//log
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "Binding request body failed", ErrDesc: err.Error()})
		}

		if err := request.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "Bad request", ErrDesc: err.Error()})
		}

		bus, err := busService.CreateBus(c.Request().Context(), request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Failed to create bus", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, bus)
	}
}

func GetBusHandler(busService Service, cfg config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		busId := c.Param("busId")

		bus, err := busService.Get(c.Request().Context(), busId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Failed to get bus", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, bus)
	}
}
