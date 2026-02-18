package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (h *Handlers) GetCarrerasAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		careers, err := h.serviceDb.GetCarrerasAll(ctx)
		if err != nil {
			log.Panic("Error en GetCarrerasAll")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Printf("Esto es careers: %+v", careers)
		json.NewEncoder(w).Encode(careers)
	}
}
