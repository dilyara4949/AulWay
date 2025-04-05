package service

import (
	"aulway/internal/domain"
	busRepo "aulway/internal/repository/bus"
	paymentRepo "aulway/internal/repository/payment"
	routeRepo "aulway/internal/repository/route"
	ticketRepo "aulway/internal/repository/ticket"
	"aulway/internal/utils/errs"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"image/jpeg"
	"math/rand"
	"time"
)

func NewTicketService(ticketRepo ticketRepo.Repository, paymentRepo paymentRepo.Repository, routeRepo routeRepo.Repository, processor PaymentProcessor, busRepo busRepo.Repository) *TicketService {
	return &TicketService{
		TicketRepo:       ticketRepo,
		RouteRepo:        routeRepo,
		PaymentRepo:      paymentRepo,
		PaymentProcessor: processor,
		BusRepo:          busRepo,
	}
}

type TicketService struct {
	TicketRepo       ticketRepo.Repository
	RouteRepo        routeRepo.Repository
	PaymentRepo      paymentRepo.Repository
	PaymentProcessor PaymentProcessor
	BusRepo          busRepo.Repository
}

//4242 4242 4242 4242 (Visa) – Succeeds
//4000 0000 0000 9995 (Declined)

func (s *TicketService) BuyTickets(ctx context.Context, userID, routeID, paymentMethodID string, quantity int, stripeKey string) ([]domain.Ticket, *domain.Bus, *domain.Route, error) {
	if quantity <= 0 {
		return nil, nil, nil, errors.New("quantity must be positive")
	}

	tx := s.TicketRepo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	route, err := s.RouteRepo.Get(ctx, routeID)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}
	if route.AvailableSeats < quantity {
		tx.Rollback()
		return nil, nil, nil, errs.ErrNoSeatsAvailable
	}

	totalAmount := route.Price * quantity

	success, transactionId, paymentErr := s.PaymentProcessor.ProcessPayment(ctx, userID, totalAmount, paymentMethodID, stripeKey)
	if paymentErr != nil {
		tx.Rollback()
		return nil, nil, nil, fmt.Errorf("payment failed: %w", paymentErr)
	}
	if !success {
		tx.Rollback()
		return nil, nil, nil, fmt.Errorf("payment was not successful")
	}

	paymentId, _ := uuid.NewV7()
	payment := &domain.Payment{
		ID:            paymentId.String(),
		UserID:        userID,
		Amount:        totalAmount,
		Status:        "successful",
		TransactionID: transactionId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = s.PaymentRepo.Create(ctx, tx, payment)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, fmt.Errorf("failed to create payment: %w", err)
	}

	var tickets []domain.Ticket

	for i := 0; i < quantity; i++ {
		ticketId, _ := uuid.NewV7()
		ticket := domain.Ticket{
			ID:            ticketId.String(),
			UserID:        userID,
			RouteID:       routeID,
			Price:         route.Price,
			Status:        "approved",
			PaymentStatus: "paid",
			CreatedAt:     time.Now(),
		}

		ticket.OrderNumber = generateOrderNumber()
		ticket.PaymentID = payment.ID

		qrCodePath, err := generateQRCode(&ticket)
		if err != nil {
			tx.Rollback()
			return nil, nil, nil, fmt.Errorf("failed to generate QR code: %w", err)
		}
		ticket.QRCode = qrCodePath

		err = s.TicketRepo.Create(ctx, tx, &ticket)
		if err != nil {
			tx.Rollback()
			return nil, nil, nil, fmt.Errorf("failed to create ticket: %w", err)
		}

		tickets = append(tickets, ticket)
	}

	err = s.RouteRepo.UpdateSeat(ctx, tx, map[string]interface{}{
		"available_seats": route.AvailableSeats - quantity,
	}, route.Id)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, fmt.Errorf("failed to update route seats: %w", err)
	}

	bus, err := s.BusRepo.Get(ctx, route.BusId)
	if err != nil {
		tx.Rollback()
		return nil, nil, nil, err
	}

	tx.Commit()
	return tickets, bus, route, nil
}

func generateQRCode(ticket *domain.Ticket) (string, error) {
	ticketJSON, err := json.Marshal(ticket)
	if err != nil {
		return "", err
	}

	qr, err := qrcode.New(string(ticketJSON), qrcode.Medium)
	if err != nil {
		return "", err
	}

	img := qr.Image(183)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		return "", err
	}

	qrBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return qrBase64, nil
}

func (s *TicketService) GetUpcomingTickets(ctx context.Context, userID string, now time.Time) ([]domain.Ticket, error) {
	return s.TicketRepo.GetUpcomingTickets(ctx, userID, now)
}

func (s *TicketService) GetPastTickets(ctx context.Context, userID string, now time.Time) ([]domain.Ticket, error) {
	return s.TicketRepo.GetPastTickets(ctx, userID, now)
}

func (s *TicketService) TicketDetails(ctx context.Context, ticketId string) (*domain.Ticket, error) {
	return s.TicketRepo.Get(ctx, ticketId)
}

func (s *TicketService) GetTicketsSortBy(
	ctx context.Context,
	sortBy, ord string,
	page, pageSize int,
) ([]domain.Ticket, error) {
	return s.TicketRepo.GetTicketsSortBy(ctx, sortBy, ord, page, pageSize)
}

func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d-%04d", time.Now().Unix(), rand.Intn(10000))
}

func (s *TicketService) CancelTicket(ctx context.Context, userID, ticketID, stripeKey string) (*domain.Ticket, string, error) {
	ticket, err := s.TicketRepo.Get(ctx, ticketID)
	if err != nil {
		return nil, "", fmt.Errorf("ticket not found: %w", err)
	}

	if ticket.UserID != userID {
		return nil, "", fmt.Errorf("unauthorized cancel attempt")
	}
	if ticket.Status == "cancelled" {
		return nil, "", errors.New("ticket already cancelled")
	}

	route, err := s.RouteRepo.Get(ctx, ticket.RouteID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch route: %w", err)
	}

	if time.Until(route.StartDate) < 24*time.Hour {
		return nil, "", fmt.Errorf("cancellation not allowed less than 24 hours before departure")
	}

	tx := s.TicketRepo.BeginTransaction()

	if ticket.PaymentStatus == "paid" {
		payment, err := s.PaymentRepo.GetByID(ctx, ticket.PaymentID)
		if err != nil {
			tx.Rollback()
			return nil, "", fmt.Errorf("failed to get payment: %w", err)
		}

		refundSuccess, refundErr := s.PaymentProcessor.Refund(ctx, payment.TransactionID, stripeKey, ticket.Price)
		if refundErr != nil || !refundSuccess {
			tx.Rollback()
			return nil, "", fmt.Errorf("refund failed: %w", refundErr)
		}

		//err = s.PaymentRepo.Update(ctx, tx, map[string]interface{}{
		//	"status": "refunded",
		//}, payment.ID)
		//if err != nil {
		//	tx.Rollback()
		//	return fmt.Errorf("failed to update payment status: %w", err)
		//}
	}

	err = s.TicketRepo.Update(ctx, tx, map[string]interface{}{
		"status":         "cancelled",
		"payment_status": "refunded",
	}, ticket.ID)
	if err != nil {
		tx.Rollback()
		return nil, "", fmt.Errorf("failed to cancel ticket: %w", err)
	}

	err = s.RouteRepo.IncrementSeats(ctx, tx, ticket.RouteID, 1)
	if err != nil {
		tx.Rollback()
		return nil, "", fmt.Errorf("failed to update seat count: %w", err)
	}

	msg := buildCancellationEmail(ticket, route)

	tx.Commit()
	return ticket, msg, nil
}

func buildCancellationEmail(ticket *domain.Ticket, route *domain.Route) string {
	return fmt.Sprintf(`<html><body style="font-family: Arial, sans-serif;">
		<h2 style="color:#dc3545;">Ваш билет был отменён</h2>
		<p>Номер билета: <strong>%s</strong></p>
		<p>Маршрут: <strong>%s → %s</strong></p>
		<p>Статус: <strong>Отменён</strong></p>
		<p>Спасибо, что пользуетесь AulWay</p>
	</body></html>`, ticket.OrderNumber, route.Departure, route.Destination)
}

func (s *TicketService) GetCancelledTickets(ctx context.Context, userID string) ([]domain.Ticket, error) {
	return s.TicketRepo.GetCancelledTickets(ctx, userID)
}

func (s *TicketService) GetAdminCancelledTickets(ctx context.Context, page, pageSize int) ([]domain.Ticket, error) {
	return s.TicketRepo.GetAdminCancelledTickets(ctx, page, pageSize)
}
