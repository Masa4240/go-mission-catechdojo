package userhandler

import (
	"encoding/json"
	"net/http"
	"time"

	usercontroller "github.com/Masa4240/go-mission-catechdojo/controller/user"
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
		err := logger.Sync()
		if err != nil {
			// fmt.Println(err)
			return
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

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Start Get User Process in service", zap.Time("now", time.Now()))
	var res = usermodel.UserGetResponse{}
	id, ok := r.Context().Value("id").(int64)
	if !ok {
		logger.Info("Fail to get id from header", zap.Time("now", time.Now()), zap.Int64("id", id))
		return
	}

	name, err := h.ctrl.GetUserService(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Error in create user", zap.Time("now", time.Now()), zap.Error(err))
		return
	}
	res.Name = *name
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Info("Finish Create User process", zap.Time("now", time.Now()))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Start update User Process in controller", zap.Time("now", time.Now()))

	var req = usermodel.UserUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Decode error", zap.Time("now", time.Now()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, ok := r.Context().Value("id").(int64)
	if !ok {
		return
	}

	if err := h.ctrl.UpdateUserService(req.Newname, int(id)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Error in service", zap.Time("now", time.Now()))
		return
	}
	logger.Info("Finish update User Process in controller", zap.Time("now", time.Now()))
}
