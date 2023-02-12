package usercontroller

import (
	"encoding/json"
	"net/http"
	"time"

	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	userservice "github.com/Masa4240/go-mission-catechdojo/service/user"
	"go.uber.org/zap"
)

type UserController struct {
	svc *userservice.UserService
}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewUserController(svc *userservice.UserService) *UserController {
	return &UserController{
		svc: svc,
	}
}

func (h *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Create User Process in service", zap.Time("now", time.Now()))
	var req = usermodel.UserResistrationRequest{}
	var resp = usermodel.UserResistrationResponse{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userinfo, err := h.svc.CreateUserService(r.Context(), req.Name)
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

func (h *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Create User Process in service", zap.Time("now", time.Now()))
	var res = usermodel.UserGetResponse{}
	id := r.Context().Value("id").(int64)

	name, err := h.svc.GetUserService(r.Context(), int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Error in create user", zap.Time("now", time.Now()))
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

func (h *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start update User Process in controller", zap.Time("now", time.Now()))

	var req = usermodel.UserUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Decode error", zap.Time("now", time.Now()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := r.Context().Value("id").(int64)

	if err := h.svc.UpdateUserService(r.Context(), req.Newname, int(id)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Error in service", zap.Time("now", time.Now()))
		return
	}
	logger.Info("Finish update User Process in controller", zap.Time("now", time.Now()))
}
