package ticket

import (
	"aulway/internal/domain"
	"aulway/internal/handler/access"
	"aulway/internal/handler/pagination"
	"aulway/internal/handler/ticket/model"
	"aulway/internal/service"
	"aulway/internal/utils/config"
	"aulway/internal/utils/errs"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"time"
)

const (
	UpcomingTicket = "upcoming"
	PastTicket     = "past"
)

type Service interface {
	BuyTickets(ctx context.Context, userID, routeID string, paymentMethodID string, ticketAmount int) ([]domain.Ticket, *domain.Bus, *domain.Route, error)
	GetUpcomingTickets(ctx context.Context, userID string, now time.Time) ([]domain.Ticket, error)
	GetPastTickets(ctx context.Context, userID string, now time.Time) ([]domain.Ticket, error)
	TicketDetails(ctx context.Context, ticketId string) (*domain.Ticket, error)
	GetTicketsSortBy(ctx context.Context, sortBy, ord string, page, pageSize int) ([]domain.Ticket, error)
}

// BuyTicketHandler processes ticket purchase requests for multiple tickets.
// @Summary      Buy tickets
// @Description  Allows a user to purchase one or more tickets for a specific route using card details.
// @Tags         tickets
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
func BuyTicketHandler(s Service, cfg config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.BuyTicketRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, errs.Err{Err: "invalid request", ErrDesc: "request binding failed"})
		}

		paymentId := c.QueryParam("payment_id")
		routeId := c.Param("routeId")
		userID := c.Get("user_id").(string)

		tickets, bus, route, err := s.BuyTickets(c.Request().Context(), userID, routeId, paymentId, req.Quantity)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "error", ErrDesc: err.Error()})
		}

		go func() {
			emailBody := buildTicketEmailBody(tickets, bus, route)
			err := service.SendEmail(req.UserEmail, "Your Bus Ticket(s)", emailBody, cfg.SMTP)
			if err != nil {
				slog.Error("failed to send ticket email", slog.String("user_id", userID), slog.String("error", err.Error()))
			}
		}()

		return c.JSON(http.StatusOK, tickets)
	}
}

func buildTicketEmailBody(tickets []domain.Ticket, bus *domain.Bus, route *domain.Route) string {
	body := `<html><body><h2>AulWay Tickets</h2><hr>`

	totalPrice := 0
	ticketCount := len(tickets)

	for _, t := range tickets {
		departureDate := route.StartDate.Format("02 Jan 2006")
		departureTime := route.StartDate.Format("15:04")
		arrivalDate := route.EndDate.Format("02 Jan 2006")
		arrivalTime := route.EndDate.Format("15:04")

		body += fmt.Sprintf(`
			<h3>Ticket ID: %s</h3>
			<p><strong>Route:</strong> %s → %s<br>
			<strong>Bus Number:</strong> %s<br>
			<strong>Departure:</strong> %s at %s<br>
			<strong>Arrival:</strong> %s at %s<br>`,
			t.ID, route.Departure, route.Destination, bus.Number,
			departureDate, departureTime,
			arrivalDate, arrivalTime,
		)

		if t.QRCode != "" {
			body += fmt.Sprintf(`<img src="data:image/png;base64,%s" alt="QR Code" style="margin-top:10px;"/><br><br>`, t.QRCode)
		}

		body += "<hr>"
		totalPrice += t.Price
	}

	body += fmt.Sprintf(`
		<h3>Total Tickets: %d</h3>
		<h3>Total Price: %d₸</h3>
		</body></html>`, ticketCount, totalPrice)

	return body
}

// GetUserTicketsHandler returns a user's past or upcoming tickets.
// @Summary      Get user tickets
// @Description  Fetches a user's past or upcoming tickets based on the type parameter.
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        userId   path      string  true  "User ID"
// @Param        type     query     string  true  "Type of tickets" Enums(upcoming, past)
// @Success      200      {array}   domain.Ticket
// @Failure      400      {object}  errs.Err "Invalid type"
// @Failure      500      {object}  errs.Err "Failed to fetch tickets"
// @Router       /api/tickets/users/{userId} [get]
func GetUserTicketsHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !access.Check(c, c.Get("user_id"), "userId") {
			return c.JSON(http.StatusForbidden, errs.Err{Err: "get tickets failed", ErrDesc: "access denied"})
		}

		userID := c.Param("userId")
		ticketType := c.QueryParam("type")

		if ticketType != UpcomingTicket && ticketType != PastTicket {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid type, must be 'upcoming' or 'past'"})
		}

		now := time.Now()

		var tickets []domain.Ticket
		var err error
		if ticketType == UpcomingTicket {
			tickets, err = service.GetUpcomingTickets(c.Request().Context(), userID, now)
		} else {
			tickets, err = service.GetPastTickets(c.Request().Context(), userID, now)
		}

		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "error", ErrDesc: "Failed to fetch tickets"})
		}

		return c.JSON(http.StatusOK, tickets)
	}
}

// GetTicketDetailsHandler returns a user's ticket detail
// @Summary      Get user ticket's detail
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        userId   path      string  true  "User ID"
// @Param        ticketId   path      string  true  "Ticket ID"
// @Success      200      {object}   domain.Ticket
// @Failure      400      {object}  errs.Err
// @Failure      500      {object}  errs.Err
// @Router       /api/tickets/users/{userId}/{ticketId} [get]
func GetTicketDetailsHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !access.Check(c, c.Get("user_id"), "userId") {
			return c.JSON(http.StatusForbidden, errs.Err{Err: "get tickets failed", ErrDesc: "access denied"})
		}

		ticketId := c.Param("ticketId")

		ticket, err := service.TicketDetails(c.Request().Context(), ticketId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "error", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, ticket)
	}
}

// GetTicketsSortByHandler retrieves sorted and paginated tickets.
// @Summary Get sorted and paginated tickets
// @Description Retrieve tickets sorted by user, start date, route, price, status, or payment status with pagination.
// @Tags tickets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sort_by query string true "Sort by field (user, start_date, route, price, status, payment_status)"
// @Param order query string false "Sorting order (asc or desc)" default(asc)
// @Param page query int false "Page number (default: 1)" default(1)
// @Param page_size query int false "Number of tickets per page (default: 30)" default(30)
// @Success 200 {array} domain.Ticket "List of tickets"
// @Failure 500 {object} errs.Err "Internal server error"
// @Router /api/tickets [get]
func GetTicketsSortByHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		sortBy := c.QueryParam("sort_by")
		ord := c.QueryParam("order")

		page, pageSize := pagination.GetPageInfo(c)

		tickets, err := service.GetTicketsSortBy(ctx, sortBy, ord, page, pageSize)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "failed to get tickets", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, tickets)
	}
}
