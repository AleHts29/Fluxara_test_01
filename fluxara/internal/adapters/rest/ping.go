package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Handlers) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("Esta bien ping\n")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(http.StatusOK)
	}
}
