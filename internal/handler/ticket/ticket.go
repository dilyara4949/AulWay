package ticket

import (
	"aulway/internal/domain"
	"aulway/internal/utils/errs"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service interface {
	BuyTicket(ctx context.Context, userID, routeID string, paymentMethodID string) (*domain.Ticket, error)
}

type PaymentService interface {
}

// BuyTicketHandler processes ticket purchase requests.
// @Summary      Buy a ticket
// @Description  Allows a user to purchase a ticket for a specific route using card details.
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Param        routeId  path      string         true  "Route ID"
// @Param payment_id query string true "Payment method id"
// @Security     BearerAuth
// @Success      201      {object}  domain.Ticket  "Successfully purchased ticket"
// @Failure      400      {object}  errs.Err       "Invalid request or request binding failed"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /api/tickets/{routeId} [post]
func BuyTicketHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		//var req domain.CardDetails
		//if err := c.Bind(&req); err != nil {
		//	return c.JSON(http.StatusBadRequest, errs.Err{Err: "invalid request", ErrDesc: "request binding failed"})
		//}

		paymentId := c.QueryParam("payment_id")

		routeId := c.Param("routeId")
		userID := c.Get("user_id").(string)

		ticket, err := service.BuyTicket(c.Request().Context(), userID, routeId, paymentId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "error", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusCreated, ticket)
	}
}
