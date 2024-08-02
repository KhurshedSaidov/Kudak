package service

import (
	"Kudak/internal/repository"
	"Kudak/models"
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

func (s *Service) CreateKindergarten(file *models.Kindergarten) error {
	return s.Repository.CreateKindergarten(file)
}

func (s *Service) GetAllKindergartens() ([]models.Kindergarten, error) {
	return s.Repository.GetAllKindergartens()
}

func (s *Service) GetKindergartenByID(id uint) (models.Kindergarten, error) {
	return s.Repository.GetKindergartenByID(id)
}

func (s *Service) DeleteKindergartenByID(id uint) error {
	return s.Repository.DeleteKindergartenByID(id)
}

func (s *Service) CreateEducationMinistry(em *models.EducationMinistry) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(em.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	em.Password = string(hashedPassword)
	return s.Repository.CreateEducationMinistry(em)
}

func (s *Service) GetAllEducationMinistries() ([]models.EducationMinistry, error) {
	return s.Repository.GetAllEducationMinistries()
}

func (s *Service) GetEducationMinistryByID(id uint) (models.EducationMinistry, error) {
	return s.Repository.GetEducationMinistryByID(id)
}

func (s *Service) UpdateEducationMinistry(em *models.EducationMinistry) error {

	oldEm, err := s.Repository.GetEducationMinistryByID(em.ID)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(em.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	em.Password = string(hashedPassword)

	if err := s.Repository.ArchiveEducationMinistry(oldEm); err != nil {
		return err
	}

	return s.Repository.UpdateEducationMinistry(em)
}

func (s *Service) DeleteEducationMinistry(id uint) error {
	oldEm, err := s.Repository.GetEducationMinistryByID(id)
	if err != nil {
		return err
	}

	if err := s.Repository.ArchiveEducationMinistry(oldEm); err != nil {
		return err
	}

	return s.Repository.DeleteEducationMinistry(id)
}

func (s *Service) CreateMainDepartment(md *models.MainDepartment) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(md.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	md.Password = string(hashedPassword)
	return s.Repository.CreateMainDepartment(md)
}

func (s *Service) GetAllMainDepartments() ([]models.MainDepartment, error) {
	return s.Repository.GetAllMainDepartments()
}

func (s *Service) GetMainDepartmentByID(id uint) (models.MainDepartment, error) {
	return s.Repository.GetMainDepartmentByID(id)
}

func (s *Service) UpdateMainDepartment(md *models.MainDepartment) error {
	oldMd, err := s.Repository.GetMainDepartmentByID(md.ID)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(md.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	md.Password = string(hashedPassword)

	if err := s.Repository.ArchiveMainDepartment(oldMd); err != nil {
		return err
	}
	return s.Repository.UpdateMainDepartment(md)
}

func (s *Service) DeleteMainDepartment(id uint) error {

	oldMd, err := s.Repository.GetMainDepartmentByID(id)
	if err != nil {
		return err
	}

	if err := s.Repository.ArchiveMainDepartment(oldMd); err != nil {
		return err
	}
	return s.Repository.DeleteMainDepartment(id)
}

func (s *Service) AddChild(kindergartenID uint, child *models.Child) error {
	return s.Repository.AddChild(kindergartenID, child)
}

func (s *Service) RecordAttendance(childID uint, present bool) error {
	return s.Repository.RecordAttendance(childID, present)
}

func (s *Service) GetAllChildren() ([]models.Child, error) {
	return s.Repository.GetAllChildren()
}

func (s *Service) GetLatestAttendance() ([]models.Attendance, error) {
	return s.Repository.GetLatestAttendance()
}

func (s *Service) UpdateAttendance(kindergartenID uint, attendance *models.Attendance) error {
	return s.Repository.UpdateAttendance(kindergartenID, attendance)
}
