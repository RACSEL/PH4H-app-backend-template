package main

import (
	"encoding/json"
	customMiddleware "ips-lacpass-backend/pkg/middleware"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	userClient "ips-lacpass-backend/internal/users/client"
	userCore "ips-lacpass-backend/internal/users/core"
	userHandler "ips-lacpass-backend/internal/users/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"

	ipsClient "ips-lacpass-backend/internal/ips/client"
	ipsCore "ips-lacpass-backend/internal/ips/core"
	ipsHandler "ips-lacpass-backend/internal/ips/handler"

	vhlClient "ips-lacpass-backend/internal/vhl/client"
	vhlCore "ips-lacpass-backend/internal/vhl/core"
	vhlHandler "ips-lacpass-backend/internal/vhl/handler"

	walletClient "ips-lacpass-backend/internal/wallet/client"
	walletCore "ips-lacpass-backend/internal/wallet/core"
	walletHandler "ips-lacpass-backend/internal/wallet/handler"

	medicationClient "ips-lacpass-backend/internal/medication/client"
	medicationCore "ips-lacpass-backend/internal/medication/core"
	medicationHandler "ips-lacpass-backend/internal/medication/handler"
)

func (a *App) loadRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			// TODO: Add configuration for CORS
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, User-Agent, "+customMiddleware.NodeNameHeader)
			next.ServeHTTP(w, r)
		})
	})

	r.Use(customMiddleware.NodeNameFromHeader)

	if strings.ToLower(a.config.LogLevel) == "debug" {
		logFormat := httplog.SchemaECS.Concise(true)
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: logFormat.ReplaceAttr,
		})).With(
			slog.String("api", "lacpass"),
			slog.String("version", "0.1.0"),
			slog.String("env", "development"),
		)
		r.Use(httplog.RequestLogger(logger, &httplog.Options{
			Level:             slog.LevelDebug,
			Schema:            httplog.SchemaECS,
			RecoverPanics:     true,
			Skip:              nil,
			LogRequestHeaders: []string{"Authorization", "Content-Type", "User-Agent", customMiddleware.NodeNameHeader},
		}))
	}

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	if a.config.UseMultipleNodes {
		r.Get("/nodes", a.handleGetNodes)
	}

	r.Route("/users", a.loadUserRoutesNoAuth)

	r.Group(func(r chi.Router) {
		authMiddleware := customMiddleware.NewAuthMiddleware(
			a.config.AuthInternalUrl,
			a.config.AuthRealm,
			a.config.AuthHostName,
		)
		authMiddleware.RefreshKeySet(24 * time.Hour)
		r.Use(authMiddleware.Authenticator)

		r.Route("/ips", a.loadIpsRoute)
		r.Route("/users/auth", a.loadUserRoutesAuth)
		r.Route("/qr", a.loadVhlRoute)
		if a.config.WalletEnabled {
			r.Route("/wallet", a.loadWalletRoutes)
		}
		r.Route("/medications", a.loadMedicationRoute)
	})

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	a.router = r
}

func (a *App) loadUserRoutesNoAuth(router chi.Router) {
	r := userClient.NewClient(
		a.config.AuthInternalUrl,
		a.config.AuthRealm,
		a.config.AuthAdminClientID,
		a.config.AuthClientSecret,
		a.config.AuthEmailRedirectURI,
		a.config.AuthEmailClientID,
		a.config.AuthEmailLifespan,
	)
	s := userCore.NewService(&r)
	h := userHandler.NewHandler(&s)
	router.Post("/", h.Create)
}

func (a *App) loadUserRoutesAuth(router chi.Router) {
	r := userClient.NewClient(
		a.config.AuthInternalUrl,
		a.config.AuthRealm,
		a.config.AuthAdminClientID,
		a.config.AuthClientSecret,
		a.config.AuthEmailRedirectURI,
		a.config.AuthEmailClientID,
		a.config.AuthEmailLifespan,
	)
	s := userCore.NewService(&r)
	h := userHandler.NewHandler(&s)
	router.Put("/update", h.Update)
}

func (a *App) loadIpsRoute(router chi.Router) {
	r := ipsClient.NewClient(a.config.FhirBaseUrl, a.config.FhirMediatorBaseUrl)
	s := ipsCore.NewService(&r)

	if a.config.UseMultipleNodes {
		for _, node := range a.config.Nodes {
			nodeClient := ipsClient.NewClient(node.FhirBaseUrl, node.FhirMediatorBaseUrl)
			s.Repositories[node.ID] = &nodeClient
		}
	}

	h := ipsHandler.NewHandler(&s)
	router.Get("/", h.Get)
	router.Post("/merge", h.Merge)
	router.Get("/icvp", h.GetICVP)
}

func (a *App) loadVhlRoute(router chi.Router) {
	r := vhlClient.NewClient(a.config.VhlBaseUrl, a.config.ICVPValidatorUrl)
	s := vhlCore.NewService(&r)

	if a.config.UseMultipleNodes {
		for _, node := range a.config.Nodes {
			nodeClient := vhlClient.NewClient(node.VhlBaseUrl, node.ICVPValidatorUrl)
			s.Clients[node.ID] = &nodeClient
		}
	}

	h := vhlHandler.NewHandler(&s)
	router.Post("/", h.Create)
	router.Post("/fetch", h.Get)
	router.Post("/validate", h.Validate)
}

func (a *App) loadWalletRoutes(router chi.Router) {
	r := walletClient.NewClient(a.config.WalletUrl, a.config.WalletIdentifier, a.config.WalletAPIKey)
	s := walletCore.NewService(&r)
	h := walletHandler.NewHandler(&s)
	router.Post("/generate-link", h.GenerateWalletLink)
}

func (a *App) handleGetNodes(w http.ResponseWriter, r *http.Request) {
	nodes := make([]map[string]string, 0, len(a.config.Nodes))
	for _, node := range a.config.Nodes {
		nodes = append(nodes, map[string]string{
			"id":   node.ID,
			"name": node.Name,
		})
	}
	json.NewEncoder(w).Encode(nodes)
}

func (a *App) loadMedicationRoute(router chi.Router) {
	r := medicationClient.NewClient(a.config.FhirBaseUrl, a.config.FhirMediatorBaseUrl)
	s := medicationCore.NewService(r)
	h := medicationHandler.NewHandler(s)
	router.Get("/", h.Get)
	router.Get("/meow", h.GetMeow)
}
