package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string `json:"password"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Database struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBname   string `json:"DBname"`
}
