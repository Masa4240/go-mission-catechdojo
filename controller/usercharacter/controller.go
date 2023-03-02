package uccontroller

import (
	"errors"
	"time"

	userservice "github.com/Masa4240/go-mission-catechdojo/service/user"

	"go.uber.org/zap"
)

type UcController struct {
	svc    *userservice.UserService
	logger *zap.Logger
}

func NewUcController(svc *userservice.UserService, logger *zap.Logger) *UcController {
	return &UcController{
		svc:    svc,
		logger: logger,
	}
}

func (c *UcController) GetUserCharacterList(req UserCharacterReq) ([]*UserCharacterRes, error) {
	c.logger.Info("Start Gacha Process", zap.Time("now", time.Now()))
	if req.ID == 0 {
		return nil, errors.New("empty id")
	}
	userinfo := userservice.UserInfo{
		Profile: &userservice.UserProfile{
			ID: req.ID,
		},
		Characters: nil,
	}

	list, err := c.svc.GetUserCharacterList(userinfo)
	if err != nil {
		c.logger.Info("ERR", zap.Time("now", time.Now()), zap.Error(err))
		return nil, err
	}
	res := convertToUserCharacterRes(*list)
	return res, nil
}
