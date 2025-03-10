package service

import (
	"aulway/internal/domain"
	paymentRepo "aulway/internal/repository/payment"
	routeRepo "aulway/internal/repository/route"
	ticketRepo "aulway/internal/repository/ticket"
	"aulway/internal/utils/errs"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"os"
	"path/filepath"
	"time"
)

func NewTicketService(ticketRepo ticketRepo.Repository, paymentRepo paymentRepo.Repository, routeRepo routeRepo.Repository, processor PaymentProcessor) *TicketService {
	return &TicketService{
		TicketRepo:       ticketRepo,
		RouteRepo:        routeRepo,
		PaymentRepo:      paymentRepo,
		PaymentProcessor: processor,
	}
}

type TicketService struct {
	TicketRepo       ticketRepo.Repository
	RouteRepo        routeRepo.Repository
	PaymentRepo      paymentRepo.Repository
	PaymentProcessor PaymentProcessor
}

//4242 4242 4242 4242 (Visa) – Succeeds
//4000 0000 0000 9995 (Declined)

func (s *TicketService) BuyTickets(ctx context.Context, userID, routeID, paymentMethodID string, quantity int) ([]domain.Ticket, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	tx := s.TicketRepo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	route, err := s.RouteRepo.Get(ctx, routeID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if route.AvailableSeats < quantity {
		tx.Rollback()
		return nil, errs.ErrNoSeatsAvailable
	}

	var tickets []domain.Ticket

	for i := 0; i < quantity; i++ {
		// Create ticket
		ticketId, _ := uuid.NewV7()
		ticket := domain.Ticket{
			ID:            ticketId.String(),
			UserID:        userID,
			RouteID:       routeID,
			Price:         route.Price,
			Status:        "awaiting",
			PaymentStatus: "pending",
			CreatedAt:     time.Now(),
		}

		err = s.TicketRepo.Create(ctx, tx, &ticket)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create ticket: %w", err)
		}

		qrCodePath, err := generateQRCode(&ticket)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to generate QR code: %w", err)
		}
		ticket.QRCode = qrCodePath

		transactionId, _ := uuid.NewV7()
		success, stripeErr := s.PaymentProcessor.ProcessPayment(ctx, userID, ticket.Price, paymentMethodID)
		if stripeErr != nil {
			tx.Rollback()
			return nil, fmt.Errorf("payment failed: %w", stripeErr)
		}
		if !success {
			tx.Rollback()
			return nil, fmt.Errorf("payment was not successful")
		}

		// Create payment record
		paymentId, _ := uuid.NewV7()
		payment := &domain.Payment{
			ID:            paymentId.String(),
			UserID:        userID,
			TicketID:      ticket.ID,
			Amount:        ticket.Price,
			Status:        "successful",
			TransactionID: transactionId.String(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		err = s.PaymentRepo.Create(ctx, tx, payment)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create payment: %w", err)
		}

		err = s.TicketRepo.Update(ctx, tx, map[string]interface{}{
			"status":         "approved",
			"payment_status": "paid",
		}, ticket.ID)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update ticket: %w", err)
		}

		tickets = append(tickets, ticket)
	}

	err = s.RouteRepo.UpdateSeat(ctx, tx, map[string]interface{}{
		"available_seats": route.AvailableSeats - quantity,
	}, route.Id)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update route seats: %w", err)
	}

	tx.Commit()
	return tickets, nil
}

func generateQRCode(ticket *domain.Ticket) (string, error) {
	qrDir := "qrcodes"
	qrPath := filepath.Join(qrDir, fmt.Sprintf("%s.png", ticket.ID))

	// Ensure directory exists
	if err := os.MkdirAll(qrDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create QR code directory: %w", err)
	}

	// QR code content - URL for ticket PDF download
	qrContent := fmt.Sprintf("https://localhos:8080/api/tickets/%s/download", ticket.ID)

	// Generate and save the QR code
	err := qrcode.WriteFile(qrContent, qrcode.Medium, 256, qrPath)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qrPath, nil
}
