package main

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort           uint16
	AuthHostName         string
	AuthInternalUrl      string
	AuthRealm            string
	AuthAdminClientID    string
	AuthClientSecret     string
	AuthEmailRedirectURI string
	AuthEmailLifespan    int
	AuthEmailClientID    string
	FhirBaseUrl          string
	VhlBaseUrl           string
	APISwagger           bool
	LogLevel             string
}

func LoadConfig() Config {
	cfg := Config{
		ServerPort:           3000,
		AuthHostName:         "http://localhost:9083",
		AuthInternalUrl:      "http://localhost:9083",
		AuthRealm:            "lacpass",
		AuthAdminClientID:    "admin-cli",
		AuthClientSecret:     "bbU4vnqhqe2AJ32XpdQVRVqfRMA82Hnu",
		AuthEmailRedirectURI: "ph4happ://open/validated-email",
		AuthEmailLifespan:    3600,
		AuthEmailClientID:    "app",
		FhirBaseUrl:          "http://lacpass.create.cl:8080",
		VhlBaseUrl:           "http://lacpass.create.cl:8182",
		APISwagger:           false,
		LogLevel:             "info",
	}

	if serverPort, exists := os.LookupEnv("API_PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}
	if authUrl, exists := os.LookupEnv("AUTH_INTERNAL_URL"); exists {
		cfg.AuthInternalUrl = authUrl
	}
	if authHostname, exists := os.LookupEnv("AUTH_HOSTNAME"); exists {
		cfg.AuthHostName = authHostname
	}
	if authRealm, exists := os.LookupEnv("AUTH_REALM"); exists {
		cfg.AuthRealm = authRealm
	}
	if authEmailClientID, exists := os.LookupEnv("AUTH_EMAIL_CLIENT_ID"); exists {
		cfg.AuthEmailClientID = authEmailClientID
	}
	if authClientSecret, exists := os.LookupEnv("AUTH_CLIENT_SECRET"); exists {
		cfg.AuthClientSecret = authClientSecret
	}
	if authEmailRedirectURI, exists := os.LookupEnv("AUTH_EMAIL_REDIRECT_URI"); exists {
		cfg.AuthEmailRedirectURI = authEmailRedirectURI
	}
	if authEmailLifespan, exists := os.LookupEnv("AUTH_EMAIL_LIFESPAN"); exists {
		if lifespan, err := strconv.Atoi(authEmailLifespan); err == nil {
			cfg.AuthEmailLifespan = lifespan
		}
	}
	if authEmailClientID, exists := os.LookupEnv("AUTH_EMAIL_CLIENT_ID"); exists {
		cfg.AuthEmailClientID = authEmailClientID
	}
	if authEmailClientID, exists := os.LookupEnv("AUTH_EMAIL_CLIENT_ID"); exists {
		cfg.AuthEmailClientID = authEmailClientID
	}
	if fhirBaseUrl, exists := os.LookupEnv("FHIR_BASE_URL"); exists {
		cfg.FhirBaseUrl = fhirBaseUrl
	}
	if vhlBaseUrl, exists := os.LookupEnv("VHL_BASE_URL"); exists {
		cfg.VhlBaseUrl = vhlBaseUrl
	}
	if apiSwagger, exists := os.LookupEnv("API_SWAGGER"); exists {
		cfg.APISwagger = apiSwagger == "true"
	}
	if logLevel, exists := os.LookupEnv("LOG_LEVEL"); exists {
		cfg.LogLevel = logLevel
	}

	return cfg
}
