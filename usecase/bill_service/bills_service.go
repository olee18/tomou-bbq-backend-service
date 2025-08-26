package bill_service

import (
	"fmt"
	"gorm.io/gorm"
	"laotop_final/models"
	"laotop_final/repositories"
	"time"
)

type BillsService interface {
	GetBillHistoryByCustomerID(req BillsReq) ([]BillResponse, error)
	GetBillHistory(req BillsReq) ([]BillResponse, error)
	GenerateBillForCustomerID(req BillsReq) ([]BillResponse, error)
	DeleteBillService(req BillDeleteReq) error
	UpdateBIllService(req UpdateBillReq) error
	GetAllBillItemsService() ([]BillItemResponse, error)
	DeleteBillitemService(req BillitemDeleteReq) error
	DeleteAllByCustomerID(customerID uint) error
	ShowBillClientByCustomerID(customerID uint) ([]BillResponse, error)
}
type billsService struct {
	billRepoSitRoy     repositories.BillRepoSitroy
	customerRepoSitory repositories.CustomerRepository
	orderRepository    repositories.OrderRePoSitroy
	db                 *gorm.DB
}

func (b *billsService) ShowBillClientByCustomerID(customerID uint) ([]BillResponse, error) {
	bills, err := b.billRepoSitRoy.GetBillByCustomerIDAndStatus(int(customerID))
	if err != nil {
		return nil, err
	}
	if len(bills) == 0 {
		return []BillResponse{}, nil
	}
	customer, err := b.customerRepoSitory.GetDataCustomerByID(int(customerID))
	if err != nil {
		return nil, err
	}

	var responses []BillResponse

	for _, bill := range bills {
		billItems, err := b.billRepoSitRoy.GetBillItemsByBillID(bill.OrderID)
		if err != nil {
			return nil, err
		}

		var itemResponses []BillItemResponse
		var totalExtraPrice int
		for _, item := range billItems {
			itemResponses = append(itemResponses, BillItemResponse{
				Name:       item.MenuName,
				Quantity:   item.Quantity,
				Price:      item.Price,
				TotalPrice: item.TotalPrice,
			})
			totalExtraPrice += item.TotalPrice
		}

		response := BillResponse{
			ID:                  int(bill.ID),
			CustomerID:          bill.CustomerID,
			TableNumber:         customer.TableNumber,
			OrderID:             bill.OrderID,
			TotalAmount:         bill.TotalAmount,
			Adults:              bill.Adults,
			Children:            bill.Children,
			TotalCustomer:       bill.TotalCustomer,
			Date:                bill.Date.Format("2006-01-02 15:04:05"),
			Items:               itemResponses,
			TotalCombinedAmount: bill.TotalAmount + totalExtraPrice,
		}

		responses = append(responses, response)
	}

	return responses, nil
}
func (b *billsService) GenerateBillForCustomerID(req BillsReq) ([]BillResponse, error) {
	existingBill, err := b.billRepoSitRoy.GetBillByCustomerIDAndStatus(int(req.CustomerID))
	if err != nil {
		return nil, err
	}
	if len(existingBill) > 0 {
		return nil, fmt.Errorf("bill already generated for this table")
	}

	customer, err := b.customerRepoSitory.GetDataCustomerByID(int(req.CustomerID))
	if err != nil {
		return nil, err
	}

	orders, err := b.orderRepository.GetOrderWithItemsByCustomerID(uint(req.CustomerID))
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("Orders not found for table ID %d", req.CustomerID)
	}
	for _, order := range orders {
		if order.OrderStatus == "PENDING" {
			return nil, fmt.Errorf("cannot generate bill: order is still pending")
		}
	}
	var totalDrinkPrice, totalTissuePrice int

	customerPrice := (customer.Adults * 99000) + (customer.Children * 49000)
	now := time.Now()
	bill := models.Bill{
		OrderID:       orders[0].ID,
		CustomerID:    uint(req.CustomerID),
		Name:          "Tomou BBQ Bill",
		TotalAmount:   customerPrice,
		Adults:        customer.Adults,
		Children:      customer.Children,
		TotalCustomer: customer.Adults + customer.Children,
		Date:          now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := b.billRepoSitRoy.InsertBill(&bill); err != nil {
		return nil, err
	}

	var billItems []models.BillItem
	var responseItems []BillItemResponse

	for _, order := range orders {
		for _, item := range order.OrderItems {
			if item.Category == nil {
				continue
			}
			if item.Category.Name == "Food" || item.Category.Name == "Sauce" || item.Category.Name == "Seafood" {
				continue
			}

			itemTotal := item.Price * item.Quantity

			billItems = append(billItems, models.BillItem{
				BillID:     bill.ID,
				CategoryID: item.Category.ID,
				OrderID:    order.ID,
				MenuName:   item.Name,
				Quantity:   item.Quantity,
				Price:      item.Price,
				TotalPrice: itemTotal,
				CreatedAt:  now,
				UpdatedAt:  now,
			})

			responseItems = append(responseItems, BillItemResponse{
				Name:       item.Name,
				Quantity:   item.Quantity,
				Price:      item.Price,
				TotalPrice: itemTotal,
			})

			switch item.Category.Name {
			case "Tissue Paper":
				totalTissuePrice += itemTotal
			case "Drink":
				totalDrinkPrice += itemTotal
			}
		}
	}

	if len(billItems) > 0 {
		if err := b.billRepoSitRoy.InsertBillItems(billItems); err != nil {
			return nil, err
		}
	}

	if err := b.billRepoSitRoy.UpdateBillStatus(bill.ID, true); err != nil {
		return nil, err
	}

	if err := b.billRepoSitRoy.UpdateCustomerBillStatus(bill.CustomerID, false); err != nil {
		return nil, err
	}

	billResponse := BillResponse{
		ID:                  int(bill.ID),
		CustomerID:          bill.CustomerID,
		TableNumber:         customer.TableNumber,
		OrderID:             bill.OrderID,
		TotalAmount:         bill.TotalAmount,
		Adults:              bill.Adults,
		Children:            bill.Children,
		TotalCustomer:       bill.TotalCustomer,
		Date:                bill.Date.Format("2006-01-02 15:04:05"),
		Items:               responseItems,
		TotalCombinedAmount: customerPrice + totalDrinkPrice + totalTissuePrice,
	}

	return []BillResponse{billResponse}, nil
}
func (b *billsService) UpdateBIllService(req UpdateBillReq) error {
	const adultPrice = 99000
	const childPrice = 49000

	bill, err := b.billRepoSitRoy.GetBillByID(uint(req.ID))
	if err != nil {
		return fmt.Errorf("bill not found: %w", err)
	}

	orderID := bill.OrderID
	customerID := bill.CustomerID
	tableNumber := bill.TableNumber

	totalAdultPrice := req.Adults * adultPrice
	totalChildPrice := req.Children * childPrice
	totalAmount := totalAdultPrice + totalChildPrice

	var newItems []models.BillItem
	now := time.Now()

	for _, item := range req.Items {
		var menu models.Menu
		if err := b.db.Where("name = ?", item.Name).First(&menu).Error; err != nil {
			return fmt.Errorf("menu not found: %s", item.Name)
		}

		totalPrice := item.Quantity * menu.Price

		newItems = append(newItems, models.BillItem{
			BillID:      uint(req.ID),
			OrderID:     orderID,
			CustomerID:  customerID,
			TableNumber: tableNumber,
			CategoryID:  menu.CategoryID,
			MenuName:    menu.Name,
			Quantity:    item.Quantity,
			Price:       menu.Price,
			TotalPrice:  totalPrice,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	return b.db.Transaction(func(tx *gorm.DB) error {
		err := b.billRepoSitRoy.UpdateBillsWithTx(tx, []models.Bill{
			{
				ID:            uint(req.ID),
				TotalAmount:   totalAmount,
				Adults:        req.Adults,
				Children:      req.Children,
				TotalCustomer: req.Adults + req.Children,
				StatusBills:   req.StatusBills,
				UpdatedAt:     now,
			},
		})
		if err != nil {
			return err
		}

		if err := b.billRepoSitRoy.DeleteBillItemsByBillIDWithTx(tx, uint(req.ID)); err != nil {
			return fmt.Errorf("failed to delete old bill items: %w", err)
		}

		if len(newItems) > 0 {
			if err := b.billRepoSitRoy.InsertBillItemsWithTx(tx, newItems); err != nil {
				return fmt.Errorf("failed to insert new bill items: %w", err)
			}
		}

		return nil
	})
}

func (b *billsService) GetBillHistory(req BillsReq) ([]BillResponse, error) {
	var resRepo []repositories.BillHistoryResult
	var err error
	switch req.TableNumber {
	case "":
		resRepo, err = b.billRepoSitRoy.GetBillHistory()
	default:
		resRepo, err = b.billRepoSitRoy.GetBillHistoryByTableNumber(req.TableNumber)
	}
	if err != nil {
		return nil, err
	}
	var res []BillResponse
	for _, v := range resRepo {
		billItems, err := b.billRepoSitRoy.GetBillItemsByBillID(v.OrderID)
		if err != nil {
			return nil, err
		}
		var itemResponses []BillItemResponse
		var sumItemTotalPrice int
		for _, item := range billItems {
			sumItemTotalPrice += item.TotalPrice
			itemResponses = append(itemResponses, BillItemResponse{
				Name:       item.MenuName,
				Quantity:   item.Quantity,
				Price:      item.Price,
				TotalPrice: item.TotalPrice,
			})
		}
		if itemResponses == nil {
			itemResponses = []BillItemResponse{}
		}
		res = append(res, BillResponse{
			ID:                  v.ID,
			CustomerID:          v.CustomerID,
			TableNumber:         v.TableNumber,
			OrderID:             v.OrderID,
			TotalAmount:         v.TotalAmount,
			TotalCustomer:       v.TotalCustomer,
			Adults:              v.Adults,
			Children:            v.Children,
			Date:                v.Date.Format("2006-01-02 15:04:05"),
			Items:               itemResponses,
			TotalCombinedAmount: v.TotalAmount + sumItemTotalPrice,
		})
	}
	if len(res) == 0 {
		return []BillResponse{}, nil
	}
	return res, nil
}

func (b *billsService) GetBillHistoryByCustomerID(req BillsReq) ([]BillResponse, error) {
	resRepo, err := b.billRepoSitRoy.GetBillHistoryByTableNumber(req.TableNumber)
	if err != nil {
		return nil, err
	}
	var res []BillResponse
	for _, v := range resRepo {
		billItems, err := b.billRepoSitRoy.GetBillItemsByBillID(v.OrderID)
		if err != nil {
			return nil, err
		}
		var itemResponses []BillItemResponse
		var sumItemTotalPrice int
		for _, item := range billItems {
			sumItemTotalPrice += item.TotalPrice
			itemResponses = append(itemResponses, BillItemResponse{
				Name:       item.MenuName,
				Quantity:   item.Quantity,
				Price:      item.Price,
				TotalPrice: item.TotalPrice,
			})
		}
		res = append(res, BillResponse{
			CustomerID:          v.CustomerID,
			TableNumber:         v.TableNumber,
			OrderID:             v.OrderID,
			TotalAmount:         v.TotalAmount,
			TotalCustomer:       v.TotalCustomer,
			Adults:              v.Adults,
			Children:            v.Children,
			Date:                v.Date.Format("2006-01-02 15:04:05"),
			Items:               itemResponses,
			TotalCombinedAmount: v.TotalAmount + sumItemTotalPrice,
		})
	}
	if len(res) == 0 {
		return []BillResponse{}, nil
	}
	return res, nil
}

func (t *billsService) DeleteBillService(req BillDeleteReq) error {
	err := t.billRepoSitRoy.DeleteBillItemsByBillID(req.ID)
	if err != nil {
		return err
	}
	err = t.billRepoSitRoy.DeleteBill(models.Bill{
		ID: req.ID,
	})
	if err != nil {
		return err
	}
	return nil
}
func (t *billsService) DeleteBillitemService(req BillitemDeleteReq) error {
	err := t.billRepoSitRoy.DeleteBillItems(models.BillItem{ID: req.ID})
	if err != nil {
		return err
	}
	return nil
}

func (b *billsService) GetAllBillItemsService() ([]BillItemResponse, error) {
	billItems, err := b.billRepoSitRoy.GetAllBillItems()
	if err != nil {
		return nil, err
	}

	var response []BillItemResponse
	for _, item := range billItems {
		response = append(response, BillItemResponse{
			ID:         item.ID,
			Name:       item.MenuName,
			Quantity:   item.Quantity,
			Price:      item.Price,
			TotalPrice: item.TotalPrice,
		})
	}
	if len(response) == 0 {
		return []BillItemResponse{}, nil
	}
	return response, nil
}

func (b *billsService) DeleteAllByCustomerID(customerID uint) error {
	return b.billRepoSitRoy.DeleteAllByCustomerID(customerID)
}

func NewBillsService(
	billRepoSitroy *repositories.BillRepoSitroy,
	customerRepoSitory *repositories.CustomerRepository,
	orderRepository *repositories.OrderRePoSitroy,
	db *gorm.DB,
) BillsService {
	return &billsService{
		billRepoSitRoy:     *billRepoSitroy,
		customerRepoSitory: *customerRepoSitory,
		orderRepository:    *orderRepository,
		db:                 db,
	}
}
