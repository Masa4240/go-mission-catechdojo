package gachaservice

import (
	ucmodel "github.com/Masa4240/go-mission-catechdojo/model/usercharacter"
)

func convertToUserCharacter(content GachaContent) []*ucmodel.UserCharacterList {
	characters := make([]*ucmodel.UserCharacterList, 0, len(content.Characters))

	for _, model := range content.Characters {
		characters = append(characters, &ucmodel.UserCharacterList{
			CharacterID: model.CharacterID,
			UserID:      content.Request.ID,
		})
	}
	return characters
}
