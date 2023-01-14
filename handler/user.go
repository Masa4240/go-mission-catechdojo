package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Masa4240/go-mission-catechdojo/model"
	"github.com/Masa4240/go-mission-catechdojo/service"
	"go.uber.org/zap"
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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Create User Process in service", zap.Time("now", time.Now()))
	var req = model.UserResistrationRequest{}
	var resp = model.UserResistrationResponse{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Empty name", zap.Time("now", time.Now()))
		return
	}
	userinfo, err := h.svc.CreateUser(r.Context(), req.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Error in create user", zap.Time("now", time.Now()))
		return
	}
	resp.Token = userinfo.Token
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Info("Finish Create User process", zap.Time("now", time.Now()))
	return
}

func (h *UserHandler) GetUserName(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Get User Process in handler", zap.Time("now", time.Now()))
	var req = model.UserGetRequest{}
	req.ID = r.Context().Value("id").(int64)
	resp, err := h.svc.GetUser(r.Context(), int(req.ID))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("error in service", zap.Time("now", time.Now()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Info("Finish Get User process", zap.Time("now", time.Now()))
	return
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start update User Process in handler", zap.Time("now", time.Now()))
	var req = model.UserUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Decode error", zap.Time("now", time.Now()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Newname == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Empty name", zap.Time("now", time.Now()))
		return
	}
	req.ID = r.Context().Value("id").(int64)
	if _, err := h.svc.UpdateUser(r.Context(), req.Newname, int(req.ID)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Error in service", zap.Time("now", time.Now()))
		return
	}
	logger.Info("Finish update User Process in handler", zap.Time("now", time.Now()))
	return
}
