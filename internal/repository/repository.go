package repository

import (
	"Kudak/models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateUser(username, hashedPassword string) error {
	newUser := &models.User{
		Username: username,
		Password: hashedPassword,
	}

	return r.DB.Create(newUser).Error
}

func (r *Repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
