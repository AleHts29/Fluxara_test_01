package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (h *Handlers) GetFullData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		resumenFull, err := h.serviceDb.GetFullData(ctx)
		if err != nil {
			log.Panic("Error en GetFullData")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Printf("Esto es resumen-full: %+v", resumenFull)
		json.NewEncoder(w).Encode(resumenFull)
	}
}
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

// func (h *Handlers) GetCarrerasResumen() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := context.Background()

// 		careersResumen, err := h.serviceDb.GetCarrerasResumen(ctx)
// 		if err != nil {
// 			log.Panic("Error en GetCarrerasAll")
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)

// 		fmt.Printf("Esto es careers-resumen: %+v", careersResumen)
// 		json.NewEncoder(w).Encode(careersResumen)
// 	}
// }
