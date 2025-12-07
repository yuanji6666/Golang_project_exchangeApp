package models

import "gorm.io/gorm"

type Article struct{
	gorm.Model
	Title string `json:"required"`
	Content string`json:"required"`
	Preview string`json:"required"`
	Likes int `gorm:"default:0"`
}