package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handlers) GetCarrerasAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		careers, err := h.serviceDb.GetCarrerasAll(ctx)
		if err != nil {
			log.Panic("Error en GetCarrerasAll")
		}

		fmt.Printf("Esto es carreras %+v \n", careers)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(careers)
	}
}

func (h *Handlers) GetCarrerasResumen() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		product, err := h.serviceDb.GetCarrerasResumen(ctx)
		if err != nil {
			log.Panic("Error en GetCarrerasAll")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(product)
	}
}

func (h *Handlers) GetCarrerasByName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		vars := mux.Vars(r)
		name := vars["name"] // name

		career, err := h.serviceDb.GetCarrerasByName(ctx, name)
		if err != nil {
			log.Panic("Error en GetCarrerasByName")
		}

		fmt.Printf("Esto es career %+v \n", career)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(career)
	}
}
