package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string `json:"password"`
	Result   bool   `json:"result"`
}

type Kindergarten struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Inn         int     `json:"inn"`
	PhoneNumber string  `json:"phone_number"`
	Subtitle    string  `json:"subtitle"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Address     string  `json:"address"`
}

type KindergartenResponse struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	Name      string  `json:"name"`
	Inn       int     `json:"inn"`
	Address   string  `json:"address"`
	Number    string  `json:"number"`
	Subtitle  string  `json:"subtitle"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type KindergartenBasicInfoResponse struct {
	ID        uint    `gorm:"primarykey"`
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

//type KindergartenPicture struct {
//	gorm.Model
//	KindergartenID uint   `json:"kindergarten_id"`
//	PicturePath    string `json:"picture_path"`
//}

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

type Child struct {
	gorm.Model
	FullName   string       `json:"full_name"`
	BirthDate  time.Time    `json:"birth_date"`
	Group      string       `json:"group"`
	Attendance []Attendance `json:"attendance" gorm:"foreignKey:ChildID"`
}

type Attendance struct {
	gorm.Model
	ChildID    uint      `json:"child_id"`
	Present    bool      `json:"present"`
	RecordedAt time.Time `json:"recorded_at"`
}

type UpdateAttendancesRequest struct {
	Attendance []Attendance `json:"attendance"`
}

type AttendanceUpdate struct {
	ChildID uint `json:"child_id"`
	Present bool `json:"present"`
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
