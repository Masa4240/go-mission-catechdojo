package gachahandler

import (
	"encoding/json"
	"net/http"
	"time"

	gachacontroller "github.com/Masa4240/go-mission-catechdojo/controller/gacha"
	"go.uber.org/zap"
)

type GachaHandler struct {
	ctrl   *gachacontroller.GachaController
	logger *zap.Logger
}

func NewGachaHandler(ctrl *gachacontroller.GachaController, logger *zap.Logger) *GachaHandler {
	return &GachaHandler{
		ctrl:   ctrl,
		logger: logger,
	}
}

func (h *GachaHandler) Gacha(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Start Gacha process", zap.Time("now", time.Now()))

	req := gachacontroller.GachaReq{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Info("Decode error", zap.Time("now", time.Now()), zap.Error(err))
		return
	}
	req.ID = int(r.Context().Value("id").(int64))

	res, err := h.ctrl.Gacha(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Info("Error in create user", zap.Time("now", time.Now()), zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// Convert w header to json
	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Info("Error in encode", zap.Time("now", time.Now()), zap.Error(err))
		return
	}

	h.logger.Info("Finish Gacha process", zap.Time("now", time.Now()))
}
