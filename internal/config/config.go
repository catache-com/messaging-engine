package config

import (
	"github.com/google/uuid"
)

var Config MessagingEngineConfig

var EngineId string

type AuthConfig struct {
	AuthMiddlewareSecretKey string `json:"auth_middleware_secret_key"`
	JwtCookieName           string `json:"jwt_cookie_name"`
	CsrfCookieName          string `json:"csrf_cookie_name"`
	CsrfHeaderName          string `json:"csrf_header_name"`
}

type MessagingEngineConfig struct {
	Host           string     `json:"host"`
	Port           int        `json:"port"`
	AllowedOrigins string     `json:"allowed_origins"` // comma seperated origins
	Auth           AuthConfig `json:"auth"`
}

func init() {
	// claim my id
	EngineId = uuid.New().String()
}
