package gachaservice

import (
	// gachacontroller "github.com/Masa4240/go-mission-catechdojo/controller/gacha"
	gachamodel "github.com/Masa4240/go-mission-catechdojo/model/gacha"
	usermodel "github.com/Masa4240/go-mission-catechdojo/model/user"
)

type GachaService struct {
	model *gachamodel.GachaModel
}

func NewGachaService(model *gachamodel.GachaModel) *GachaService {
	return &GachaService{
		model: model,
	}
}

func (s *GachaService) Gacha(newCharacters []*gachamodel.UserCharacterList) ([]*gachamodel.GachaResponse, error) {
	if err := s.model.RegisterCharacters(newCharacters); err != nil {
		return nil, err
	}
	var newResCharacters []*gachamodel.GachaResponse
	newResCharacters = resConverter(newCharacters)
	return newResCharacters, nil
}

func (s *GachaService) GetUserCharacterList(id int) ([]*gachamodel.GachaResponse, error) {
	var req usermodel.UserLists
	req.ID = uint(id)

	list, err := s.model.GetCharaterList(&req)
	if err != nil {
		return nil, err
	}
	res := resConverter(list)
	return res, nil
}

func resConverter(characterList []*gachamodel.UserCharacterList) []*gachamodel.GachaResponse {
	resCharacters := []*gachamodel.GachaResponse{}
	for i := 0; i < len(characterList); i++ {
		resCharacter := gachamodel.GachaResponse{}
		resCharacter.CharacterID = characterList[i].CharacterID
		// resCharacter.Name = gachacontroller.CharacterMasterDataMVC[characterList[i].CharacterID].Name
		resCharacters = append(resCharacters, &resCharacter)
	}
	return resCharacters
}

func (s *GachaService) AddCharacter(name, rank, desc string, weight int) error {
	newCharacter := gachamodel.CharacterLists{}
	newCharacter.Name = name
	newCharacter.Rank = rank
	newCharacter.Desc = desc
	newCharacter.Weight = weight
	if err := s.model.AddNewCharacter(&newCharacter); err != nil {
		return err
	}
	return nil
}

func (s *GachaService) GetForamalCharacterList() ([]*gachamodel.CharacterLists, error) {
	res, err := s.model.GetForamalCharacterList()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *GachaService) GetRankMasterData() ([]*gachamodel.RankRatio, error) {
	res, err := s.model.GetRankRatio()
	if err != nil {
		return nil, err
	}
	return res, nil
}
