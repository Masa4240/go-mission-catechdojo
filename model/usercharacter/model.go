package ucmodel

import (
	"time"

	charactermodel "github.com/Masa4240/go-mission-catechdojo/model/character"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type UcModel struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUcModel(db *gorm.DB, logger *zap.Logger) *UcModel {
	return &UcModel{
		db:     db,
		logger: logger,
	}
}

func (m *UcModel) GetCharaterList(user *usermodel.UserList) ([]*charactermodel.CharacterData, error) {
	m.logger.Info("Start Gacha Process", zap.Time("now", time.Now()))

	var res []*charactermodel.CharacterData
	var tempres []charactermodel.CharacterData
	tableName := "user_characterlist"
	tx := m.db.Begin()
	if err := tx.Table(tableName).Where("`user_id` = ?", user.ID).Find(&res).Error; err != nil {
		m.logger.Info("Fail to get char list from db", zap.Time("now", time.Now()), zap.Error(err))
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	for i := 0; i < len(tempres); i++ {
		res = append(res, &tempres[i])
	}

	return res, nil
}

func (m *UcModel) RegisterCharacters(characters []*UserCharacterList) error {
	m.logger.Info("Start Gacha Process", zap.Time("now", time.Now()))
	tx := m.db.Begin()

	for i := 0; i < len(characters); i++ {
		res := tx.Table("user_characterlist").Create(characters[i])
		if res.Error != nil {
			m.logger.Info("Failed to register new char to user db", zap.Time("now", time.Now()), zap.Error(res.Error))
			tx.Rollback()
			return res.Error
		}
	}
	tx.Commit()
	return nil
}
