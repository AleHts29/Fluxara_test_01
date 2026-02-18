package repos

import (
	"context"
	"fluxara/internal/domain"

	"github.com/mercadopago/sdk-go/pkg/payment"
)

type MpReporer interface {
	CreatePayment(ctx context.Context, order domain.Order) (*domain.PaymentLink, error)
	GetPayment(ctx context.Context, paymentID string) (*payment.Response, error)
}
