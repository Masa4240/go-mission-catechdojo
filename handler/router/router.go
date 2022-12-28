package router

import (
	"net/http"

	"github.com/Masa4240/go-mission-catechdojo/handler"
	"github.com/Masa4240/go-mission-catechdojo/middleware"
	"github.com/Masa4240/go-mission-catechdojo/service"

	"github.com/go-chi/chi"
	//"github.com/go-chi/chi/v5"
	"github.com/jinzhu/gorm"
)

func NewRouter(userDB *gorm.DB) http.Handler {
	//func NewRouter(userDB *gorm.DB) *http.ServeMux {
	// register routes
	// mux := http.NewServeMux()
	// mux := chi.NewRouter()
	// // mux.HandleFunc("/todos", handler.NewTODOHandler(svc).ServeHTTP)
	// mux.Handle("/user/create", middleware.Recovery(userHandler))
	// mux.Handle("/user/get", middleware.Recovery(middleware.TokenValidation(userHandler)))
	// mux.Handle("/user/update", middleware.Recovery(middleware.TokenValidation(userHandler)))

	healthzHandler := handler.NewHealthzHandler()
	//mux.HandleFunc("/healthz", healthzHandler.ServeHTTP)

	r := chi.NewRouter()
	r.Use(middleware.Recovery)
	svc := service.NewUserService(userDB)
	userHandler := handler.NewUserHandler(svc)

	r.Route("/user", func(r chi.Router) {
		r.Post("/create", userHandler.CreateUser)
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.TokenValidation)
			r.Get("/get", userHandler.GetUserName)
			r.Put("/update", userHandler.UpdateUser)
		})
	})

	r.Get("/healthz", healthzHandler.ServeHTTP)

	return r
}
