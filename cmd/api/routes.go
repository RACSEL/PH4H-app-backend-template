package main

import (
	customMiddleware "ips-lacpass-backend/pkg/middleware"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	httpSwagger "github.com/swaggo/http-swagger"

	userClient "ips-lacpass-backend/internal/users/client"
	userCore "ips-lacpass-backend/internal/users/core"
	userHandler "ips-lacpass-backend/internal/users/handler"

	ipsClient "ips-lacpass-backend/internal/ips/client"
	ipsCore "ips-lacpass-backend/internal/ips/core"
	ipsHandler "ips-lacpass-backend/internal/ips/handler"

	vhlClient "ips-lacpass-backend/internal/vhl/client"
	vhlCore "ips-lacpass-backend/internal/vhl/core"
	vhlHandler "ips-lacpass-backend/internal/vhl/handler"
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
	r := ipsClient.NewClient(a.config.FhirBaseUrl)
	s := ipsCore.NewService(&r)
	h := ipsHandler.NewHandler(&s)
	router.Get("/", h.Get)
	router.Post("/merge", h.Merge)
}

func (a *App) loadVhlRoute(router chi.Router) {
	r := vhlClient.NewClient(a.config.VhlBaseUrl)
	s := vhlCore.NewService(&r)
	h := vhlHandler.NewHandler(&s)
	router.Post("/", h.Create)
	router.Post("/fetch", h.Get)
}
