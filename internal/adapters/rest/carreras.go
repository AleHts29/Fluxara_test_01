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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Printf("Esto es careers: %+v", careers)
		json.NewEncoder(w).Encode(careers)
	}
}

func (h *Handlers) GetCarrerasResumen() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		careersResumen, err := h.serviceDb.GetCarrerasResumen(ctx)
		if err != nil {
			log.Panic("Error en GetCarrerasAll")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Printf("Esto es careers-resumen: %+v", careersResumen)
		json.NewEncoder(w).Encode(careersResumen)
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

		fmt.Printf("Esto es career: %+v", career)
		json.NewEncoder(w).Encode(career)
	}
}
