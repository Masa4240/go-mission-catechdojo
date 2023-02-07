package model

import "github.com/jinzhu/gorm"

type (
	CharLists struct {
		// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
		gorm.Model
		CharID int
		Name   string
		Rank   string
		Weight int
		// Level int
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
		// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
		gorm.Model
		CharacterID int
		UserID      int
		Name        string
		// ここはマスターリストからとれるから重複になる。正規化。二重管理になってしまう マスターデータが変わったときに追従できない
		// Name   string
		Rank   string
		Level  int
		Desc   string
		Weight int
	}

	UserCharacterList struct {
		// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
		gorm.Model
		CharacterID int
		UserID      int
		// ここはマスターリストからとれるから重複になる。正規化。二重管理になってしまう マスターデータが変わったときに追従できない
		// Name   string
		// Rank   string
		// Level int
	}
	NewCharacterReq struct {
		Name   string
		Rank   string
		Level  int
		Desc   string
		Weight int
	}
)
