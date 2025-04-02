package service

import (
	"context"
	"fmt"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, userID string, amount int, paymentMethodID string) (bool, error)
	Refund(ctx context.Context, transactionID string, amount int) (bool, error)
}

func NewPaymentProcessor() PaymentProcessor {
	return &stripePaymentProcessor{}
}

type stripePaymentProcessor struct{}

func (s *stripePaymentProcessor) Refund(ctx context.Context, transactionID string, amount int) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s *stripePaymentProcessor) ProcessPayment(ctx context.Context, userID string, amount int, paymentMethodID string) (bool, error) {
	stripe.Key = "stripeKey"

	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(int64(amount * 100)), // Convert to cents
		Currency:      stripe.String("usd"),
		PaymentMethod: stripe.String(paymentMethodID),
		Confirm:       stripe.Bool(true),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return false, fmt.Errorf("stripe payment failed: %w", err)
	}

	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		return false, fmt.Errorf("payment not successful")
	}

	return true, nil
}
