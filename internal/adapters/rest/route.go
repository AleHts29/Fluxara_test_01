package rest

import (
	"fluxara/internal/config"
	serviceDb "fluxara/internal/services/repos/db"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Host string
	Port string
}

type Handlers struct {
	serviceDb *serviceDb.DbService
	handlers  map[string]map[string]http.HandlerFunc
}

func NewHandlers(serviceDb *serviceDb.DbService) *Handlers {
	h := &Handlers{
		serviceDb: serviceDb,
		handlers:  make(map[string]map[string]http.HandlerFunc),
	}

	// registra handlers necesarios - los que se necesiten
	// h.RegisterHandler("GET", "/products", h.GetProductsAll())
	// h.RegisterHandler("GET", "/products/{id}", h.GetProduct())

	// artes
	h.RegisterHandler("GET", "/resumen-full", h.GetFullData())
	h.RegisterHandler("GET", "/carreras", h.GetCarrerasAll())
	h.RegisterHandler("GET", "/ping", h.Ping())

	return h
}
func (h *Handlers) RegisterHandler(method string, endopoint string, handler http.HandlerFunc) {
	if _, exists := h.handlers[method]; !exists {
		h.handlers[method] = make(map[string]http.HandlerFunc)
	}

	// agregado de handler por método y ruta
	h.handlers[method][endopoint] = handler
}

func NewRouter(configs *config.Config, handlers *Handlers) {
	s := Server{
		Host: configs.Server.Host,
		Port: configs.Server.Port,
	}

	router := mux.NewRouter()
	router.Use()

	// asigna handlers dinámicamente según correspondan hallan sido creados
	for method, routes := range handlers.handlers {
		for endpoint, handler := range routes {
			router.HandleFunc(endpoint, handler).Methods(method)
			log.Printf("HandleFunc added, method %s, endpoint %s", method, endpoint)
		}
	}

	http.Handle("/", router)

	// arranca el listener del server
	currentHost := fmt.Sprintf("%s:%s", s.Host, s.Port)
	// currentHost := "0.0.0.0:8099"
	log.Printf("Starting server %s\n", currentHost)
	if err := http.ListenAndServe(currentHost, router); err != nil {
		log.Panic("Error en ListenAdnServe")
	}
}
