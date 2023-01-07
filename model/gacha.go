package model

type (
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
