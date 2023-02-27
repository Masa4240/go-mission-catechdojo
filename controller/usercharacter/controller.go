package uccontroller

import (
	"errors"
	"time"

	userservice "github.com/Masa4240/go-mission-catechdojo/service/user"
	ucservice "github.com/Masa4240/go-mission-catechdojo/service/usercharacter"
	"github.com/jinzhu/copier"

	"go.uber.org/zap"
)

type UcController struct {
	svc    *ucservice.UcService
	logger *zap.Logger
}

func NewUcController(svc *ucservice.UcService, logger *zap.Logger) *UcController {
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
	var userinfo = userservice.UserInfo{}
	userinfo.ID = req.ID

	list, err := c.svc.GetUserCharacterList(userinfo)
	if err != nil {
		c.logger.Info("ERR", zap.Time("now", time.Now()), zap.Error(err))
		return nil, err
	}
	var res []*UserCharacterRes
	if err := copier.Copy(&res, &list); err != nil {
		return nil, err
	}
	return res, nil
}
