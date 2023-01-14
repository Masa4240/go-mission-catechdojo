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
		//Level int
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
		CharID int    `json:"characterID"`
		Name   string `json:"name"`
	}
)
