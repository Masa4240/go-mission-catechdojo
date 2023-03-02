package ucmodel

import (
	"time"

	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type UcModel struct {
	logger *zap.Logger
}

func NewUcModel(logger *zap.Logger) *UcModel {
	return &UcModel{
		logger: logger,
	}
}

func (m *UcModel) GetCharaterList(db *gorm.DB, user *usermodel.UserList) ([]*UserCharacterList, error) {
	m.logger.Info("Start Gacha Process", zap.Time("now", time.Now()))

	var res []*UserCharacterList
	var tempres []UserCharacterList
	tableName := "user_characterlist"
	if err := db.Table(tableName).Where("`user_id` = ?", user.ID).Find(&res).Error; err != nil {
		m.logger.Info("Fail to get char list from db", zap.Time("now", time.Now()), zap.Error(err))
		return nil, err
	}

	for i := 0; i < len(tempres); i++ {
		res = append(res, &tempres[i])
	}

	return res, nil
}

func (m *UcModel) RegisterCharacters(db *gorm.DB, characters []*UserCharacterList) error {
	m.logger.Info("Start Gacha Process", zap.Time("now", time.Now()))

	for i := 0; i < len(characters); i++ {
		res := db.Table("user_characterlist").Create(characters[i])
		if res.Error != nil {
			m.logger.Info("Failed to register new char to user db", zap.Time("now", time.Now()), zap.Error(res.Error))
			return res.Error
		}
	}
	return nil
}
