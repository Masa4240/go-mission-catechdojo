package router

import (
	"net/http"

	"github.com/Masa4240/go-mission-catechdojo/handler"
)

func NewRouter() *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	// svc := service.NewTODOService(todoDB)
	// mux.HandleFunc("/todos", handler.NewTODOHandler(svc).ServeHTTP)
	//mux.HandleFunc("/do-panic", handler.NewPanicHandler().ServeHTTP)
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
