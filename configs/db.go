package configs

import (
	"Kudak/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDatabase() error {
	config, err := InitConfigs()
	if err != nil {
		return err
	}

	dbUrl := fmt.Sprintf(
		"host = %s port = %s user = %s password = %s dbname = %s",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.DBname)

	DB, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Kindergarten{})
	DB.AutoMigrate(&models.KindergartenPicture{})
	DB.AutoMigrate(&models.EducationMinistry{})
	DB.AutoMigrate(&models.MainDepartment{})

	return nil
}
