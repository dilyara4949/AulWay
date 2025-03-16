package bus

import (
	"aulway/internal/domain"
	"aulway/internal/handler/bus/model"
	"aulway/internal/handler/pagination"
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
	GetBusesList(ctx context.Context, page, pageSize int) ([]domain.Bus, error)
	DeleteBus(ctx context.Context, id string) error
}

// CreateBusHandler
// @Description Create Bus
// @Tags bus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param requestBody body model.CreateRequest true "Request Body"
// @Success 200 {object} domain.Bus  "Success"
// @Failure 400 {object} errs.Err
// @Failure 500 {object} errs.Err
// @Router /api/buses [post]
func CreateBusHandler(busService Service, cfg config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
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

// GetBusHandler
// @Summary Get Bus
// @Description Retrieve a bus by its ID
// @Tags bus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param busId path string true "Bus ID"
// @Success 200 {object} domain.Bus "Success"
// @Failure 400 {object} errs.Err "Bad Request"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /api/buses/{busId} [get]
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

// GetBusesListHandler
// @Summary Get Buses List
// @Description Retrieve a list of buses
// @Tags bus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number for pagination (default: 1)"
// @Param pageSize query int false "Page size for pagination (default: 30)"
// @Success 200 {array} []domain.Bus "List of buses"
// @Failure 400 {object} errs.Err "Bad Request"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /api/buses [get]
func GetBusesListHandler(busService Service, _ config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, pageSize := pagination.GetPageInfo(c)

		routes, err := busService.GetBusesList(c.Request().Context(), page, pageSize)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Failed to get buses", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, routes)
	}
}

// DeleteBusHandler
// @Summary Delete bus by id
// @Tags bus
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param busId path string true "Bus ID"
// @Success 200 {array} string "bus deleted"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /api/buses/{busId} [delete]
func DeleteBusHandler(busService Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		busId := c.Param("busId")

		err := busService.DeleteBus(c.Request().Context(), busId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Failed to delete bus", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, "bus deleted")
	}
}
