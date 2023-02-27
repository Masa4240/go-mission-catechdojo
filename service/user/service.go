package userservice

import (
	"errors"
	"time"

	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	"github.com/form3tech-oss/jwt-go"
	"go.uber.org/zap"
)

type UserService struct {
	model  *usermodel.UserModel
	logger *zap.Logger
}

func NewUserService(model *usermodel.UserModel, logger *zap.Logger) *UserService {
	return &UserService{
		model:  model,
		logger: logger,
	}
}

func (s *UserService) CreateUser(info UserInfo) (*UserInfo, error) {
	var newUser usermodel.UserList
	newUser.Name = info.Name
	res, err := s.model.CreateUser(&newUser)
	if err != nil {
		s.logger.Info("Create User Service", zap.Time("now", time.Now()), zap.Error(err))
		return nil, err
	}
	// Create Token
	token := jwt.New(jwt.SigningMethodHS256)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Info("Fail to create token", zap.Time("now", time.Now()))
		return nil, errors.New("fail to create token")
	}
	claims["id"] = res.ID
	tokenString, _ := token.SignedString([]byte("SIGNINGKEY"))
	info.Token = tokenString
	return &info, nil
}

func (s *UserService) GetUserByID(info UserInfo) (*UserInfo, error) {
	var user usermodel.UserList
	user.ID = uint(info.ID)
	res, err := s.model.GetUserByID(&user)
	if err != nil {
		s.logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		return nil, errors.New("id not found")
	}
	info.Name = res.Name
	return &info, nil
}

func (s *UserService) DuplicationCheck(info UserInfo) error {
	var user usermodel.UserList
	user.Name = info.Name
	res, err := s.model.GetUserByName(&user)
	if err != nil {
		s.logger.Info("Get User By Name Error", zap.Time("now", time.Now()), zap.Error(err))
		return errors.New("not found")
	}
	if len(res) != 0 {
		s.logger.Info("Duplicated", zap.Time("now", time.Now()), zap.Uint("ID", res[0].ID))
		return errors.New("not found")
	}
	return nil
}

func (s *UserService) UpdateUser(info UserInfo) error {
	var user usermodel.UserList
	user.ID = uint(info.ID)
	user.Name = info.Name
	if err := s.model.UpdateUser(&user); err != nil {
		return errors.New("fail to confirm new db")
	}
	return nil
}
