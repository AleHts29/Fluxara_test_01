package rest

import "net/http"

func (h *Handlers) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
