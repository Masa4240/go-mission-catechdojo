package userservice

import (
	charactermodel "github.com/Masa4240/go-mission-catechdojo/model/character"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	ucmodel "github.com/Masa4240/go-mission-catechdojo/model/usercharacter"
)

func convertToUserInfo(userModel *usermodel.UserList, characterList []*ucmodel.UserCharacterList) *UserInfo {
	characters := make([]*UserCharacter, 0, len(characterList))
	for _, model := range characterList {
		characters = append(characters, &UserCharacter{
			ID:   model.CharacterID,
			Name: charactermodel.GetCharacterNameyID(model.CharacterID),
		})
	}
	return &UserInfo{
		Profile: &UserProfile{
			ID:   int(userModel.ID),
			Name: userModel.Name,
		},
		Characters: characters,
	}
}
