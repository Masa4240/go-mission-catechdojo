package uchandler

import (
	"encoding/json"
	"net/http"
	"time"

	uccontroller "github.com/Masa4240/go-mission-catechdojo/controller/usercharacter"

	"go.uber.org/zap"
)

type UcHandler struct {
	ctrl   *uccontroller.UcController
	logger *zap.Logger
}

func NewUcHandler(ctrl *uccontroller.UcController, logger *zap.Logger) *UcHandler {
	return &UcHandler{
		ctrl:   ctrl,
		logger: logger,
	}
}

func (h *UcHandler) GetCharacterList(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Start get character List process", zap.Time("now", time.Now()))
	var req = uccontroller.UserCharacterReq{}
	req.ID = int(r.Context().Value("id").(int64))
	res, err := h.ctrl.GetUserCharacterList(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.logger.Info("Error in Get User Character", zap.Time("now", time.Now()), zap.Error(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Info("Encode error", zap.Time("now", time.Now()), zap.Error(err))
		return
	}
	h.logger.Info("Finish Get character process", zap.Time("now", time.Now()))
}
