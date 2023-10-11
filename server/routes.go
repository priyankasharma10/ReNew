package server

import (
	"github.com/go-chi/chi"
)

// Update InjectRoutes to use the modified srv.register
func (srv *Server) InjectRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", srv.HealthCheck)
	r.Route("/api", func(api chi.Router) {
		api.Use(srv.MiddlewareProvider.Default()...)
		api.Route("/public", func(public chi.Router) {
			public.Post("/login", srv.loginWithEmailPassword)
		})
		api.Route("/admin", func(admin chi.Router) {
			admin.Use(srv.MiddlewareProvider.Middleware())
			admin.Post("/register", srv.register)
		})

		// Other API routes can be added here if needed
	})

	return r
}
