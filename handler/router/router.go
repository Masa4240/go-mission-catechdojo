package router

import (
	"database/sql"
	"net/http"

	"github.com/Masa4240/go-mission-catechdojo/handler"
	"github.com/Masa4240/go-mission-catechdojo/middleware"
	"github.com/Masa4240/go-mission-catechdojo/service"
)

func NewRouter(userDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	svc := service.NewUserService(userDB)
	userHandler := handler.NewUserHandler(svc)
	// mux.HandleFunc("/todos", handler.NewTODOHandler(svc).ServeHTTP)
	mux.Handle("/user/create", middleware.Recovery(userHandler))
	mux.Handle("/user/get", middleware.Recovery(middleware.TokenValidation(userHandler)))
	mux.Handle("/user/update", middleware.Recovery(middleware.TokenValidation(userHandler)))

	healthzHandler := handler.NewHealthzHandler()
	mux.HandleFunc("/healthz", healthzHandler.ServeHTTP)
	return mux
}
