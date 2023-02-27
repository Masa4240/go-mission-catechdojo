package ucmodel

import "github.com/jinzhu/gorm"

type (
	UserCharacterList struct {
		gorm.Model
		CharacterID int
		UserID      int
	}
)
