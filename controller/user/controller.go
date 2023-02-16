package usercontroller

import (
	"errors"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	userservice "github.com/Masa4240/go-mission-catechdojo/service/user"
	"go.uber.org/zap"
	"time"
	"unicode/utf8"
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
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}(logger)
	logger.Info("Start Create User Process in service", zap.Time("now", time.Now()), zap.String("new name is", newName))

	if err := c.nameValidation(newName); err != nil {
		logger.Info("Wrong Name", zap.Time("now", time.Now()), zap.String("new name is", newName))
		return nil, err
	}
	var newUser usermodel.UserLists
	newUser.Name = newName

	if len(newName) == 0 {
		err := errors.New("null name")
		return nil, err
	}
	maxLength := 10
	if utf8.RuneCountInString(newName) > maxLength {
		err := errors.New("too long name")
		return nil, err
	}

	res, err := c.svc.CreateUser(&newUser)
	if err != nil {
		return nil, err
	}
	logger.Info("Finish Create User Service Registration", zap.Time("now", time.Now()), zap.Uint("New ID", res.ID))

	return res, nil
}

func (c *UserController) nameValidation(newName string) error {
	// Duplication check
	var user usermodel.UserLists
	user.Name = newName
	res, err := c.svc.GetUserByName(&user)
	if err != nil {
		return errors.New("fail to get user from db")
		// return err
	}
	if len(res) != 0 {
		return errors.New("duplicated")
		// return err
	}
	return nil
}
