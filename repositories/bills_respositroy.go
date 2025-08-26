package repositories

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"laotop_final/models"
)

type BillHistoryResult struct {
	OrderID       uint
	TableNumber   string
	CategoryID    uint
	CategoryName  string
	MenuID        uint
	MenuName      string
	TotalCustomer int
	Adults        int
	Children      int
	TotalAmount   int
	Date          time.Time
	Name          string
	DrinkPrice    int
	DrinkQuantity int
	CustomerID    uint
	TotalPrice    int
	Quantity      int
	Price         int
	ID            int
}

type BillRePositroy interface {
	GetBillHistoryByTableNumber(tableNumber string) ([]BillHistoryResult, error)
	GetBillHistory() ([]BillHistoryResult, error)
	InsertBill(bill *models.Bill) error
	InsertBillItems(items []models.BillItem) error
	GetBillItemsByBillID(orderID uint) ([]models.BillItem, error)
	UpdateBillStatus(billID uint, status bool) error
	DeleteBill(req models.Bill) error
	UpdateBills(req []models.Bill) error
	GetAllBillItems() ([]models.BillItem, error)
	DeleteAllByCustomerID(customerID uint) error
}

type BillRepoSitroy struct {
	Db *gorm.DB
}

func (r *BillRepoSitroy) GetBillByCustomerIDAndStatus(CustomerID int) ([]models.Bill, error) {
	var bills []models.Bill
	err := r.Db.Where("customer_id = ?", CustomerID).Find(&bills).Error
	return bills, err
}
func (b *BillRepoSitroy) GetAllBillItems() ([]models.BillItem, error) {
	var billItems []models.BillItem
	err := b.Db.Preload("Category").Preload("Order").Find(&billItems).Error
	if err != nil {
		return nil, err
	}
	return billItems, nil
}
func (b *BillRepoSitroy) GetBillHistory() ([]BillHistoryResult, error) {
	var bills []BillHistoryResult
	err := b.Db.Table("bills AS bd").
		Select(`
			bd.order_id,
			bd.customer_id,
			t.table_number,
			bd.total_customer,
			bd.adults,
			bd.children,
			bd.total_amount,
			bd.date,
			bd.id
		`).
		Joins("LEFT JOIN data_customers t ON bd.customer_id = t.id").
		Order("bd.date DESC").
		Scan(&bills).Error
	if err != nil {
		return nil, err
	}
	return bills, nil
}

func (b *BillRepoSitroy) GetBillHistoryByTableNumber(tableNumber string) ([]BillHistoryResult, error) {
	var bills []BillHistoryResult
	err := b.Db.Table("bills AS bd").
		Select(`
			bd.order_id,
			bd.customer_id,
			t.table_number,
			c.id AS category_id,
			c.name AS category_name,
			bd.menu_id,
			m.name AS menu_name,
			bd.total_customer,
			bd.adults,
			bd.children,
			bd.total_amount,
			bd.date,
bd.id
		`).
		Joins("LEFT JOIN data_customers t ON bd.customer_id = t.id").
		Joins("LEFT JOIN menus m ON bd.menu_id = m.id").
		Joins("LEFT JOIN categories c ON m.category_id = c.id").
		Where("t.table_number = ?", tableNumber).
		Order("bd.date DESC").
		Scan(&bills).Error
	if err != nil {
		return nil, err
	}
	return bills, nil
}

func (b *BillRepoSitroy) InsertBill(bill *models.Bill) error {
	return b.Db.Create(bill).Error
}

func (b *BillRepoSitroy) InsertBillItems(items []models.BillItem) error {
	return b.Db.Create(&items).Error
}
func (r *BillRepoSitroy) UpdateBillsWithTx(tx *gorm.DB, bills []models.Bill) error {
	for _, bill := range bills {
		result := tx.Model(&models.Bill{}).
			Where("id = ?", bill.ID).
			Select("total_amount", "adults", "children", "total_customer", "status_bills", "updated_at").
			Updates(bill)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("bill with ID %d not found", bill.ID)
		}
	}
	return nil
}

func (r *BillRepoSitroy) DeleteBillItemsByBillIDWithTx(tx *gorm.DB, billID uint) error {
	return tx.Where("bill_id = ?", billID).Delete(&models.BillItem{}).Error
}

func (r *BillRepoSitroy) InsertBillItemsWithTx(tx *gorm.DB, items []models.BillItem) error {
	return tx.Create(&items).Error
}
func (b *BillRepoSitroy) UpdateBills(bills []models.Bill) error {
	for _, bill := range bills {
		tx := b.Db.Model(&models.Bill{}).
			Where("id = ?", bill.ID).
			Select("total_amount", "adults", "children", "total_customer", "status_bills", "updated_at").
			Updates(bill)

		if tx.Error != nil {
			return tx.Error
		}
		if tx.RowsAffected == 0 {
			return fmt.Errorf("bill with ID %d not found", bill.ID)
		}
	}
	return nil
}
func (r *BillRepoSitroy) GetBillByID(id uint) (*models.Bill, error) {
	var bill models.Bill
	err := r.Db.Preload("BillItems").First(&bill, id).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}
func (b *BillRepoSitroy) GetBillItemsByBillID(orderID uint) ([]models.BillItem, error) {
	var billItems []models.BillItem
	err := b.Db.Table("bill_items").Joins("JOIN bills ON bills.id = bill_items.bill_id").
		Where("bills.order_id = ?", orderID).Find(&billItems).Error
	return billItems, err
}

func (b *BillRepoSitroy) UpdateCustomerBillStatus(customerID uint, canOrder bool) error {
	return b.Db.Table("data_customers").
		Where("id = ?", customerID).
		Updates(map[string]interface{}{
			"can_order": canOrder,
			"status":    false,
		}).Error
}
func (b *BillRepoSitroy) UpdateBillStatus(billID uint, status bool) error {
	return b.Db.Model(&models.Bill{}).Where("id = ?", billID).Update("status_bills", status).Error
}

func (b *BillRepoSitroy) DeleteBill(req models.Bill) error {
	tx := b.Db.Model(&models.Bill{}).Where("id = ?", req.ID).Delete(req)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("ID not found")
	}
	return nil
}
func (b *BillRepoSitroy) DeleteBillItems(req models.BillItem) error {
	tx := b.Db.Model(&models.BillItem{}).Where("id = ?", req.ID).Delete(req)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("ID not found")
	}
	return nil
}

func (b *BillRepoSitroy) DeleteBillItemsByBillID(billID uint) error {
	return b.Db.Where("bill_id = ?", billID).Delete(&models.BillItem{}).Error
}

func (b *BillRepoSitroy) DeleteAllByCustomerID(customerID uint) error {
	return b.Db.Transaction(func(tx *gorm.DB) error {
		var bills []models.Bill
		if err := tx.Where("customer_id = ?", customerID).Find(&bills).Error; err != nil {
			return err
		}
		for _, bill := range bills {
			if err := tx.Where("bill_id = ?", bill.ID).Delete(&models.BillItem{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("customer_id = ?", customerID).Delete(&models.Bill{}).Error; err != nil {
			return err
		}
		var orders []models.Order
		if err := tx.Where("customer_id = ?", customerID).Find(&orders).Error; err != nil {
			return err
		}
		for _, order := range orders {
			if err := tx.Where("order_id = ?", order.ID).Delete(&models.OrderItem{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("customer_id = ?", customerID).Delete(&models.Order{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func NewBillRePositroy(db *gorm.DB) *BillRepoSitroy {
	db.AutoMigrate(
	//&models.Bill{},
	//&models.BillItem{},
	)
	return &BillRepoSitroy{Db: db}
}
