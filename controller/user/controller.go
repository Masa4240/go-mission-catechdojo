package usercontroller

import (
	"errors"
	"time"
	"unicode/utf8"

	userservice "github.com/Masa4240/go-mission-catechdojo/service/user"
	"go.uber.org/zap"
)

type UserController struct {
	svc    *userservice.UserService
	logger *zap.Logger
}

func NewUserController(svc *userservice.UserService, logger *zap.Logger) *UserController {
	return &UserController{
		svc:    svc,
		logger: logger,
	}
}

func (c *UserController) CreateUser(req UserResistrationRequest) (*UserResistrationResponse, error) {
	c.logger.Info("Start Create User Process in controller", zap.Time("now", time.Now()), zap.String("new name is", req.Name))
	var info userservice.UserInfo
	info.Name = req.Name
	if err := c.nameValidation(info); err != nil {
		c.logger.Info("Wrong Name", zap.Time("now", time.Now()), zap.String("new name is", req.Name))
		return nil, err
	}
	res, err := c.svc.CreateUser(info)
	if err != nil {
		return nil, err
	}
	var response = UserResistrationResponse{}
	response.Token = res.Token
	return &response, nil
}

func (c *UserController) GetUser(req UserGetRequest) (*UserGetResponse, error) {
	c.logger.Info("Start to get User name in Controller", zap.Time("now", time.Now()), zap.Int64("Requested ID", req.ID))
	var info userservice.UserInfo
	info.ID = int(req.ID)
	res, err := c.svc.GetUserByID(info)
	if err != nil {
		c.logger.Info("ID not found", zap.Time("now", time.Now()), zap.Int64("req ID", req.ID), zap.Error(err))
		return nil, errors.New("id not found")
	}
	var response = UserGetResponse{}
	response.Name = res.Name
	return &response, nil
}

func (c *UserController) UpdateUserService(req UserUpdateRequest) error {
	var info userservice.UserInfo
	info.Name = req.Newname
	info.ID = int(req.ID)
	if err := c.nameValidation(info); err != nil {
		c.logger.Info("Wrong Name", zap.Time("now", time.Now()), zap.String("new name is", info.Name))
		return err
	}
	if err := c.svc.UpdateUser(info); err != nil {
		return err
	}
	return nil
}

func (c *UserController) nameValidation(info userservice.UserInfo) error {
	if len(info.Name) == 0 {
		err := errors.New("null name")
		return err
	}
	maxLength := 10
	if utf8.RuneCountInString(info.Name) > maxLength {
		err := errors.New("too long name")
		return err
	}
	// Duplication check
	if err := c.svc.DuplicationCheck(info); err != nil {
		return errors.New("duplicated")
	}
	return nil
}
