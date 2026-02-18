package mpadapter

import (
	"fluxara/internal/config"

	configMp "github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

type MPAdapter struct {
	prefClient    preference.Client
	paymentClient payment.Client
}

func NewMPAdapter(configs *config.Config) (*MPAdapter, error) {
	token := configs.MercadoPago.Token
	cfg, err := configMp.New(token)
	if err != nil {
		return nil, err
	}

	return &MPAdapter{
		prefClient:    preference.NewClient(cfg),
		paymentClient: payment.NewClient(cfg),
	}, nil
}
