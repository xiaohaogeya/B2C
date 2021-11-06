package models

import (
	_ "github.com/jinzhu/gorm"
)

type Administrator struct {
	Id       int
	Username string
	Password string
	Mobile   string
	Email    string
	Status   int
	RoleId   int `gorm:"roleid"`
	AddTime  int
	IsSuper  int
	Role     Role `gorm:"foreignkey:Id;association_foreignkey:RoleId"`
}

func (Administrator) TableName() string {
	return "administrator"
}
