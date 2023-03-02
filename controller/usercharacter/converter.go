package uccontroller

import (
	// usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
	// ucmodel "github.com/Masa4240/go-mission-catechdojo/model/usercharacter"
	userservice "github.com/Masa4240/go-mission-catechdojo/service/user"
)

func convertToUserCharacterRes(info userservice.UserInfo) []*UserCharacterRes {
	res := make([]*UserCharacterRes, 0, len(info.Characters))
	for _, model := range info.Characters {
		res = append(res, &UserCharacterRes{
			CharacterID: model.ID,
			Name:        model.Name,
		})
	}
	return res
}
