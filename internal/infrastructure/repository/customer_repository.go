package repository

import (
	"github.com/gabrielcamargo/oficina-api/internal/domain/entity"
	domainrepo "github.com/gabrielcamargo/oficina-api/internal/domain/repository"
	"github.com/gabrielcamargo/oficina-api/internal/infrastructure/database/model"
	"gorm.io/gorm"
)

type GormCustomerRepository struct {
	db *gorm.DB
}

func NewGormCustomerRepository(db *gorm.DB) *GormCustomerRepository {
	return &GormCustomerRepository{db: db}
}

func (r *GormCustomerRepository) Create(customer *entity.Customer) error {
	return r.db.Create(model.CustomerModelFromDomain(customer)).Error
}

func (r *GormCustomerRepository) Update(customer *entity.Customer) error {
	return r.db.Save(model.CustomerModelFromDomain(customer)).Error
}

func (r *GormCustomerRepository) Delete(id string) error {
	return r.db.Delete(&model.CustomerModel{}, "id = ?", id).Error
}

func (r *GormCustomerRepository) FindByID(id string) (*entity.Customer, error) {
	var m model.CustomerModel
	err := r.db.First(&m, "id = ?", id).Error
	if isNotFound(err) {
		return nil, domainrepo.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *GormCustomerRepository) FindByDocument(document string) (*entity.Customer, error) {
	var m model.CustomerModel
	err := r.db.First(&m, "document = ?", document).Error
	if isNotFound(err) {
		return nil, domainrepo.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *GormCustomerRepository) FindAll() ([]*entity.Customer, error) {
	var models []model.CustomerModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	customers := make([]*entity.Customer, len(models))
	for i, m := range models {
		customers[i] = m.ToDomain()
	}
	return customers, nil
}
