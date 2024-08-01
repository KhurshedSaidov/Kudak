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

func (r *Repository) AddKindergarten(kindergarten *models.Kindergarten) error {
	return r.DB.Create(kindergarten).Error
}

func (r *Repository) GetAllKindergartens() (models.Kindergarten, error) {
	var kindergartens models.Kindergarten
	err := r.DB.Select("id", "name", "inn", "address", "number", "description").Find(&kindergartens).Error
	return kindergartens, err
}

func (r *Repository) GetKindergartenByID(id uint) (models.Kindergarten, error) {
	var kindergarten models.Kindergarten
	if err := r.DB.Select("name", "picture").First(&kindergarten, id).Error; err != nil {
		return kindergarten, err
	}
	return kindergarten, nil
}
func (r *Repository) CreateEducationMinistry(em *models.EducationMinistry) error {
	return r.DB.Create(em).Error
}

func (r *Repository) GetAllEducationMinistries() ([]models.EducationMinistry, error) {
	var ministries []models.EducationMinistry
	err := r.DB.Select("id, name, inn, number").Find(&ministries).Error
	return ministries, err
}

func (r *Repository) GetEducationMinistryByID(id uint) (models.EducationMinistry, error) {
	var em models.EducationMinistry
	err := r.DB.Select("id, name, inn, number").First(&em, id).Error
	return em, err
}

func (r *Repository) UpdateEducationMinistry(em *models.EducationMinistry) error {
	return r.DB.Save(em).Error
}

func (r *Repository) DeleteEducationMinistry(id uint) error {
	return r.DB.Delete(&models.EducationMinistry{}, id).Error
}

func (r *Repository) ArchiveEducationMinistry(oldEm models.EducationMinistry) error {
	archive := models.EducationMinistryArchive{
		OriginalID:  oldEm.ID,
		Name:        oldEm.Name,
		Password:    oldEm.Password,
		Inn:         oldEm.Inn,
		PhoneNumber: oldEm.PhoneNumber,
	}
	return r.DB.Create(&archive).Error
}

func (r *Repository) CreateMainDepartment(md *models.MainDepartment) error {
	return r.DB.Create(md).Error
}

func (r *Repository) GetAllMainDepartments() ([]models.MainDepartment, error) {
	var departments []models.MainDepartment
	err := r.DB.Select("id, name, inn, number").Find(&departments).Error
	return departments, err
}

func (r *Repository) GetMainDepartmentByID(id uint) (models.MainDepartment, error) {
	var md models.MainDepartment
	err := r.DB.Select("id, name, inn, number").First(&md, id).Error
	return md, err
}

func (r *Repository) UpdateMainDepartment(md *models.MainDepartment) error {
	return r.DB.Save(md).Error
}

func (r *Repository) ArchiveMainDepartment(oldDm models.MainDepartment) error {
	archive := models.MainDepartmentArchive{
		OriginalID:  oldDm.ID,
		Name:        oldDm.Name,
		Password:    oldDm.Password,
		Inn:         oldDm.Inn,
		PhoneNumber: oldDm.PhoneNumber,
	}
	return r.DB.Create(&archive).Error
}

func (r *Repository) DeleteMainDepartment(id uint) error {
	return r.DB.Delete(&models.MainDepartment{}, id).Error
}
