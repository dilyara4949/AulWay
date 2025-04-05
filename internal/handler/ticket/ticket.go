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
	BuyTickets(ctx context.Context, userID, routeID string, paymentMethodID string, ticketAmount int, stripeKey string) ([]domain.Ticket, *domain.Bus, *domain.Route, error)
	GetUpcomingTickets(ctx context.Context, userID string, now time.Time) ([]domain.Ticket, error)
	GetPastTickets(ctx context.Context, userID string, now time.Time) ([]domain.Ticket, error)
	TicketDetails(ctx context.Context, ticketId string) (*domain.Ticket, error)
	GetTicketsSortBy(ctx context.Context, sortBy, ord string, page, pageSize int) ([]domain.Ticket, error)
	CancelTicket(ctx context.Context, userID, ticketID, stripeKey string) (*domain.Ticket, string, error)
	GetCancelledTickets(ctx context.Context, userID string) ([]domain.Ticket, error)
	GetAdminCancelledTickets(ctx context.Context, page, pageSize int) ([]domain.Ticket, error)
}

// BuyTicketHandler processes ticket purchase requests for multiple tickets.
// @Summary      Buy tickets
// @Description  Allows a user to purchase one or more tickets for a specific route using card details.
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        routeId  path      string                     true  "Route ID"
// @Param        payment_id query  string                      true  "Payment method ID - pm_card_visa"
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

		tickets, bus, route, err := s.BuyTickets(c.Request().Context(), userID, routeId, paymentId, req.Quantity, cfg.StripeKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "error", ErrDesc: err.Error()})
		}

		go func() {
			emailBody := buildTicketEmailBody(tickets, bus, route)
			err := service.SendEmailWithQR(req.UserEmail, "Your Bus Ticket(s)", tickets, cfg.SMTP, emailBody)
			if err != nil {
				slog.Error("failed to send ticket email", slog.String("user_id", userID), slog.String("error", err.Error()))
			}
		}()

		return c.JSON(http.StatusOK, tickets)
	}
}

func buildTicketEmailBody(tickets []domain.Ticket, bus *domain.Bus, route *domain.Route) string {
	if len(tickets) == 0 {
		return "<html><body><p>Билеты не найдены.</p></body></html>"
	}

	orderNumber := tickets[0].OrderNumber
	loc := time.FixedZone("Almaty", 5*60*60)
	departureDate := route.StartDate.In(loc).Format("02 Jan 2006")
	departureTime := route.StartDate.In(loc).Format("15:04")
	arrivalDate := route.EndDate.In(loc).Format("02 Jan 2006")
	arrivalTime := route.EndDate.In(loc).Format("15:04")

	body := `<html><body style="font-family: Arial, sans-serif;">`
	body += `<h2 style="color:#2d89ef;">Подтверждение покупки билетов – AulWay</h2><hr>`
	body += fmt.Sprintf(`<p><strong>Номер заказа:</strong> %s</p>`, orderNumber)
	body += fmt.Sprintf(`<p><strong>Маршрут:</strong> %s → %s<br>
<strong>Автобус №:</strong> %s<br>
<strong>Отправление:</strong> %s в %s (GMT+05 Алматы)<br>
<strong>Прибытие:</strong> %s в %s (GMT+05 Алматы)<br>
<strong>Адрес посадки:</strong> %s<br>
<strong>Адрес высадки:</strong> %s</p>`,
		route.Departure, route.Destination, bus.Number,
		departureDate, departureTime,
		arrivalDate, arrivalTime,
		route.DepartureLocation, route.DestinationLocation,
	)

	body += `<table border="1" cellpadding="10" cellspacing="0" style="border-collapse: collapse; margin-top: 20px;">
	<thead>
		<tr>
			<th>Место</th>
			<th>Цена</th>
			<th>QR-код</th>
		</tr>
	</thead>
	<tbody>`

	totalPrice := 0

	for i, t := range tickets {
		body += "<tr>"
		body += fmt.Sprintf("<td>Билет #%d</td>", i+1)
		body += fmt.Sprintf("<td>%d₸</td>", t.Price)

		if t.QRCode != "" {
			cid := fmt.Sprintf("qr%d.png", i+1)
			body += fmt.Sprintf(`<td><img src="cid:%s" alt="QR-код" style="max-width:120px;"/></td>`, cid)
		} else {
			body += "<td>Нет</td>"
		}

		body += "</tr>"
		totalPrice += t.Price
	}

	body += "</tbody></table>"

	body += fmt.Sprintf(`<p style="margin-top:20px;"><strong>Всего билетов:</strong> %d<br><strong>Общая сумма:</strong> %d₸</p>`,
		len(tickets), totalPrice)

	body += `<p style="margin-top:30px;">Спасибо за покупку!<br>Хорошей поездки с AulWay 😊</p>`
	body += `</body></html>`

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

// CancelTicketHandler allows the user to cancel their ticket
// @Summary Cancel ticket
// @Description Cancels a ticket by ID if it belongs to the user
// @Tags tickets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param ticketId path string true "Ticket ID"
// @Param email query string true "email for sending mail"
// @Success 200 {object} string "Cancellation successful"
// @Failure 400 {object} errs.Err
// @Failure 403 {object} errs.Err "Access denied"
// @Failure 500 {object} errs.Err
// @Router /api/tickets/users/{userId}/{ticketId}/cancel [put]
func CancelTicketHandler(cfg config.Config, s Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !access.Check(c, c.Get("user_id"), "userId") {
			return c.JSON(http.StatusForbidden, errs.Err{Err: "cancel failed", ErrDesc: "access denied"})
		}

		userID := c.Param("userId")
		ticketID := c.Param("ticketId")

		email := c.QueryParam("email")

		_, msg, err := s.CancelTicket(c.Request().Context(), userID, ticketID, cfg.StripeKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "cancel error", ErrDesc: err.Error()})
		}

		go func() {
			err := service.SendEmail(email, "Билет отменён", msg, cfg.SMTP)
			if err != nil {
				slog.Error("failed to send cancellation email", slog.String("user_id", userID), slog.String("error", err.Error()))
			}
		}()

		return c.JSON(http.StatusOK, map[string]string{"message": "Ticket successfully cancelled"})
	}
}

// GetCancelledTicketsHandler user's cancelled tickets
// @Summary Get cancelled tickets
// @Tags tickets
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Success 200 {array} domain.Ticket
// @Failure 403 {object} errs.Err
// @Failure 500 {object} errs.Err
// @Router /api/tickets/users/{userId}/cancelled [get]
func GetCancelledTicketsHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !access.Check(c, c.Get("user_id"), "userId") {
			return c.JSON(http.StatusForbidden, errs.Err{Err: "access denied", ErrDesc: "You are not allowed to view these tickets"})
		}

		userID := c.Param("userId")

		tickets, err := service.GetCancelledTickets(c.Request().Context(), userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "failed to fetch cancelled tickets", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, tickets)
	}
}

// GetAdminCancelledTicketsHandler cancelled tickets
// @Summary Get cancelled tickets
// @Tags tickets
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)" default(1)
// @Param page_size query int false "Number of tickets per page (default: 30)" default(30)
// @Success 200 {array} domain.Ticket
// @Failure 403 {object} errs.Err
// @Failure 500 {object} errs.Err
// @Router /api/tickets/users/cancelled [get]
func GetAdminCancelledTicketsHandler(service Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		page, pageSize := pagination.GetPageInfo(c)

		tickets, err := service.GetAdminCancelledTickets(c.Request().Context(), page, pageSize)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, errs.Err{Err: "failed to fetch cancelled tickets", ErrDesc: err.Error()})
		}

		return c.JSON(http.StatusOK, tickets)
	}
}
