package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Masa4240/go-mission-catechdojo/model"
	"github.com/Masa4240/go-mission-catechdojo/service"
)

// A HealthzHandler implements health check endpoint.
type UserHandler struct {
	svc *service.UserService
}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) Create(ctx context.Context, req *model.UserResistrationRequest) (*model.UserResistrationResponse, error) {
	//	fmt.Fprint(nil, "Hello TODO!")
	_, err := h.svc.CreateUser(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	// return &model.CreateTODOResponse{TODO: *todo}, nil
	return &model.UserResistrationResponse{}, nil
}

// ServeHTTP implements http.Handler interface.
func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	println("user handler")
	if r.Method == "POST" {
		var req = model.UserResistrationRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "decode error")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Name is empty")
			//fmt.Fprintln(w, r)
			return
		}
		fmt.Fprintln(w, req.Name)
		resp, err := h.Create(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, err, resp)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Incorrect request")
}
