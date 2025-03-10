package stripe

import (
	"github.com/stretchr/testify/mock"
)

type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) ProcessPayment(amount int, currency string) (bool, string) {
	args := m.Called(amount, currency)
	return args.Bool(0), args.String(1)
}
