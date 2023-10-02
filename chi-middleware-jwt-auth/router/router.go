package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kokweikhong/go-playground/chi-middleware-jwt-auth/controller"
	custom_middleware "github.com/kokweikhong/go-playground/chi-middleware-jwt-auth/middleware"
)

func InitRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	baseController := controller.NewBaseController()

	r.Route("/", func(r chi.Router) {
		// aunthentication middleware
		r.Use(custom_middleware.AuthMiddleware)
		// chi basic auth

		r.Get("/", baseController.Index)

	})
	r.Get("/login", baseController.UserForm)
	r.Post("/login", baseController.Login)
	r.Get("/logout", baseController.Logout)

	return r
}
