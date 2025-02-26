package route

import (
	"aulway/internal/domain"
	busServ "aulway/internal/handler/bus"
	"aulway/internal/handler/pagination"
	"aulway/internal/handler/route/model"
	"aulway/internal/utils/config"
	"aulway/internal/utils/errs"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type Service interface {
	CreateRoute(ctx context.Context, request model.CreateRouteRequest, availableSeats int) (*domain.Route, error)
	GetRoute(ctx context.Context, id string) (*domain.Route, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, req model.UpdateRouteRequest, id string) error
	GetRoutesList(ctx context.Context, departure, destination string, date time.Time, passengers, page, pageSize int) ([]domain.Route, int, error)
}

// CreateRouteHandler
// @Summary Create Route
// @Description Create a new bus route
// @Tags route
// @Accept json
// @Produce json
// @Param requestBody body model.CreateRouteRequest true "Route creation request"
// @Success 200 {object} domain.Route "Success"
// @Failure 400 {object} errs.Err "Bad Request"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /route [post]
func CreateRouteHandler(routeService Service, busService busServ.Service, _ config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var request model.CreateRouteRequest

		if err := c.Bind(&request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errs.Err{Err: "Binding request body failed", ErrDesc: err.Error()})
		}

		if err := request.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "Invalid request", ErrDesc: err.Error()})
		}

		bus, err := busService.Get(c.Request().Context(), request.BusId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "Get bus failed", ErrDesc: err.Error()})
		}

		route, err := routeService.CreateRoute(c.Request().Context(), request, bus.TotalSeats)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Create route failed", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, route)
	}
}

// GetRouteHandler
// @Summary Get Route
// @Description Retrieve a route by its ID
// @Tags route
// @Accept json
// @Produce json
// @Param routeId path string true "Route ID"
// @Success 200 {object} domain.Route "Success"
// @Failure 400 {object} errs.Err "Bad Request"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /route/{routeId} [get]
func GetRouteHandler(routeService Service, _ config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		routeId := c.Param("routeId")

		bus, err := routeService.GetRoute(c.Request().Context(), routeId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Failed to get bus", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, bus)
	}
}

// DeleteRouteHandler
// @Summary Delete Route
// @Description Delete a route by its ID
// @Tags route
// @Accept json
// @Produce json
// @Param routeId path string true "Route ID"
// @Success 200 {string} string "Success"
// @Failure 400 {object} errs.Err "Bad Request"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /route/{routeId} [delete]
func DeleteRouteHandler(routeService Service, _ config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		routeId := c.Param("routeId")

		err := routeService.Delete(c.Request().Context(), routeId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Failed to delete route", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, nil)
	}
}

// UpdateRouteHandler
// @Summary Update Route
// @Description Update an existing route by its ID
// @Tags route
// @Accept json
// @Produce json
// @Param routeId path string true "Route ID"
// @Param requestBody body model.UpdateRouteRequest true "Update Route Request Body"
// @Success 200 {object} string "Route updated successfully"
// @Failure 400 {object} errs.Err "Bad Request"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /route/{routeId} [put]
func UpdateRouteHandler(routeService Service, _ config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		routeId := c.Param("routeId")

		var request model.UpdateRouteRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "Binding request body failed", ErrDesc: err.Error()})
		}

		err := routeService.Update(c.Request().Context(), request, routeId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Failed to update route", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, nil)
	}
}

// GetRoutesListHandler
// @Summary Get Routes List
// @Description Retrieve a list of available routes based on filters
// @Tags route
// @Accept json
// @Produce json
// @Param departure query string true "Departure location"
// @Param destination query string true "Destination location"
// @Param date query string true "Travel date (format: YYYY-MM-DD)"
// @Param passengers query int true "Number of passengers"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param pageSize query int false "Page size for pagination (default: 30)"
// @Success 200 {array} []domain.Route "List of routes"
// @Failure 400 {object} errs.Err "Bad Request"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /route [get]
func GetRoutesListHandler(routeService Service, _ config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		departure := c.QueryParam("departure")
		destination := c.QueryParam("destination")
		dateStr := c.QueryParam("date")
		passengersStr := c.QueryParam("passengers")

		if departure == "" || destination == "" || dateStr == "" || passengersStr == "" {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "Failed to get routes", ErrDesc: "Missing required query parameters"})
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "Failed to get routes", ErrDesc: "Invalid date format, expected YYYY-MM-DD"})
		}

		passengers, err := strconv.Atoi(passengersStr)
		if err != nil || passengers <= 0 {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "Failed to get routes", ErrDesc: "Invalid passengers count"})
		}

		page, pageSize := pagination.GetPageInfo(c)

		routes, _, err := routeService.GetRoutesList(c.Request().Context(), departure, destination, date, passengers, page, pageSize)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Failed to get routes", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, routes)
	}
}
