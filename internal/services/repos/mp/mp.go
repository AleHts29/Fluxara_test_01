package mp

import (
	"context"
	"fluxara/internal/domain"
	"fluxara/internal/ports/repos"
	"strconv"

	"github.com/mercadopago/sdk-go/pkg/payment"
)

type MpService struct {
	db repos.DbReporerGergal
	mp repos.MpReporer
}

func NewDbService(db repos.DbReporerGergal, mp repos.MpReporer) *MpService {
	return &MpService{
		db: db,
		mp: mp,
	}
}

func (s *MpService) CreatePayment(ctx context.Context, order domain.Order) (*domain.PaymentLink, error) {
	createPay, err := s.mp.CreatePayment(ctx, order)
	if err != nil {
		return createPay, err
	}

	return createPay, err
}

func (s *MpService) ProcessWebhook(ctx context.Context, paymentID string) error {
	payment, err := s.mp.GetPayment(ctx, paymentID)
	if err != nil {
		return err
	}

	if payment.Status != "approved" {
		return nil
	}

	orderID, err := strconv.Atoi(payment.ExternalReference)
	if err != nil {
		return err
	}

	return s.db.MarkOrderPaid(ctx, orderID)
}

func (s *MpService) GetPayment(ctx context.Context, paymentID string) (*payment.Response, error) {
	return s.mp.GetPayment(ctx, paymentID)
}

func (s *MpService) MarkOrderPaid(ctx context.Context, orderID int) error {
	return s.db.MarkOrderPaid(ctx, orderID)
}
