package favorite

import (
	"aulway/internal/domain"
	"aulway/internal/handler/access"
	"aulway/internal/handler/favorites/model/model"
	"aulway/internal/handler/pagination"
	"aulway/internal/utils/errs"
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type Service interface {
	Add(ctx context.Context, favorite *domain.FavoriteRoute) error
	Remove(ctx context.Context, favorite *domain.FavoriteRoute) error
	GetFavoriteRoutes(ctx context.Context, userId string, page, pageSize int) ([]domain.Route, error)
}

// AddFavoriteHandler
// @Summary Add to Favorites
// @Description Adds a route to the user's favorites
// @Tags favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param request body model.AddFavoriteRequest true "Favorite Route Request"
// @Success 200 {object} map[string]string "Added to favorites"
// @Failure 400 {object} errs.Err "Bad Request"
// @Failure 403 {object} errs.Err "Access Denied"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /api/users/{userId}/favorites [post]
func AddFavoriteHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !access.Check(c, c.Get("user_id"), "userId") {
			return c.JSON(http.StatusForbidden, errs.Err{Err: "add favorite failed", ErrDesc: "access denied"})
		}

		var req model.AddFavoriteRequest
		if err := c.Bind(&req); err != nil || req.RouteID == "" {
			log.Printf("error: %v", err)
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "invalid request", ErrDesc: "route_id is required"})
		}

		userID := c.Param("userId")

		fav := &domain.FavoriteRoute{
			UserID:  userID,
			RouteID: req.RouteID,
		}

		if err := service.Add(c.Request().Context(), fav); err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "failed to add to favorites", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "added to favorites"})
	}
}

// RemoveFavoriteHandler
// @Summary Remove from Favorites
// @Description Removes a route from the user's favorites
// @Tags favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param routeId path string true "Route ID"
// @Success 200 {object} map[string]string "Removed from favorites"
// @Failure 400 {object} errs.Err "Bad Request"
// @Failure 403 {object} errs.Err "Access Denied"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /api/users/{userId}/favorites/{routeId} [delete]
func RemoveFavoriteHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !access.Check(c, c.Get("user_id"), "userId") {
			return c.JSON(http.StatusForbidden, errs.Err{Err: "remove favorite failed", ErrDesc: "access denied"})
		}

		userID := c.Param("userId")
		routeID := c.Param("routeId")
		if routeID == "" {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "invalid route_id", ErrDesc: "route_id is required"})
		}

		fav := &domain.FavoriteRoute{
			UserID:  userID,
			RouteID: routeID,
		}

		err := service.Remove(c.Request().Context(), fav)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "failed to remove from favorites", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "removed from favorites"})
	}
}

// GetFavoritesHandler
// @Summary Get Favorite Routes
// @Description Returns a list of user's favorite routes with route and bus info
// @Tags favorites
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param pageSize query int false "Page size for pagination (default: 30)"
// @Success 200 {array} domain.Route "List of favorite routes"
// @Failure 403 {object} errs.Err "Access Denied"
// @Failure 500 {object} errs.Err "Internal Server Error"
// @Router /api/users/{userId}/favorites [get]
func GetFavoritesHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !access.Check(c, c.Get("user_id"), "userId") {
			return c.JSON(http.StatusForbidden, errs.Err{Err: "get favorites failed", ErrDesc: "access denied"})
		}

		userID := c.Param("userId")

		page, pageSize := pagination.GetPageInfo(c)

		routes, err := service.GetFavoriteRoutes(c.Request().Context(), userID, page, pageSize)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "failed to get favorites", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, routes)
	}
}
