package userservice

import (
	"errors"
	"time"

	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	ucmodel "github.com/Masa4240/go-mission-catechdojo/model/usercharacter"
	"github.com/form3tech-oss/jwt-go"
	"go.uber.org/zap"
)

type UserService struct {
	ucmodel   *ucmodel.UcModel
	usermodel *usermodel.UserModel
	logger    *zap.Logger
}

func NewUserService(ucmodel *ucmodel.UcModel, usermodel *usermodel.UserModel, logger *zap.Logger) *UserService {
	return &UserService{
		ucmodel:   ucmodel,
		usermodel: usermodel,
		logger:    logger,
	}
}

func (s *UserService) CreateUser(info UserInfo) (*UserInfo, error) {
	// var newUser usermodel.UserList
	req := usermodel.UserList{
		Name: info.Profile.Name,
	}

	// newUser.Name = info.Name
	res, err := s.usermodel.CreateUser(&req)
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
	info.Profile.Token = tokenString
	return &info, nil
}

func (s *UserService) GetUserByID(info UserInfo) (*UserInfo, error) {
	user := usermodel.UserList{}
	user.ID = uint(info.Profile.ID)

	res, err := s.usermodel.GetUserByID(&user)
	if err != nil {
		s.logger.Info("ID Not Found", zap.Time("now", time.Now()), zap.Error(err))
		return nil, errors.New("id not found")
	}
	info.Profile.Name = res.Name
	// info.Name = res.Name
	return &info, nil
}

func (s *UserService) DuplicationCheck(info UserInfo) error {
	user := usermodel.UserList{
		Name: info.Profile.Name,
	}
	// user.Name = info.Name
	res, err := s.usermodel.GetUserByName(&user)
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
	user := usermodel.UserList{
		// ID:   uint(info.Profile.ID),
		Name: info.Profile.Name,
	}
	user.ID = uint(info.Profile.ID)

	if err := s.usermodel.UpdateUser(&user); err != nil {
		return errors.New("fail to confirm new db")
	}
	return nil
}

func (s *UserService) GetUserCharacterList(info UserInfo) (*UserInfo, error) {
	req := usermodel.UserList{
		// ID: uint(info.Profile.ID),
	}
	req.ID = uint(info.Profile.ID)

	list, err := s.ucmodel.GetCharaterList(&req)
	if err != nil {
		return nil, err
	}
	res := convertToUserInfo(&req, list)
	return res, nil
}
