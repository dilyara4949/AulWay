package service

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

func NewFPaymentProcessor() PaymentProcessor {
	return &fakePaymentProcessor{}
}

type fakePaymentProcessor struct{}

func (p *fakePaymentProcessor) ProcessPayment(ctx context.Context, userID string, amount int, paymentMethodID string) (bool, error) {
	time.Sleep(1 * time.Second)

	if rand.Intn(100) < 90 {
		return true, nil
	}

	return false, errors.New("payment failure")
}

func (p *fakePaymentProcessor) Refund(ctx context.Context, transactionID string, amount int) (bool, error) {
	time.Sleep(500 * time.Millisecond)

	if rand.Intn(100) < 95 {
		return true, nil
	}

	return false, errors.New("refund failure")
}
