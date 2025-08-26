package repositories

import (
	"errors"
	"gorm.io/gorm"
	"laotop_final/models"
)

type OrderRePoSitroy interface {
	GetOrder() ([]models.Order, error)
	GetOrderWithItemsByCustomerID(customerID uint) ([]models.Order, error)
	CreateOrder(order models.Order, items []models.OrderItem) error
	UpdateOrder(req models.Order) error
	DeleteOrder(req models.Order) error
	CreateOrderItems(items []models.OrderItem) error
	UpdateOrderItemStatus(orderItem models.Order) error
	GetOrderWithItems() ([]models.Order, error)
	GetCategoryByID(ID uint) (*models.Category, error)
	GetCategory() ([]models.Category, error)
	CreateCategory(req models.Category) error
	UpdateCategory(req models.Category) error
	DeleteCategory(req models.Category) error
	GetMenuByCategoryID(categoryID uint) ([]models.Menu, error)
	GetMenuByID(id uint) (*models.Menu, error)
	GetMenu() ([]models.Menu, error)
	CreateMenu(req models.Menu) error
	UpdateMenu(req models.Menu) error
	DeleteMenu(req models.Menu) error
	GetAllOrderItems() ([]models.OrderItem, error)
	GetAllOrdersWithItems() ([]models.Order, error)
	DeleteOrderItemByCustomerAndMenu(customerID uint, menuID uint) error
}
type orderRepository struct {
	db *gorm.DB
}

func (o *orderRepository) UpdateOrderItemStatus(orderItem models.Order) error {
	tx := o.db.Model(&models.Order{}).Where("id = ?", orderItem.ID)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("order item ID not found")
	}
	return nil
}
func (o *orderRepository) GetAllOrderItems() ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := o.db.
		Preload("Category").
		Preload("Order").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (o *orderRepository) CreateOrder(order models.Order, items []models.OrderItem) error {
	tx := o.db.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range items {
		items[i].OrderID = order.ID
		if items[i].Order.OrderStatus == "" {
			items[i].Order.OrderStatus = "PENDING"
		}
		if err := tx.Create(&items[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
func (o *orderRepository) GetOrderWithItems() ([]models.Order, error) {
	var orders []models.Order
	err := o.db.Preload("OrderItems.Category").Find(&orders).Error
	return orders, err
}

func (o *orderRepository) GetAllOrdersWithItems() ([]models.Order, error) {
	var orders []models.Order
	err := o.db.Preload("OrderItems.Category").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) CreateOrderItems(items []models.OrderItem) error {
	return r.db.Create(&items).Error
}

func (o *orderRepository) GetOrderWithItemsByCustomerID(customerID uint) ([]models.Order, error) {
	var orders []models.Order
	err := o.db.Where("customer_id = ?", customerID).
		Preload("OrderItems.Category").
		Find(&orders).Error
	return orders, err
}

func (o *orderRepository) GetMenuByID(id uint) (*models.Menu, error) {
	var menu models.Menu
	err := o.db.Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (o *orderRepository) GetOrder() ([]models.Order, error) {
	var orders []models.Order
	err := o.db.Preload("Category").
		Joins("JOIN data_customers ON data_customers.id = orders.customer_id").
		Where("orders.menu_status = ?", "PENDING").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
func (o *orderRepository) UpdateOrder(req models.Order) error {
	tx := o.db.Model(&models.Order{}).Where("id = ?", req.ID).Updates(&req)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("order ID not found")
	}
	return nil
}
func (o *orderRepository) DeleteOrder(req models.Order) error {
	tx := o.db.Model(&models.Order{}).Where("id = ?", req.ID).Delete(&req)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("order ID not found")
	}
	return nil
}

func (o *orderRepository) GetCategoryByID(ID uint) (*models.Category, error) {
	var category models.Category
	err := o.db.Where("id = ?", ID).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (o *orderRepository) GetCategory() (res []models.Category, err error) {
	err = o.db.Find(&res).Error
	if err == nil {
		return res, nil
	}
	return res, nil
}

func (o *orderRepository) CreateCategory(req models.Category) error {
	err := o.db.Create(&req).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *orderRepository) UpdateCategory(req models.Category) error {
	tx := o.db.Model(&models.Category{}).Where("id = ?", req.ID).Updates(&req)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New(" ID not found")
	}
	return nil
}

func (o *orderRepository) DeleteCategory(req models.Category) error {
	tx := o.db.Model(&models.Category{}).Where("id = ?", req.ID).Delete(&req)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("ID not found")
	}
	return nil
}
func (o *orderRepository) DeleteOrderItemByCustomerAndMenu(customerID, menuID uint) error {
	subQuery := o.db.Model(&models.Order{}).Select("id").Where("customer_id = ?", customerID)
	tx := o.db.Where("order_id IN (?) AND menu_id = ?", subQuery, menuID).Delete(&models.OrderItem{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("order item not found")
	}
	return nil
}

func (o *orderRepository) GetMenu() ([]models.Menu, error) {
	var menus []models.Menu
	err := o.db.Preload("Category").Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (o *orderRepository) GetMenuByCategoryID(categoryID uint) ([]models.Menu, error) {
	var menus []models.Menu
	err := o.db.Where("category_id = ?", categoryID).Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func (o *orderRepository) CreateMenu(req models.Menu) error {
	err := o.db.Create(&req).Error
	if err != nil {
		return err
	}
	return nil
}

func (o *orderRepository) UpdateMenu(req models.Menu) error {
	tx := o.db.Model(&models.Menu{}).
		Where("id = ?", req.ID).
		Select("Name", "Image", "Price", "Status", "CategoryID").
		Updates(req)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("ID not found")
	}
	return nil
}

func (o *orderRepository) DeleteMenu(req models.Menu) error {
	tx := o.db.Model(&models.Menu{}).Where("id = ?", req.ID).Delete(&req)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New(" ID not found")
	}
	return nil
}
func NewOrderRepository(db *gorm.DB) OrderRePoSitroy {

	return &orderRepository{db: db}
}
