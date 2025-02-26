package healthz

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func CheckHealth() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	}
}
