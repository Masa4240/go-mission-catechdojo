package gachacontroller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Masa4240/go-mission-catechdojo/model"
	gachaservice "github.com/Masa4240/go-mission-catechdojo/service/gacha"
	"go.uber.org/zap"
)

type GachaController struct {
	svc *gachaservice.GachaService
}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewGachaController(svc *gachaservice.GachaService) *GachaController {
	return &GachaController{
		svc: svc,
	}
}

// これはHandlerにあるべき、Viewの下にこれを持ってくる。Serviceの呼び出しをControllerから呼び出す
func (h *GachaController) Gacha(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Gacha process", zap.Time("now", time.Now()))

	req := model.GachaReq{}
	req.ID = int(r.Context().Value("id").(int64))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.svc.Gacha(r.Context(), int(r.Context().Value("id").(int64)), req.Times)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Error in create user", zap.Time("now", time.Now()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Convert w header to json
	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("Finish Gacha process", zap.Time("now", time.Now()))
}

func (h *GachaController) AddCharacter(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Gacha process", zap.Time("now", time.Now()))

	var req = model.NewCharacterReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.svc.AddCharacter(r.Context(), req.Name, req.Rank, req.Desc, req.Weight)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Error in create user", zap.Time("now", time.Now()))
		return
	}

	logger.Info("Finish Gacha process", zap.Time("now", time.Now()))
}

func (h *GachaController) GetCharacterList(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Gacha process", zap.Time("now", time.Now()))

	// var res = []*model.GachaResponse{}
	id, ok := int(r.Context().Value("id").(int64))
	if !ok {
		return
	}
	res, err := h.svc.GetUserCharacterList(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Error in create user", zap.Time("now", time.Now()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Info("Finish Gacha process", zap.Time("now", time.Now()))
}
