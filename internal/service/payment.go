package service

import (
	"context"
	"fmt"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

const stripeKey = "sk_test_51R0O8UQr8wSBOpQomnUQjRXcuxTXA5XVJOHeCxfY8caeGpGG3zmLk9zlo2lMAN5CI3vUy3dy2zSsE1MfhiHBnlre00EUHCkRi1"

type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, userID string, amount int, paymentMethodID string) (bool, error)
}

func NewPaymentProcessor() PaymentProcessor {
	return &stripePaymentProcessor{}
}

type stripePaymentProcessor struct{}

func (s *stripePaymentProcessor) ProcessPayment(ctx context.Context, userID string, amount int, paymentMethodID string) (bool, error) {
	stripe.Key = stripeKey

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

//func (s *stripePaymentProcessor) ProcessPayment(ctx context.Context, userID string, amount int, card domain.CardDetails) (bool, error) {
//	stripe.Key = stripeKey
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
