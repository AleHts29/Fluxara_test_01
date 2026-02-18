package rest

import (
	"fluxara/internal/config"
	serviceDb "fluxara/internal/services/repos/db"
	serviceDbGergal "fluxara/internal/services/repos/dbGergal"
	serviceMp "fluxara/internal/services/repos/mp"
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
	serviceDb       *serviceDb.DbService
	serviceDbGergal *serviceDbGergal.DbService
	serviceMp       *serviceMp.MpService

	handlers map[string]map[string]http.HandlerFunc
}

func NewHandlers(serviceDb *serviceDb.DbService, serviceDbGergal *serviceDbGergal.DbService, serviceMp *serviceMp.MpService) *Handlers {
	h := &Handlers{
		serviceDb:       serviceDb,
		serviceDbGergal: serviceDbGergal,
		serviceMp:       serviceMp,
		handlers:        make(map[string]map[string]http.HandlerFunc),
	}
	h.RegisterHandler("GET", "/ping", h.Ping())

	// abm
	h.RegisterHandler("GET", "/abm/resumen-full", h.GetFullData("abm"))
	h.RegisterHandler("GET", "/abm/carreras", h.GetCarrerasAll())

	// gergal
	h.RegisterHandler("GET", "/gergal/catalog-full", h.GetFullData("gergal"))
	h.RegisterHandler("GET", "/gergal/deivery-zones", h.GetDeliveryZones())
	// gergal-ordenes-pagos
	h.RegisterHandler("POST", "/gergal/orders", h.CreateOrder())
	h.RegisterHandler("POST", "/gergal/payments/webhook", h.Mpwebhook())
	h.RegisterHandler("POST", "/gergal/webhook/mercadopago", h.MercadoPagoWebhook())
	// payment-status
	h.RegisterHandler("GET", "/gergal/payments/success", h.PaymentSuccess())
	h.RegisterHandler("GET", "/gergal/payments/failure", h.PaymentFailure())
	h.RegisterHandler("GET", "/gergal/payments/pending", h.PaymentPending())

	// h.RegisterHandler("POST", "/gergal/payments/link", h.GetDeliveryZones())
	// h.RegisterHandler("POST", "/gergal/orders/previews", h.GetDeliveryZones())
	// h.RegisterHandler("GET", "/gergal/orders/{id}", h.GetDeliveryZones())

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
