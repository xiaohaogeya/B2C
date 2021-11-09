package models

import (
	_ "github.com/jinzhu/gorm"
)

type RoleAuth struct {
	AuthId int
	RoleId   int
}

func (RoleAuth) TableName() string {
	return "role_auth"
}
