package charactermodel

import "github.com/jinzhu/gorm"

type (
	CharacterList struct {
		gorm.Model
		// CharacterID int
		Name   string
		Rank   string
		Weight int
	}

	CharacterData struct {
		CharacterID int
		Name        string
	}
)
