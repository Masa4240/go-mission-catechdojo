package gachamodel

import "github.com/jinzhu/gorm"

type (
	CharLists struct {
		// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
		gorm.Model
		CharID int
		Name   string
		Rank   string
		Weight int
	}
	NewCharReq struct {
		Name   string `json:"name"`
		Rank   string `json:"rank"`
		Weight int    `json:"weight"`
		Desc   string `json:"desc"`
	}
	GachaReq struct {
		ID    int `json:"userid"`
		Times int `json:"times"`
	}
	GachaResponse struct {
		CharacterID int    `json:"characterID"`
		Name        string `json:"name"`
	}
	RankRatio struct {
		Ranklevel string
		Weight    int
	}

	CharacterLists struct {
		gorm.Model
		CharacterID int
		UserID      int
		Name        string
		Rank        string
		Level       int
		Desc        string
		Weight      int
	}

	UserCharacterList struct {
		gorm.Model
		CharacterID int
		UserID      int
	}
	NewCharacterReq struct {
		Name   string
		Rank   string
		Level  int
		Desc   string
		Weight int
	}
)
