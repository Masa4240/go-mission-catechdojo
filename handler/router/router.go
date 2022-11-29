package router

import (
	"database/sql"
	"net/http"

	"github.com/Masa4240/go-mission-catechdojo/auth"
	"github.com/Masa4240/go-mission-catechdojo/handler"
	"github.com/Masa4240/go-mission-catechdojo/service"
)

func NewRouter(userDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	svc := service.NewTODOService(userDB)
	// mux.HandleFunc("/todos", handler.NewTODOHandler(svc).ServeHTTP)
	mux.HandleFunc("/user/create", handler.NewUserHandler(svc).ServeHTTP)
	mux.Handle("/auth", auth.GetTokenHandler)
	// healthzHandler := handler.NewHealthzHandler()
	// mux.HandleFunc("/healthz", healthzHandler.ServeHTTP)
	healthzHandler := handler.NewHealthzHandler()
	mux.HandleFunc("/healthz", healthzHandler.ServeHTTP)
	//mux.Handle("/healthz2", middleware.BasicAuth(middleware.DetectOS(middleware.AccessLogOutput1(healthzHandler))))
	// panicHandler := handler.NewPanicHandler()
	// mux.HandleFunc("/do-panic1", panicHandler.ServeHTTP)
	// mux.Handle("/do-panic", middleware.Recovery(panicHandler))

	return mux
}
