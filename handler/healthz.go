package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Masa4240/go-mission-catechdojo/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct {
}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = &model.HealthzResponse{}
	response := &model.HealthzResponse{
		Message: "OK",
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		log.Println(err)
	}
}

func Handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello Handler!")
}
