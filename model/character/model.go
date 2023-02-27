package charactermodel

import (
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type CharacterModel struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewCharacterModel(db *gorm.DB, logger *zap.Logger) *CharacterModel {
	return &CharacterModel{
		db:     db,
		logger: logger,
	}
}

var (
	CharacterMasterData map[int]CharacterData
	srCharacters        []*CharacterData
	rCharacters         []*CharacterData
	nCharacters         []*CharacterData
)

func (s *CharacterModel) GetForamalCharacterList() error {
	s.logger.Info("Start Gacha Process", zap.Time("now", time.Now()))
	characters := []CharacterList{}
	tx := s.db.Begin()

	if err := s.db.Table("formal_character_list").Find(&characters).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	CharacterMasterData = map[int]CharacterData{}
	for i := 0; i < len(characters); i++ {
		var tempCharacter = CharacterData{}
		tempCharacter.CharacterID = int(characters[i].ID)
		tempCharacter.Name = characters[i].Name
		CharacterMasterData[int(characters[i].ID)] = tempCharacter
		if characters[i].Rank == "SR" {
			for j := 0; j < characters[i].Weight; j++ {
				srCharacters = append(srCharacters, &tempCharacter)
			}
		}
		if characters[i].Rank == "R" {
			for j := 0; j < characters[i].Weight; j++ {
				rCharacters = append(rCharacters, &tempCharacter)
			}
		}
		if characters[i].Rank == "N" {
			for j := 0; j < characters[i].Weight; j++ {
				nCharacters = append(nCharacters, &tempCharacter)
			}
		}
	}
	return nil
}

func GetCharacterByID(data CharacterData) CharacterData {
	return CharacterMasterData[data.CharacterID]
}

func GetSRCharacters() []*CharacterData {
	return srCharacters
}

func GetRCharacters() []*CharacterData {
	return rCharacters
}
func GetNCharacters() []*CharacterData {
	return nCharacters
}
