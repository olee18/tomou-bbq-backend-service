package repositories

import (
	"errors"
	"gorm.io/gorm"
	"laotop_final/models"
)

type CustomerRepository interface {
	GetDataCustomerByID(id int) (*models.DataCustomer, error)
	GetDataCustomer() ([]models.DataCustomer, error)
	CreateCustomer(req models.DataCustomer) (uint, error)
	UpdateCustomer(req models.DataCustomer) error
	DeleteCustomer(req models.DataCustomer) error
	IsTableNumberActive(tableNumber string) (bool, error)
}
type customerRepository struct {
	db *gorm.DB
}

func (t *customerRepository) GetDataCustomer() ([]models.DataCustomer, error) {
	var res []models.DataCustomer
	err := t.db.Order("created_at DESC").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *customerRepository) CreateCustomer(req models.DataCustomer) (uint, error) {
	err := t.db.Create(&req).Error
	if err != nil {
		return 0, err
	}
	return req.ID, nil
}

func (t *customerRepository) UpdateCustomer(req models.DataCustomer) error {
	tx := t.db.Model(&models.DataCustomer{}).Where("id = ?", req.ID).Updates(req)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("customer ID not found")
	}
	return nil
}

func (t *customerRepository) DeleteCustomer(req models.DataCustomer) error {
	tx := t.db.Model(&models.DataCustomer{}).Where("id = ?", req.ID).Delete(req)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("customer ID not found")
	}
	return nil
}

func (t *customerRepository) IsTableNumberActive(tableNumber string) (bool, error) {
	var count int64
	err := t.db.Model(&models.DataCustomer{}).
		Where("table_number = ? AND status = ?", tableNumber, true).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (t *customerRepository) GetDataCustomerByID(id int) (*models.DataCustomer, error) {
	var customer models.DataCustomer
	err := t.db.
		Where("id = ?", id).
		Order("created_at DESC").
		First(&customer).Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}
