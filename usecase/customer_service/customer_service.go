package customer_service

import (
	"errors"
	"laotop_final/models"
	"laotop_final/repositories"
	"time"
)

type CustomerService interface {
	GetCustomer() ([]DataCustomerResponse, error)
	CreateCustomer(req CreateDataCustomerReq) (uint, error)
	UpdateCustomer(req UpdateDataCustomerReq) error
	DeleteCustomer(req DeleteDataCustomerReq) error
	GetCustomerByID(req DataCustomerByIdReq) ([]DataCustomerResponse, error)
}
type customerService struct {
	customerRepoSitory repositories.CustomerRepository
}

func (t *customerService) GetCustomer() ([]DataCustomerResponse, error) {
	resRepo, err := t.customerRepoSitory.GetDataCustomer()
	if err != nil {
		return nil, err
	}
	var res []DataCustomerResponse
	for _, req := range resRepo {
		res = append(res, DataCustomerResponse{
			ID:            req.ID,
			TableNumber:   req.TableNumber,
			Adults:        req.Adults,
			Children:      req.Children,
			TotalCustomer: req.Adults + req.Children,
			Status:        req.Status,
			CanOrder:      req.CanOrder,
			CreatedAt:     req.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     req.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if len(res) == 0 {
		return []DataCustomerResponse{}, nil
	}
	return res, nil
}
func (t *customerService) GetCustomerByID(req DataCustomerByIdReq) ([]DataCustomerResponse, error) {
	data, err := t.customerRepoSitory.GetDataCustomerByID(req.ID)
	if err != nil {
		return nil, err
	}
	var res []DataCustomerResponse

	if data == nil {
		return []DataCustomerResponse{}, nil
	}

	res = append(res, DataCustomerResponse{
		ID:            data.ID,
		TableNumber:   data.TableNumber,
		Adults:        data.Adults,
		Children:      data.Children,
		TotalCustomer: data.Adults + data.Children,
		Status:        data.Status,
		CanOrder:      data.CanOrder,
		CreatedAt:     data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     data.UpdatedAt.Format("2006-01-02 15:04:05"),
	})

	if len(res) == 0 {
		return []DataCustomerResponse{}, nil
	}
	return res, nil
}

func (t *customerService) CreateCustomer(req CreateDataCustomerReq) (uint, error) {
	isActive, err := t.customerRepoSitory.IsTableNumberActive(req.TableNumber)
	if err != nil {
		return 0, err
	}
	if isActive {
		return 0, errors.New("table_number is already in use")
	}

	customer := models.DataCustomer{
		TableNumber:   req.TableNumber,
		Adults:        req.Adults,
		Children:      req.Children,
		TotalCustomer: req.Adults + req.Children,
		Status:        true,
		CreatedAt:     time.Now(),
	}

	id, err := t.customerRepoSitory.CreateCustomer(customer)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (t *customerService) UpdateCustomer(req UpdateDataCustomerReq) error {
	existing, err := t.customerRepoSitory.GetDataCustomerByID(int(req.ID))
	if err != nil {
		return err
	}
	err = t.customerRepoSitory.UpdateCustomer(models.DataCustomer{
		ID:            req.ID,
		TableNumber:   req.TableNumber,
		Adults:        req.Adults,
		Children:      req.Children,
		TotalCustomer: req.Adults + req.Children,
		Status:        existing.Status,
		UpdatedAt:     time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (t *customerService) DeleteCustomer(req DeleteDataCustomerReq) error {
	err := t.customerRepoSitory.DeleteCustomer(models.DataCustomer{
		ID: req.ID,
	})
	if err == nil {
		return err
	}
	return nil
}
func NewCustomerService(
	customerRepoSitory *repositories.CustomerRepository,
) CustomerService {
	return &customerService{
		customerRepoSitory: *customerRepoSitory,
	}
}
