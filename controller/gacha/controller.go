package gachacontroller

import (
	"time"

	gachaservice "github.com/Masa4240/go-mission-catechdojo/service/gacha"

	"go.uber.org/zap"
)

type GachaController struct {
	svc    *gachaservice.GachaService
	logger *zap.Logger
}

func NewGachaController(svc *gachaservice.GachaService, logger *zap.Logger) *GachaController {
	return &GachaController{
		svc:    svc,
		logger: logger,
	}
}

func (c *GachaController) Gacha(req GachaReq) ([]*GachaResponse, error) {
	c.logger.Info("Start Gacha Process in controller", zap.Time("now", time.Now()))
	// // Gacha
	greq := gachaservice.GachaRequest{
		ID:    req.ID,
		Times: req.Times,
	}
	result, err := c.svc.Gacha(greq)
	if err != nil {
		c.logger.Info("Error in service", zap.Time("now", time.Now()), zap.Error(err))
		return nil, err
	}
	// res := convertToGachaResponse(*result)
	return convertToGachaResponse(*result), nil
}
