package models

import (
	_ "github.com/jinzhu/gorm"
)

type ProductTypeAttribute struct {
	Id        int    `json:"id"`
	CateId    int    `json:"cate_id"`
	Title     string `json:"title"`
	AttrType  int    `json:"attr_type"`
	AttrValue string `json:"attr_value"`
	Status    int    `json:"status"`
	Sort      int    `json:"sort"`
	AddTime   int    `json:"add_time"`
}

func (ProductTypeAttribute) TableName() string {
	return "product_type_attribute"
}
