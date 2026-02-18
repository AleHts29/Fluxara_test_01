package mpadapter

import (
	"context"
	"fluxara/internal/domain"
	"fmt"
	"strconv"

	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

func (m *MPAdapter) CreatePayment(ctx context.Context, order domain.Order) (*domain.PaymentLink, error) {
	fmt.Printf("Ingreso a Create Payment -.-------------\n")

	req := preference.Request{
		ExternalReference: strconv.Itoa(order.ID),
		Items: []preference.ItemRequest{
			{
				Title:     fmt.Sprintf("Pedido #%d", order.ID),
				Quantity:  1,
				UnitPrice: float64(order.Total),
			},
		},
		AutoReturn: "approved",
		BackURLs: &preference.BackURLsRequest{
			Success: "http://localhost:8099/gergal/payments/success",
			Failure: "http://localhost:8099/gergal/payments/failure",
			Pending: "http://localhost:8099/gergal/payments/pending",
		},
	}

	pref, err := m.prefClient.Create(ctx, req)
	if err != nil {
		fmt.Printf("Error en .Create -.-------------\n")
		return nil, err
	}

	fmt.Printf("Retorno create Link-.-------------\n")
	return &domain.PaymentLink{
		OrderID: order.ID,
		URL:     pref.InitPoint,
	}, nil
}

func (m *MPAdapter) GetPayment(ctx context.Context, paymentID string) (*payment.Response, error) {
	id, err := strconv.ParseInt(paymentID, 10, 64)
	if err != nil {
		return nil, err
	}

	return m.paymentClient.Get(ctx, int(id))
}
