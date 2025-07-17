package main

import (
	"ips-lacpass-backend/internal/core"
	"ips-lacpass-backend/internal/repository/fhir"
	"ips-lacpass-backend/internal/repository/keycloak"
	"ips-lacpass-backend/internal/repository/vhl"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	_ "ips-lacpass-backend/internal/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	httpSwagger "github.com/swaggo/http-swagger"

	"ips-lacpass-backend/internal/handler"
	customMiddleware "ips-lacpass-backend/internal/middleware"
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
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, User-Agent")
			next.ServeHTTP(w, r)
		})
	})

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
			LogRequestHeaders: []string{"Authorization", "Content-Type", "User-Agent"},
		}))
	}

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

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
	})

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	if a.config.APISwagger {
		r.Get("/swagger/*", httpSwagger.Handler())
	}

	a.router = r
}

func (a *App) loadUserRoutesNoAuth(router chi.Router) {
	r := keycloak.NewKeycloakClient(
		a.config.AuthInternalUrl,
		a.config.AuthRealm,
		a.config.AuthAdminClientID,
		a.config.AuthClientSecret,
		a.config.AuthEmailRedirectURI,
		a.config.AuthEmailClientID,
		a.config.AuthEmailLifespan,
	)
	s := core.NewUserService(r)
	h := handler.NewUserHandler(s)
	router.Post("/", h.Create)
}

func (a *App) loadUserRoutesAuth(router chi.Router) {
	r := keycloak.NewKeycloakClient(
		a.config.AuthInternalUrl,
		a.config.AuthRealm,
		a.config.AuthAdminClientID,
		a.config.AuthClientSecret,
		a.config.AuthEmailRedirectURI,
		a.config.AuthEmailClientID,
		a.config.AuthEmailLifespan,
	)
	s := core.NewUserService(r)
	h := handler.NewUserHandler(s)
	router.Put("/update", h.Update)
}

func (a *App) loadIpsRoute(router chi.Router) {
	r := fhir.FhirRepository{
		Client:  &http.Client{},
		BaseURL: a.config.FhirBaseUrl,
	}
	s := core.NewFhirService(r)
	h := handler.NewIpsHandler(s)
	router.Get("/", h.Get)
}

func (a *App) loadVhlRoute(router chi.Router) {
	r := vhl.VhlRepository{
		Client:  &http.Client{},
		BaseURL: a.config.VhlBaseUrl,
	}
	s := core.NewVhlService(r)
	h := handler.NewVhlHandler(s)
	router.Post("/", h.Create)
}
