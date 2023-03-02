package gachacontroller

import gachaservice "github.com/Masa4240/go-mission-catechdojo/service/gacha"

func convertToGachaResponse(info gachaservice.GachaContent) []*GachaResponse {
	res := make([]*GachaResponse, 0, len(info.Characters))
	for _, model := range info.Characters {
		res = append(res, &GachaResponse{
			CharacterID: model.CharacterID,
			Name:        model.Name,
		})
	}
	return res
}
