package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string `json:"password"`
	Result   bool   `json:"result"`
}

type Kindergarten struct {
	gorm.Model
	Name        string                `json:"name"`
	Inn         int                   `json:"inn"`
	Address     string                `json:"address"`
	Number      int                   `json:"number"`
	Picture     []KindergartenPicture `json:"pictures" gorm:"foreignKey:KindergartenID"`
	Description string                `json:"description"`
}

type KindergartenPicture struct {
	gorm.Model
	KindergartenID uint   `json:"kindergarten_id"`
	PicturePath    string `json:"picture_path"`
}

type EducationMinistry struct {
	gorm.Model
	Name        string `json:"name"`
	Password    string `json:"-"`
	Inn         int    `json:"inn"`
	PhoneNumber int    `json:"phone_number"`
}

type EducationMinistryArchive struct {
	gorm.Model
	OriginalID  uint   `json:"original_id"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Inn         int    `json:"inn"`
	PhoneNumber int    `json:"phoneNumber"`
}

type MainDepartment struct {
	gorm.Model
	Name        string `json:"name"`
	Password    string `json:"-"`
	Inn         int    `json:"inn"`
	PhoneNumber int    `json:"phone_number"`
}

type MainDepartmentArchive struct {
	gorm.Model
	OriginalID  uint   `json:"original_id"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Inn         int    `json:"inn"`
	PhoneNumber int    `json:"phoneNumber"`
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
