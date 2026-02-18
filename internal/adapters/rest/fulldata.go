package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (h *Handlers) GetFullData(entity string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if entity == "abm" {
			resumenFull, err := h.serviceDb.GetFullData(ctx)
			if err != nil {
				log.Panic("Error en GetFullData")
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Printf("Esto es resumen-full: %+v", resumenFull)
			json.NewEncoder(w).Encode(resumenFull)
		}
		if entity == "gergal" {
			catalogFull, err := h.serviceDbGergal.GetCatalog(ctx)
			if err != nil {
				log.Panic("Error en GetFullData")
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Printf("Esto es resumen-full: %+v", catalogFull)
			json.NewEncoder(w).Encode(catalogFull)
		}

	}
}
