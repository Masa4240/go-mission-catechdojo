package gachahandler

import (
	"encoding/json"
	"net/http"
	"time"

	gachacontroller "github.com/Masa4240/go-mission-catechdojo/controller/gacha"
	gachamodel "github.com/Masa4240/go-mission-catechdojo/model/gacha"
	"go.uber.org/zap"
)

type GachaHandler struct {
	ctrl *gachacontroller.GachaController
}

func NewGachaHandler(ctrl *gachacontroller.GachaController) *GachaHandler {
	return &GachaHandler{
		ctrl: ctrl,
	}
}

// これはHandlerにあるべき、Viewの下にこれを持ってくる。Serviceの呼び出しをControllerから呼び出す.
func (h *GachaHandler) Gacha(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)
	logger.Info("Start Gacha process", zap.Time("now", time.Now()))

	req := gachamodel.GachaReq{}
	req.ID = int(r.Context().Value("id").(int64))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.ctrl.Gacha(int(r.Context().Value("id").(int64)), req.Times)
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

func (h *GachaHandler) AddCharacter(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)
	logger.Info("Start Gacha process", zap.Time("now", time.Now()))

	var req = gachamodel.NewCharacterReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.ctrl.AddCharacter(req.Name, req.Rank, req.Desc, req.Weight)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Info("Error in create user", zap.Time("now", time.Now()))
		return
	}

	logger.Info("Finish Gacha process", zap.Time("now", time.Now()))
}

func (h *GachaHandler) GetCharacterList(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}(logger)
	logger.Info("Start Gacha process", zap.Time("now", time.Now()))

	// var res = []*model.GachaResponse{}
	id := int(r.Context().Value("id").(int64))
	// if !ok {
	// 	return
	// }
	res, err := h.ctrl.GetUserCharacterList(id)
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
