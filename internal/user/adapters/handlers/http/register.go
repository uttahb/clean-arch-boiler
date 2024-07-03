package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterAuthHTTPEndpoints(r chi.Router, h *Handler) {

	authRouter := chi.NewRouter()

	authRouter.Post("/signup", h.SignUp)
	authRouter.Post("/login", h.Login)
	authRouter.Get("/ping", h.Ping)

	authenticatedRouter := chi.NewRouter()
	authenticatedRouter.Use(h.MiddlewareValidateAccessToken)
	authenticatedRouter.Get("/me", h.Me)
	authRouter.Get("/refresh-access", h.RefreshAccess)
	r.Mount("/", authenticatedRouter)
	// Mounting the new Sub Router on the main router
	r.Mount("/auth", authRouter)

}
