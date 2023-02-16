package userhandler

import (
	"encoding/json"
	usercontroller "github.com/Masa4240/go-mission-catechdojo/controller/user"
	"net/http"
	"time"

	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	"go.uber.org/zap"
)

type UserHandler struct {
	ctrl *usercontroller.UserController
}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewUserHandler(ctrl *usercontroller.UserController) *UserHandler {
	return &UserHandler{
		ctrl: ctrl,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}(logger)
	logger.Info("Start Create User Process in service", zap.Time("now", time.Now()))
	var req = usermodel.UserResistrationRequest{}
	var resp = usermodel.UserResistrationResponse{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userinfo, err := h.ctrl.CreateUserService(req.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Error in create user", zap.Time("now", time.Now()))
		return
	}
	resp.Token = *userinfo
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Info("Finish Create User process", zap.Time("now", time.Now()))
}
