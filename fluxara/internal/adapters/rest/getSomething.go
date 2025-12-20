package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handlers) GetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		vars := mux.Vars(r)
		id := vars["id"] // id

		product, err := h.serviceDb.GetProduct(ctx, id)
		if err != nil {
			log.Panic("ERror en GetProduct")
		}

		fmt.Printf("Esto es product %+v \n", product)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(product)
	}
}
