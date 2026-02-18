package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
