package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Masa4240/go-mission-catechdojo/model"
	"github.com/Masa4240/go-mission-catechdojo/service"
	"go.uber.org/zap"
)

type GachaHandler struct {
	svc *service.GachaService
}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewGachaHandler(svc *service.GachaService) *GachaHandler {
	return &GachaHandler{
		svc: svc,
	}
}

func (h *GachaHandler) Gacha(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Gacha process", zap.Time("now", time.Now()))

	req := model.GachaReq{}
	var res = []model.GachaResponse{}
	req.ID = int(r.Context().Value("id").(int64))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "decode error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.svc.Gacha(r.Context(), int(r.Context().Value("id").(int64)), req.Times)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//fmt.Fprintln(w, err)
		logger.Info("Error in create user", zap.Time("now", time.Now()))
		return
	}
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("Finish Gacha process", zap.Time("now", time.Now()))
	return
}

func (h *GachaHandler) AddCharacter(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Gacha process", zap.Time("now", time.Now()))

	var req = model.NewCharReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "decode error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.svc.AddCharacter(r.Context(), req.Name, req.Rank, req.Desc, req.Weight)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//fmt.Fprintln(w, err)
		logger.Info("Error in create user", zap.Time("now", time.Now()))
		return
	}

	logger.Info("Finish Gacha process", zap.Time("now", time.Now()))
	return
}

func (h *GachaHandler) GetCharsList(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Gacha process", zap.Time("now", time.Now()))

	var res = []model.GachaResponse{}

	res, err := h.svc.GetCharsList(r.Context(), int(r.Context().Value("id").(int64)))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//fmt.Fprintln(w, err)
		logger.Info("Error in create user", zap.Time("now", time.Now()))
		return
	}
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("Finish Gacha process", zap.Time("now", time.Now()))
	return
}
