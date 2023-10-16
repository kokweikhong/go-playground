package router

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kokweikhong/go-playground/auth-nextjs-golang/backend/internal/controller"
	myMiddleWare "github.com/kokweikhong/go-playground/auth-nextjs-golang/backend/internal/middleware"
)

func Init() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		// MaxAge:           300,
	}))

	userController := controller.NewUserController()

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			data := map[string]interface{}{
				"data": "hello world",
			}
			json.NewEncoder(w).Encode(data)
		})

		r.Post("/login", userController.Login)

		r.Route("/users", func(r chi.Router) {
			r.Use(myMiddleWare.AuthMiddleware)
			r.Get("/", userController.ListUsers)
			r.Post("/", userController.CreateUser)           // POST /users
			r.Get("/{userId}", userController.GetUser)       // GET /users/123
			r.Put("/{userId}", userController.UpdateUser)    // PUT /users/123
			r.Delete("/{userId}", userController.DeleteUser) // DELETE /users/123
		})
	})
	return r
}
