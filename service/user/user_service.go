package userservice

import (
	"context"
	"errors"
	"time"
	"unicode/utf8"

	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	"github.com/form3tech-oss/jwt-go"
	"go.uber.org/zap"
)

type UserService struct {
	svc *usermodel.UserModel
}

func NewUserService(svc *usermodel.UserModel) *UserService {
	return &UserService{
		svc: svc,
	}
}

func (s *UserService) CreateUserService(ctx context.Context, newName string) (*string, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start Create User Process in service", zap.Time("now", time.Now()), zap.String("new name is", newName))

	if err := s.nameValidation(ctx, newName); err != nil {
		logger.Info("Wrong Name", zap.Time("now", time.Now()), zap.String("new name is", newName))
		return nil, err
	}
	var newUser usermodel.UserLists
	newUser.Name = newName

	res, err := s.svc.CreateUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}
	logger.Info("Finish Create User Service Registration", zap.Time("now", time.Now()), zap.Uint("New ID", res.ID))

	// Create Token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = res.ID
	tokenString, _ := token.SignedString([]byte("SIGNINGKEY"))

	return &tokenString, nil
}

func (s *UserService) GetUserService(ctx context.Context, reqID int) (*string, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("Start to get User name", zap.Time("now", time.Now()), zap.Int("Requested ID", reqID))
	var user usermodel.UserLists
	user.ID = uint(reqID)

	res, err := s.svc.GetUserById(ctx, &user)
	if err != nil {
		logger.Info("ID not found", zap.Time("now", time.Now()), zap.Int("req ID", reqID), zap.Error(err))
		err := errors.New("id not found")
		return nil, err
	}
	return &res.Name, nil
}

func (s *UserService) UpdateUserService(ctx context.Context, newName string, reqID int) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	if err := s.nameValidation(ctx, newName); err != nil {
		logger.Info("Wrong Name", zap.Time("now", time.Now()), zap.String("new name is", newName))
		return err
	}
	var userInfo usermodel.UserLists
	userInfo.Name = newName
	userInfo.ID = uint(reqID)

	// tx := s.db.Begin()
	// 	tx := s.db.Lock()をSelectするタイミングで呼ぶ必要がある
	if err := s.svc.UpdateUser(ctx, &userInfo); err != nil {
		logger.Info("Fail to update new DB", zap.Time("now", time.Now()), zap.String("name", newName), zap.Int("id", reqID))
		err := errors.New("fail to confirm new db")
		//	tx.Rollback()
		return err
	}
	// tx.Commit()
	return nil
}

func (s *UserService) nameValidation(ctx context.Context, newName string) error {
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
	res, err := s.svc.GetUserByName(ctx, &user)
	if err != nil {
		err := errors.New("fail to get user from db")
		return err
	}
	if len(res) != 0 {
		err := errors.New("duplicated")
		return err
	}
	return nil
}
