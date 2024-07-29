package service

import (
	"Kudak/internal/repository"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUsernameTaken = errors.New("Пользователь уже существует")

type Service struct {
	Repository *repository.Repository
}

func (s *Service) CreateUser(login, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.Repository.GetUserByUsername(login)
	if err == nil {
		return ErrUsernameTaken
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	err = s.Repository.CreateUser(login, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Authenticiate(username, password string) (bool, error) {
	user, err := s.Repository.GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, err
	}

	return true, nil
}
