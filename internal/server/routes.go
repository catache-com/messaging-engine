package server

import (
	"net/http"
)

type Route struct {
	Name        string
	Description string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var MessagingEngineProtectedRoutes = Routes{}

var MessagingEngineApiRoutes = Routes{}

var MessagingEngineOpenRoutes = Routes{
	// ---------- probing ----------
	Route{
		Name:    "HealthCheck",
		Method:  "GET",
		Pattern: "/healthz",
		HandlerFunc: func(writer http.ResponseWriter, r *http.Request) {
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte("OK"))
			if err != nil {
				return
			}
		},
	},

	// ---------- messaging engine ----------

	Route{
		Name:        "connect",
		Method:      "GET",
		Pattern:     "/connect",
		HandlerFunc: AcceptConnection,
	},
}
