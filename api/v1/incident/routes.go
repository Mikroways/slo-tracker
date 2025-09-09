package incident

import (
	"net/http"

	"slo-tracker/api"
	"slo-tracker/middleware"
	appstore "slo-tracker/store"

	"github.com/go-chi/chi"
)

// store holds shared store conn from the api
var store *appstore.Conn

// Init initializes all the v1 routes
func Init(r chi.Router) {

	store = api.Store

	r.Method(http.MethodGet, "/", api.Handler(getAllIncidentsHandler))
	r.Method(http.MethodPost, "/", api.Handler(createIncidentHandler))
	r.Route("/webhook", webhookSubRoutes)
	r.With(middleware.IncidentRequired).
		Route("/{incidentID:[0-9]+}", incidentIDSubRoutes)
}

// ROUTE: {host}/api/v1/incident/falsepositive/:incidentID
func InitFalsePositive(r chi.Router) {
	r.With(middleware.IncidentRequired).
		Method(http.MethodPatch, "/{incidentID:[0-9]+}", api.Handler(updateFalsePositive))
}

// ROUTE: {host}/api/v1/incident/:incidentID/*
func incidentIDSubRoutes(r chi.Router) {
	r.Method(http.MethodGet, "/", api.Handler(getIncidentHandler))
	r.Method(http.MethodPatch, "/", api.Handler(updateIncidentHandler))
}

func webhookSubRoutes(r chi.Router) {
	r.Method(http.MethodPost, "/newrelic", api.Handler(createNewrelicIncidentHandler))
	r.Method(http.MethodPost, "/prometheus", api.Handler(createPromIncidentHandler))
	r.Method(http.MethodPost, "/pingdom", api.Handler(createPingdomIncidentHandler))
	r.Method(http.MethodPost, "/datadog", api.Handler(createDatadogIncidentHandler))
	r.Method(http.MethodPost, "/grafana", api.Handler(createGrafanaIncidentHandler))
}
