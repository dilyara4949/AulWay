package service

import (
	"aulway/internal/domain"
	paymentRepo "aulway/internal/repository/payment"
	routeRepo "aulway/internal/repository/route"
	ticketRepo "aulway/internal/repository/ticket"
	"aulway/internal/utils/errs"
	"context"
	"fmt"
	"github.com/google/uuid"
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

func (s *TicketService) BuyTicket(ctx context.Context, userID, routeID string, paymentMethodID string) (*domain.Ticket, error) {
	tx := s.TicketRepo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Fetch route and check seat availability
	route, err := s.RouteRepo.Get(ctx, routeID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if route.AvailableSeats <= 0 {
		tx.Rollback()
		return nil, errs.ErrNoSeatsAvailable
	}

	// Create ticket
	ticketId, _ := uuid.NewV7()
	ticket := &domain.Ticket{
		ID:            ticketId.String(),
		UserID:        userID,
		RouteID:       routeID,
		Price:         route.Price,
		Status:        "awaiting",
		PaymentStatus: "pending",
		CreatedAt:     time.Now(),
	}

	err = s.TicketRepo.Create(ctx, tx, ticket)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Process payment **before** committing ticket
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

	err = s.RouteRepo.UpdateSeat(ctx, tx, map[string]interface{}{
		"available_seats": route.AvailableSeats - 1,
	}, route.Id)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update route seats: %w", err)
	}

	tx.Commit()
	return ticket, nil
}

//func (s *TicketService) BuyTicket(ctx context.Context, userID, routeID string, paymentMethodID string) (*domain.Ticket, error) {
//	tx := s.TicketRepo.BeginTransaction()
//	defer func() {
//		if r := recover(); r != nil {
//			tx.Rollback()
//		}
//	}()
//
//	route, err := s.RouteRepo.Get(ctx, routeID)
//	if err != nil {
//		tx.Rollback()
//		return nil, err
//	}
//	if route.AvailableSeats <= 0 {
//		tx.Rollback()
//		return nil, errs.ErrNoSeatsAvailable
//	}
//
//	ticketId, _ := uuid.NewV7()
//	ticket := &domain.Ticket{
//		ID:            ticketId.String(),
//		UserID:        userID,
//		RouteID:       routeID,
//		Price:         route.Price,
//		Status:        "awaiting",
//		PaymentStatus: "pending",
//		CreatedAt:     time.Now(),
//	}
//
//	err = s.TicketRepo.Create(ctx, tx, ticket)
//	if err != nil {
//		tx.Rollback()
//		return nil, fmt.Errorf("failed to create ticket: %w", err)
//	}
//
//	tx.Commit()
//
//	transactionId, _ := uuid.NewV7()
//	success, stripeErr := s.PaymentProcessor.ProcessPayment(ctx, userID, ticket.Price, paymentMethodID)
//	if stripeErr != nil {
//		tx.Rollback()
//		return nil, fmt.Errorf("payment failed: %w", stripeErr)
//	}
//
//	if !success {
//		tx.Rollback()
//		return nil, fmt.Errorf("payment was not successful")
//	}
//
//	paymentId, _ := uuid.NewV7()
//	payment := &domain.Payment{
//		ID:            paymentId.String(),
//		UserID:        userID,
//		TicketID:      ticket.ID,
//		Amount:        ticket.Price,
//		Status:        "successful",
//		TransactionID: transactionId.String(),
//		CreatedAt:     time.Now(),
//		UpdatedAt:     time.Now(),
//	}
//
//	err = s.PaymentRepo.Create(ctx, payment)
//	if err != nil {
//		tx.Rollback()
//		return nil, fmt.Errorf("failed to create payment: %w", err)
//	}
//
//	err = s.TicketRepo.Update(ctx, tx, map[string]interface{}{
//		"status":         "approved",
//		"payment_status": "paid",
//	}, ticket.ID)
//	if err != nil {
//		tx.Rollback()
//		return nil, fmt.Errorf("failed to update ticket: %w", err)
//	}
//
//	err = s.RouteRepo.Update(ctx, map[string]interface{}{
//		"available_seats": route.AvailableSeats - 1,
//	}, route.Id)
//	if err != nil {
//		tx.Rollback()
//		return nil, fmt.Errorf("failed to update route seats: %w", err)
//	}
//
//	tx.Commit()
//	return ticket, nil
//}

//func (s *TicketService) ProcessPayment(ctx context.Context, userID string, amount int, card domain.CardDetails) (bool, error) {
//	stripe.Key = "sk_test_your_secret_key"
//
//	pmParams := &stripe.PaymentMethodParams{
//		Type: stripe.String("card"),
//		Card: &stripe.PaymentMethodCardParams{
//			Number:   stripe.String(card.Number),
//			ExpMonth: stripe.Int64(card.ExpMonth),
//			ExpYear:  stripe.Int64(card.ExpYear),
//			CVC:      stripe.String(card.CVC),
//		},
//	}
//
//	paymentMethod, err := paymentmethod.New(pmParams)
//	if err != nil {
//		return false, fmt.Errorf("failed to create payment method: %w", err)
//	}
//
//	params := &stripe.PaymentIntentParams{
//		Amount:        stripe.Int64(int64(amount * 100)),
//		Currency:      stripe.String("usd"),
//		PaymentMethod: stripe.String(paymentMethod.ID),
//		Confirm:       stripe.Bool(true),
//	}
//
//	pi, err := paymentintent.New(params)
//	if err != nil {
//		return false, fmt.Errorf("stripe payment failed: %w", err)
//	}
//
//	if pi.Status != stripe.PaymentIntentStatusSucceeded {
//		return false, fmt.Errorf("payment not successful")
//	}
//
//	return true, nil
//}
