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

func (s *UserService) GetUserByID(id int) (*string, error) {
	var user usermodel.UserLists
	user.ID = uint(id)
	res, err := s.model.GetUserByID(&user)
	if err != nil {
		return nil, errors.New("id not found")
	}
	return &res.Name, nil
}

func (s *UserService) GetUserByName(name string) error {
	var user usermodel.UserLists
	user.Name = name
	res, err := s.model.GetUserByName(&user)
	if err != nil {
		return errors.New("not found")
	}
	if len(res) != 0 {
		return errors.New("not found")
	}
	return nil
}

func (s *UserService) UpdateUser(newName string, id int) error {
	var user usermodel.UserLists
	user.ID = uint(id)
	user.Name = newName
	if err := s.model.UpdateUser(&user); err != nil {
		return errors.New("fail to confirm new db")
	}
	return nil
}
