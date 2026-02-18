package rest

import (
	"context"
	"encoding/json"
	"fluxara/internal/domain"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func (h *Handlers) GetDeliveryZones() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		deliveryZones, err := h.serviceDbGergal.GetDeliveryZones(ctx)
		if err != nil {
			log.Panic("Error en GetDeliveryZones")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Printf("Esto es delivery-zones: %+v", deliveryZones)
		json.NewEncoder(w).Encode(deliveryZones)

	}
}

func (h *Handlers) CreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req domain.CreateOrderRequest
		fmt.Printf("Leyendo -.-------------")
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "json inválido", http.StatusBadRequest)
			return
		}

		fmt.Printf("Leído -.-------------")
		// 1) Crear orden en DB
		order, err := h.serviceDbGergal.CreateOrder(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Printf("Orden creada -.-------------\n")
		// 2) Generar link de pago
		payment, err := h.serviceMp.CreatePayment(ctx, *order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Link creados -.-------------")

		resp := domain.CreateOrderResponse{
			Order:   order,
			Payment: payment,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

func (h *Handlers) Mpwebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var notif struct {
			Type string `json:"type"`
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}

		if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		if notif.Type != "payment" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if err := h.serviceMp.ProcessWebhook(ctx, notif.Data.ID); err != nil {
			log.Println("Webhook error:", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}
func (h *Handlers) MercadoPagoWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payload struct {
			Type string `json:"type"`
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		// MP manda muchos tipos de eventos, solo nos interesa payment
		if payload.Type != "payment" {
			w.WriteHeader(http.StatusOK)
			return
		}

		ctx := r.Context()

		payment, err := h.serviceMp.GetPayment(ctx, payload.Data.ID)
		if err != nil {
			http.Error(w, "error getting payment", http.StatusInternalServerError)
			return
		}

		if payment.Status == "approved" {
			orderID, err := strconv.Atoi(payment.ExternalReference)
			if err != nil {
				http.Error(w, "invalid external_reference", http.StatusBadRequest)
				return
			}

			if err := h.serviceMp.MarkOrderPaid(ctx, orderID); err != nil {
				http.Error(w, "error updating order", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}

}

func (h *Handlers) PaymentSuccess() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Pago aprobado. Gracias por tu compra."))
	}
}

func (h *Handlers) PaymentFailure() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("El pago fue rechazado. Intente nuevamente."))
	}
}

func (h *Handlers) PaymentPending() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("El pago se encuentra pendiente de confirmación."))
	}
}
