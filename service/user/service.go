package userservice

import (
	"errors"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	"github.com/form3tech-oss/jwt-go"
)

type UserService struct {
	model *usermodel.UserModel
}

func NewUserService(model *usermodel.UserModel) *UserService {
	return &UserService{
		model: model,
	}
}

func (s *UserService) CreateUser(newName string) (*string, error) {
	var newUser usermodel.UserLists
	newUser.Name = newName

	res, err := s.model.CreateUser(&newUser)
	if err != nil {
		return nil, err
	}

	// Create Token
	token := jwt.New(jwt.SigningMethodHS256)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("fail to create token")
	}
	claims["id"] = res.ID
	tokenString, _ := token.SignedString([]byte("SIGNINGKEY"))

	return &tokenString, nil
}
