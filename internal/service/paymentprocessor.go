package service

import "context"

type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, userID string, amount int, paymentMethodID, stripeKey string) (bool, string, error)
	Refund(ctx context.Context, transactionID, stripeKey string, amount int) (bool, error)
}
