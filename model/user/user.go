package usermodel

import "github.com/jinzhu/gorm"

type (
	UserList struct {
		// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
		gorm.Model
		Name string
	}
)
