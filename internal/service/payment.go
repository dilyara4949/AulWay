package service

import (
	"context"
	"fmt"
	"github.com/stripe/stripe-go/v76"
	_ "github.com/stripe/stripe-go/v76/charge" // ⚠️ Required to avoid runtime issues
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"
)

type StripeProcessor struct{}

func NewStripeProcessor() PaymentProcessor {
	return &StripeProcessor{}
}

func (s *StripeProcessor) Refund(ctx context.Context, transactionID, stripeKey string, amount int) (bool, error) {
	stripe.Key = stripeKey

	refundParams := &stripe.RefundParams{
		PaymentIntent: stripe.String(transactionID),
		Amount:        stripe.Int64(int64(amount * 100)),
	}

	_, err := refund.New(refundParams)
	if err != nil {
		return false, fmt.Errorf("stripe refund error: %w", err)
	}

	return true, nil
}

func (s *StripeProcessor) ProcessPayment(ctx context.Context, userID string, amount int, paymentMethodID string, stripeKey string) (bool, string, error) {
	stripe.Key = stripeKey

	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(int64(amount * 100)),
		Currency:      stripe.String(string(stripe.CurrencyKZT)),
		PaymentMethod: stripe.String(paymentMethodID),
		Confirm:       stripe.Bool(true),

		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled:        stripe.Bool(true),
			AllowRedirects: stripe.String("never"),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return false, "", fmt.Errorf("stripe payment error: %w", err)
	}

	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		return false, "", fmt.Errorf("payment not successful: %s", pi.Status)
	}

	return true, pi.ID, nil
}
