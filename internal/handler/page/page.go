package page

import (
	"aulway/internal/domain"
	"aulway/internal/handler/page/model"
	"aulway/internal/utils/errs"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service interface {
	GetPage(ctx context.Context, title string) (domain.Page, error)
	UpdatePage(ctx context.Context, title, content string) error
}

// GetPageHandler retrieves a page by title.
// @Summary      Get a page
// @Description  Retrieves a page's content by title. Available titles: "about_us", "privacy_policy", "help_support".
// @Tags         pages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        title   path      string  true  "Page Title"
// @Success      200     {object}  domain.Page
// @Failure      404     {object}  map[string]string "Page not found"
// @Router       /api/pages/{title} [get]
func GetPageHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		title := c.Param("title") // "about_us", "privacy_policy", "support"

		page, err := service.GetPage(c.Request().Context(), title)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Page not found"})
		}

		return c.JSON(http.StatusOK, page)
	}
}

// UpdatePageHandler updates a page content.
// @Summary      Update a page
// @Description  Updates the content of a page by title. Only admins can update.Available titles: "about_us", "privacy_policy", "help_support".
// @Tags         pages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        title   path      string  true  "Page Title"
// @Param        request body      model.UpdatePageRequest  true  "Page content"
// @Success      200     {object}  map[string]string "Page updated successfully"
// @Failure      400     {object}  map[string]string "Invalid request"
// @Failure      500     {object}  map[string]string "Failed to update page"
// @Router       /api/pages/{title} [put]
func UpdatePageHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		title := c.Param("title") // "about_us", "privacy_policy", "support"

		var req model.UpdatePageRequest

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "Invalid request", ErrDesc: err.Error()})
		}

		if err := service.UpdatePage(c.Request().Context(), title, req.Content); err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "Failed to update page", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Page updated successfully"})
	}
}
