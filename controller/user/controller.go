package usercontroller

import (
	"errors"
	"time"
	"unicode/utf8"

	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	userservice "github.com/Masa4240/go-mission-catechdojo/service/user"
	"go.uber.org/zap"
)

type UserController struct {
	svc *userservice.UserService
}

func NewUserController(svc *userservice.UserService) *UserController {
	return &UserController{
		svc: svc,
	}
}

func (c *UserController) CreateUserService(newName string) (*string, error) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Start Create User Process in service", zap.Time("now", time.Now()), zap.String("new name is", newName))

	if err := c.nameValidation(newName); err != nil {
		logger.Info("Wrong Name", zap.Time("now", time.Now()), zap.String("new name is", newName))
		return nil, err
	}
	res, err := c.svc.CreateUser(newName)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *UserController) GetUserService(reqID int) (*string, error) {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	logger.Info("Start to get User name", zap.Time("now", time.Now()), zap.Int("Requested ID", reqID))

	res, err := c.svc.GetUserByID(reqID)
	if err != nil {
		logger.Info("ID not found", zap.Time("now", time.Now()), zap.Int("req ID", reqID), zap.Error(err))
		return nil, errors.New("id not found")
	}
	return res, nil
}

func (c *UserController) UpdateUserService(newName string, reqID int) error {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	if err := c.nameValidation(newName); err != nil {
		logger.Info("Wrong Name", zap.Time("now", time.Now()), zap.String("new name is", newName))
		return err
	}

	if err := c.svc.UpdateUser(newName, reqID); err != nil {
		return err
	}
	return nil
}

func (c *UserController) nameValidation(newName string) error {
	if len(newName) == 0 {
		err := errors.New("null name")
		return err
	}
	maxLength := 10
	if utf8.RuneCountInString(newName) > maxLength {
		err := errors.New("too long name")
		return err
	}
	// Duplication check
	var user usermodel.UserLists
	user.Name = newName
	if err := c.svc.GetUserByName(newName); err != nil {
		return errors.New("duplicated")
	}
	return nil
}
