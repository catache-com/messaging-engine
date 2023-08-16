package server

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"messaging-engine/internal/config"
	"net/http"
	"strings"
)

type AuthMiddleware func(next http.Handler, secretKey interface{}, authConfig config.AuthConfig) http.Handler

func NewRouter(authMiddleware AuthMiddleware) *mux.Router {

	// CORS config
	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(config.Config.AllowedOrigins, ","),
		ExposedHeaders:   []string{config.Config.Auth.CsrfHeaderName},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	r := mux.NewRouter().StrictSlash(true)
	r.Use(c.Handler)

	probingRoutes := []Route{
		{
			Name:    "HealthCheck",
			Method:  "GET",
			Pattern: "/healthz",
			HandlerFunc: func(writer http.ResponseWriter, r *http.Request) {
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte("OK"))
			},
		},
		{
			Name:    "GatewayHealthCheck",
			Method:  "GET",
			Pattern: "/gateway/catache-user-service/healthz",
			HandlerFunc: func(writer http.ResponseWriter, r *http.Request) {
				writer.WriteHeader(http.StatusOK)
				_, _ = writer.Write([]byte("OK"))
			},
		},
	}

	for _, route := range probingRoutes {
		r.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	for _, route := range MessagingEngineOpenRoutes {
		r.Methods(strings.Split(route.Method, ",")...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	authConfig := config.AuthConfig{
		JwtCookieName:  config.Config.Auth.JwtCookieName,
		CsrfCookieName: config.Config.Auth.CsrfCookieName,
		CsrfHeaderName: config.Config.Auth.CsrfHeaderName,
	}

	for _, route := range MessagingEngineProtectedRoutes {
		r.Methods(strings.Split(route.Method, ",")...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(authMiddleware(route.HandlerFunc, config.Config.Auth.AuthMiddlewareSecretKey, authConfig))
	}

	return r
}
