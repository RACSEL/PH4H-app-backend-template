package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type NodeConfig struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	FhirBaseUrl         string `json:"FHIR_BASE_URL"`
	FhirMediatorBaseUrl string `json:"FHIR_MEDIATOR_BASE_URL"`
	VhlBaseUrl          string `json:"VHL_BASE_URL"`
	ICVPValidatorUrl    string `json:"ICVP_VALIDATOR_URL"`
}

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
	FhirMediatorBaseUrl  string
	LogLevel             string
	WalletEnabled        bool
	WalletUrl            string
	WalletIdentifier     string
	WalletAPIKey         string
	ICVPValidatorUrl     string
	UseMultipleNodes     bool
	Nodes                []NodeConfig
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
		FhirMediatorBaseUrl:  "http://lacpass.create.cl:3000",
		LogLevel:             "info",
		WalletEnabled:        false,
		WalletUrl:            "https://conectathon-balancer.izer.tech/",
		WalletIdentifier:     "test",
		WalletAPIKey:         "",
		ICVPValidatorUrl:     "http://lacpass.create.cl:7089",
		UseMultipleNodes:     false,
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
	if fhirMediatorBaseUrl, exists := os.LookupEnv("FHIR_MEDIATOR_BASE_URL"); exists {
		cfg.FhirMediatorBaseUrl = fhirMediatorBaseUrl
	}
	if vhlBaseUrl, exists := os.LookupEnv("VHL_BASE_URL"); exists {
		cfg.VhlBaseUrl = vhlBaseUrl
	}
	if logLevel, exists := os.LookupEnv("LOG_LEVEL"); exists {
		cfg.LogLevel = logLevel
	}

	if walletEnabled, exists := os.LookupEnv("WALLET_ENABLED"); exists {
		cfg.WalletEnabled = walletEnabled == "true"
	}
	if walletUrl, exists := os.LookupEnv("WALLET_URL"); exists {
		cfg.WalletUrl = walletUrl
	}
	if walletIdentifier, exists := os.LookupEnv("WALLET_IDENTIFIER"); exists {
		cfg.WalletIdentifier = walletIdentifier
	}

	if walletAPIKey, exists := os.LookupEnv("WALLET_API_KEY"); exists {
		cfg.WalletAPIKey = walletAPIKey
	}

	if icvpValidatorUrl, exists := os.LookupEnv("ICVP_VALIDATOR_URL"); exists {
		cfg.ICVPValidatorUrl = icvpValidatorUrl
	}

	if useMultipleNodes, exists := os.LookupEnv("USE_MULTIPLE_NODES"); exists {
		cfg.UseMultipleNodes = useMultipleNodes == "1" || useMultipleNodes == "true"
	}

	if cfg.UseMultipleNodes {
		nodesFile := "node-services.json"
		data, err := os.ReadFile(nodesFile)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", nodesFile, err)
		} else {
			var nodes []NodeConfig
			if err := json.Unmarshal(data, &nodes); err != nil {
				fmt.Printf("Error unmarshaling %s: %v\n", nodesFile, err)
			} else {
				for i := range nodes {
					if nodes[i].ID == "" {
						nodes[i].ID = nodes[i].Name
					}
				}
				cfg.Nodes = nodes
			}
		}
	}

	return cfg
}
