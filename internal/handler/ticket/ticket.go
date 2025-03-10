package ticket

import (
	"aulway/internal/domain"
	"aulway/internal/handler/ticket/model"
	"aulway/internal/utils/errs"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service interface {
	BuyTickets(ctx context.Context, userID, routeID string, paymentMethodID string, ticketAmount int) ([]domain.Ticket, error)
}

// BuyTicketHandler processes ticket purchase requests for multiple tickets.
// @Summary      Buy tickets
// @Description  Allows a user to purchase one or more tickets for a specific route using card details.
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Param        routeId  path      string                     true  "Route ID"
// @Param        payment_id query  string                      true  "Payment method ID"
// @Param        requestBody body   model.BuyTicketRequest     true  "Buy Ticket Request Body"
// @Security     BearerAuth
// @Success      200      {array}   domain.Ticket             "Successfully purchased tickets"
// @Failure      400      {object}  errs.Err                  "Invalid request or request binding failed"
// @Failure      500      {object}  map[string]string         "Internal server error"
// @Router       /api/tickets/{routeId} [post]
func BuyTicketHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.BuyTicketRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "invalid request", ErrDesc: "request binding failed"})
		}

		paymentId := c.QueryParam("payment_id")
		routeId := c.Param("routeId")
		userID := c.Get("user_id").(string)

		tickets, err := service.BuyTickets(c.Request().Context(), userID, routeId, paymentId, req.Quantity)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "error", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, tickets)
	}
}

// DownloadTicketHandler serves the ticket as a PDF file
//func DownloadTicketHandler(service Service) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		ticketID := c.Param("ticketID")
//
//		// Fetch ticket from DB
//		ticket, err := service.GetTicket(c.Request().Context(), ticketID)
//		if err != nil {
//			return c.JSON(http.StatusNotFound, map[string]string{"error": "Ticket not found"})
//		}
//
//		// Generate PDF
//		pdf := gofpdf.New("P", "mm", "A4", "")
//		pdf.AddPage()
//		pdf.SetFont("Arial", "B", 16)
//		pdf.Cell(40, 10, "Ticket Details")
//
//		pdf.Ln(10) // New line
//		pdf.SetFont("Arial", "", 12)
//		pdf.Cell(40, 10, fmt.Sprintf("Ticket ID: %s", ticket.ID))
//		pdf.Ln(8)
//		pdf.Cell(40, 10, fmt.Sprintf("User ID: %s", ticket.UserID))
//		pdf.Ln(8)
//		pdf.Cell(40, 10, fmt.Sprintf("Route ID: %s", ticket.RouteID))
//		pdf.Ln(8)
//		pdf.Cell(40, 10, fmt.Sprintf("Price: $%d", ticket.Price))
//		pdf.Ln(8)
//		pdf.Cell(40, 10, fmt.Sprintf("Status: %s", ticket.Status))
//		pdf.Ln(8)
//		pdf.Cell(40, 10, fmt.Sprintf("Payment Status: %s", ticket.PaymentStatus))
//		pdf.Ln(8)
//		pdf.Cell(40, 10, fmt.Sprintf("Generated on: %s", time.Now().Format(time.RFC1123)))
//
//		// Serve PDF as response
//		c.Response().Header().Set("Content-Type", "application/pdf")
//		c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=ticket_%s.pdf", ticket.ID))
//		return pdf.Output(c.Response().Writer)
//	}
//}
