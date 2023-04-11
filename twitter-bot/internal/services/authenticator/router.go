package authenticator

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"common/data/store"
	"twitter-bot/internal/services/authenticator/middlewares"
	"twitter-bot/internal/web"
	"twitter-bot/internal/web/ctx"
	"twitter-bot/internal/web/handlers"
)

// TODO: work on making this more readable
func (l *service) setupRouter() {
	l.router = chi.NewRouter()

	l.router.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		middlewares.Context(
			ctx.CtxLog(l.log),
			ctx.CtxDataProvider(store.New(l.cfg)),
			ctx.CtxOAuth(l.cfg.OAuthConfig()),
		),
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
			AllowedHeaders:   []string{web.HeaderAuthorization, web.HeaderContentType},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
	)

	l.router.Get("/healthcheck", handlers.GetHealthcheck)

	l.router.Get("/oauth/login", handlers.GetOAuth)
	l.router.Get("/oauth/callback", handlers.GetOAuthCallback)
}
