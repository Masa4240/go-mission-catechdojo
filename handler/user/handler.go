package userhandler

import (
	"encoding/json"
	"net/http"
	"time"

	usercontroller "github.com/Masa4240/go-mission-catechdojo/controller/user"
	"go.uber.org/zap"
)

type UserHandler struct {
	ctrl   *usercontroller.UserController
	logger *zap.Logger
}

func NewUserHandler(
	ctrl *usercontroller.UserController,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		ctrl:   ctrl,
		logger: logger,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Start Create User Process in service", zap.Time("now", time.Now()))
	var req = usercontroller.UserResistrationRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Info("Decode Error", zap.Time("now", time.Now()), zap.Error(err))
		return
	}
	res, err := h.ctrl.CreateUser(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Info("Error in create user", zap.Time("now", time.Now()), zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Info("Error in controller", zap.Time("now", time.Now()), zap.Error(err))
		return
	}
	h.logger.Info("Finish Create User process", zap.Time("now", time.Now()))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Start Get User Process in handler", zap.Time("now", time.Now()))
	id, ok := r.Context().Value("id").(int64)
	if !ok {
		h.logger.Info("Fail to get id from header", zap.Time("now", time.Now()), zap.Int64("id", id))
		return
	}
	var req usercontroller.UserGetRequest
	req.ID = id
	res, err := h.ctrl.GetUser(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Info("Error in get user", zap.Time("now", time.Now()), zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Info("Error in controller", zap.Time("now", time.Now()), zap.Error(err))
		return
	}
	h.logger.Info("Finish Get User process", zap.Time("now", time.Now()))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Start update User Process in handler", zap.Time("now", time.Now()))
	var req = usercontroller.UserUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Info("Decode error", zap.Time("now", time.Now()), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, ok := r.Context().Value("id").(int64)
	if !ok {
		return
	}
	req.ID = id

	if err := h.ctrl.UpdateUserService(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Info("Error in controller", zap.Time("now", time.Now()), zap.Error(err))
		return
	}
	h.logger.Info("Finish update User Process in handler", zap.Time("now", time.Now()))
}
