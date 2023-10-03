package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kokweikhong/go-playground/htmx/controller"
)

func InitRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// static file server
	fs := http.FileServer(http.Dir("./public"))
	r.Handle("/public/*", http.StripPrefix("/public", fs))

	baseController := controller.NewBaseController()

	r.Get("/", baseController.Index)
	// r.Get("/expenses-category/{id}", baseController.GetExpensesCategory)
	r.Get("/expenses-category/{formType}", baseController.NewExpensesCategory)
	r.Route("/expenses-category", func(r chi.Router) {
		r.Post("/", baseController.CreateExpensesCategory)
		r.Put("/", baseController.UpdateExpensesCategory)
		r.Delete("/{id}", baseController.DeleteExpensesCategory)
	})

	return r
}
