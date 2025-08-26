package order_service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"laotop_final/models"
	"laotop_final/repositories"
	"strings"
)

type OrderService interface {
	GetOrderHistory() ([]OrderResponse, error)
	GetOrderByCustomerID(req OrderReq) ([]OrderResponse, error)
	CreateOrder(req CreateOrder) (uint, error)
	UpdateOrder(req UpdateOrderReq) error
	DeleteOrder(req DeleteOrderReq) error
	UpdateOrderStatus(req OrderItemUpdateStatusReq) error
	CreateCategory(req CreateCategoryReq) error
	UpdateCategory(req UpdateCategoryReq) error
	DeleteCategory(req DeleteCategoryReq) error

	GetCategoryByID(id int) (CategoryResponse, error)
	GetCategory() ([]CategoryResponse, error)
	GetMenuByCategoryID(categoryID uint) ([]MenuResponse, error)
	GetAllOrderItems() ([]OrderItemResponse, error)
	GetMenuByID(id int) (MenuResponse, error)
	GetMenu() ([]MenuResponse, error)
	CreateMenu(req CreateMenuReq) error
	UpdateMenu(req UpdateMenuReq) error
	DeleteMenu(req DeleteMenuReq) error
	DeleteOrderItem(req DeleteOrderItemReq) error
	GetAllOrdersNoFilter() ([]OrderResponse, error)
}

type orderService struct {
	db                 *gorm.DB
	orderRepository    repositories.OrderRePoSitroy
	customerRepoSitory repositories.CustomerRepository
}

func (o *orderService) UpdateOrderStatus(req OrderItemUpdateStatusReq) error {
	rderItem := models.Order{
		ID:          req.OrderItemID,
		OrderStatus: req.OrderStatus,
	}
	return o.orderRepository.UpdateOrderItemStatus(rderItem)
}
func (o *orderService) GetAllOrderItems() ([]OrderItemResponse, error) {
	items, err := o.orderRepository.GetAllOrderItems()
	if err != nil {
		return nil, err
	}

	var response []OrderItemResponse
	for _, item := range items {
		response = append(response, OrderItemResponse{
			ID:           item.ID,
			CustomerID:   item.Order.CustomerID,
			TableNumber:  item.Order.TableNumber,
			MenuID:       item.MenuID,
			Name:         item.Name,
			Image:        item.Image,
			CategoryID:   item.CategoryID,
			Price:        item.Price,
			Quantity:     item.Quantity,
			CategoryName: item.Category.Name,
		})
	}
	if len(response) == 0 {
		return []OrderItemResponse{}, nil
	}
	return response, nil
}
func (o *orderService) GetOrderByCustomerID(req OrderReq) ([]OrderResponse, error) {
	orders, err := o.orderRepository.GetOrderWithItemsByCustomerID(req.CustomerID)
	if err != nil {
		return nil, err
	}

	var res []OrderResponse
	for _, order := range orders {
		var orderItemsResp []OrderItemResponse
		for _, item := range order.OrderItems {
			categoryName := ""
			if item.Category != nil {
				categoryName = item.Category.Name
			}

			orderItemsResp = append(orderItemsResp, OrderItemResponse{
				ID:           item.ID,
				MenuID:       item.MenuID,
				Name:         item.Name,
				Image:        item.Image,
				CategoryID:   item.CategoryID,
				CategoryName: categoryName,
				Price:        item.Price,
				Quantity:     item.Quantity,
			})
		}

		res = append(res, OrderResponse{
			ID:          order.ID,
			CustomerID:  order.CustomerID,
			OrderStatus: order.OrderStatus,
			TableNumber: order.TableNumber,
			CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
			OrderItems:  orderItemsResp,
		})
	}
	if len(res) == 0 {
		return []OrderResponse{}, nil
	}
	return res, nil

}

func (o *orderService) CreateOrder(req CreateOrder) (uint, error) {
	customer, err := o.customerRepoSitory.GetDataCustomerByID(int(*req.CustomerID))
	if err != nil {
		return 0, fmt.Errorf("Customer ID Not Found %d", *req.CustomerID)
	}
	if !customer.CanOrder {
		return 0, fmt.Errorf("This Customer has already checked out")
	}

	var order models.Order
	var existingOrder models.Order
	err = o.db.Where("customer_id = ? AND order_status = ?", *req.CustomerID, "PENDING").First(&existingOrder).Error
	if err == nil {
		order = existingOrder
	} else if err == gorm.ErrRecordNotFound {
		order = models.Order{
			CustomerID:  *req.CustomerID,
			TableNumber: customer.TableNumber,
			OrderStatus: "PENDING",
		}
		if err := o.db.Create(&order).Error; err != nil {
			return 0, fmt.Errorf("cannot create order: %w", err)
		}
	} else {
		return 0, err
	}

	var items []models.OrderItem
	for _, itemReq := range req.Items {
		var menu models.Menu
		if err := o.db.First(&menu, *itemReq.MenuID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return 0, fmt.Errorf("Menu ID %d not found", *itemReq.MenuID)
			}
			return 0, err
		}

		if !menu.Status {
			return 0, fmt.Errorf("Menu '%s' is currently unavailable for ordering", menu.Name)
		}
		items = append(items, models.OrderItem{
			OrderID:    order.ID,
			MenuID:     *itemReq.MenuID,
			Name:       menu.Name,
			Image:      menu.Image,
			CategoryID: menu.CategoryID,
			Price:      menu.Price,
			Quantity:   itemReq.Quantity,
		})
	}

	if err := o.orderRepository.CreateOrderItems(items); err != nil {
		return 0, fmt.Errorf("cannot create order items: %w", err)
	}

	return order.ID, nil
}

func (o *orderService) UpdateOrder(req UpdateOrderReq) error {
	err := o.orderRepository.UpdateOrder(models.Order{
		ID:          req.ID,
		OrderStatus: "SUCCESS",
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *orderService) GetMenuByID(id int) (MenuResponse, error) {
	req, err := o.orderRepository.GetCategoryByID(uint(id))
	if err != nil {
		return MenuResponse{}, err
	}
	Response := MenuResponse{
		ID:         req.ID,
		CategoryID: req.ID,
		Name:       req.Name,
		CreatedAt:  req.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  req.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return Response, nil
}

func (o *orderService) DeleteOrder(req DeleteOrderReq) error {
	err := o.orderRepository.DeleteOrder(models.Order{
		ID: req.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *orderService) GetOrderHistory() ([]OrderResponse, error) {
	orders, err := o.orderRepository.GetOrderWithItems()
	if err != nil {
		return nil, err
	}

	var responses []OrderResponse

	for _, order := range orders {
		if strings.ToUpper(strings.TrimSpace(order.OrderStatus)) != "PENDING" {
			continue
		}

		var customer models.DataCustomer
		if err := o.db.First(&customer, order.CustomerID).Error; err != nil {
			return nil, fmt.Errorf("customer ID %d not found: %w", order.CustomerID, err)
		}

		orderItems := make([]OrderItemResponse, 0, len(order.OrderItems))
		for _, item := range order.OrderItems {
			categoryName := ""
			if item.Category != nil {
				categoryName = item.Category.Name
			}
			orderItems = append(orderItems, OrderItemResponse{
				ID:           item.ID,
				MenuID:       item.MenuID,
				Name:         item.Name,
				Image:        item.Image,
				CategoryID:   item.CategoryID,
				CategoryName: categoryName,
				Price:        item.Price,
				Quantity:     item.Quantity,
			})
		}

		if len(orderItems) == 0 {
			continue
		}

		responses = append(responses, OrderResponse{
			ID:          order.ID,
			CustomerID:  order.CustomerID,
			OrderStatus: order.OrderStatus,
			TableNumber: customer.TableNumber,
			CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
			OrderItems:  orderItems,
		})
	}
	if len(responses) == 0 {
		return []OrderResponse{}, nil
	}

	return responses, nil
}

func (o *orderService) CreateCategory(req CreateCategoryReq) error {
	err := o.orderRepository.CreateCategory(models.Category{
		Name: req.Name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *orderService) UpdateCategory(req UpdateCategoryReq) error {
	err := o.orderRepository.UpdateCategory(models.Category{
		ID:   req.ID,
		Name: req.Name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *orderService) DeleteCategory(req DeleteCategoryReq) error {
	err := o.orderRepository.DeleteCategory(models.Category{
		ID: req.ID,
	})
	if err != nil {
		return err
	}
	return nil
}
func (o *orderService) DeleteOrderItem(req DeleteOrderItemReq) error {
	return o.orderRepository.DeleteOrderItemByCustomerAndMenu(req.CustomerID, req.MenuID)
}

func (o *orderService) GetCategoryByID(id int) (CategoryResponse, error) {
	category, err := o.orderRepository.GetCategoryByID(uint(id))
	if err != nil {
		return CategoryResponse{}, err
	}
	categoryResponse := CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return categoryResponse, nil

}
func (o *orderService) GetCategory() ([]CategoryResponse, error) {
	resReo, err := o.orderRepository.GetCategory()
	if err != nil {
		return nil, err
	}
	var res []CategoryResponse
	for _, req := range resReo {
		res = append(res, CategoryResponse{
			ID:        req.ID,
			Name:      req.Name,
			CreatedAt: req.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: req.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	if len(res) == 0 {
		return []CategoryResponse{}, nil
	}
	return res, nil

}
func (o *orderService) GetMenuByCategoryID(categoryID uint) ([]MenuResponse, error) {
	menus, err := o.orderRepository.GetMenuByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}
	var res []MenuResponse
	for _, menu := range menus {
		res = append(res, MenuResponse{
			ID:         menu.ID,
			CategoryID: menu.CategoryID,
			Name:       menu.Name,
			Image:      menu.Image,
			Price:      menu.Price,
			Status:     menu.Status,
			CreatedAt:  menu.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  menu.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return res, nil
}
func (o *orderService) GetMenu() ([]MenuResponse, error) {
	resRepo, err := o.orderRepository.GetMenu()
	if err != nil {
		return nil, err
	}

	var res []MenuResponse
	for _, req := range resRepo {
		res = append(res, MenuResponse{
			ID:           req.ID,
			CategoryID:   req.CategoryID,
			CategoryName: req.Category.Name,
			Image:        req.Image,
			Name:         req.Name,
			Price:        req.Price,
			Status:       req.Status,
			CreatedAt:    req.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    req.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if len(res) == 0 {
		return []MenuResponse{}, nil
	}
	return res, nil
}

func (o *orderService) CreateMenu(req CreateMenuReq) error {
	err := o.orderRepository.CreateMenu(models.Menu{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Image:      req.Image,
		Price:      req.Price,
		Status:     true,
	})
	if err != nil {

		return err
	}
	return nil
}

func (o *orderService) UpdateMenu(req UpdateMenuReq) error {
	var existing models.Menu
	if err := o.db.First(&existing, req.ID).Error; err != nil {
		return errors.New("menu not found")
	}

	if req.Image == "" {
		req.Image = existing.Image
	}

	return o.orderRepository.UpdateMenu(models.Menu{
		ID:         req.ID,
		Name:       req.Name,
		Image:      req.Image,
		Price:      req.Price,
		Status:     req.Status,
		CategoryID: req.CategoryID,
	})
}

func (o *orderService) DeleteMenu(req DeleteMenuReq) error {
	err := o.orderRepository.DeleteMenu(models.Menu{
		ID: req.ID,
	})
	if err == nil {
		return nil
	}
	return err
}
func (o *orderService) GetAllOrdersNoFilter() ([]OrderResponse, error) {
	orders, err := o.orderRepository.GetAllOrdersWithItems()
	if err != nil {
		return nil, err
	}

	customerMap := make(map[uint]*OrderResponse)

	for _, order := range orders {
		var customer models.DataCustomer
		if err := o.db.First(&customer, order.CustomerID).Error; err != nil {
			return nil, fmt.Errorf("customer ID %d not found: %w", order.CustomerID, err)
		}

		if _, exists := customerMap[order.CustomerID]; !exists {
			customerMap[order.CustomerID] = &OrderResponse{
				CustomerID:  order.CustomerID,
				OrderStatus: order.OrderStatus,
				TableNumber: customer.TableNumber,
				CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
				OrderItems:  []OrderItemResponse{},
			}
		}

		for _, item := range order.OrderItems {
			categoryName := ""
			if item.Category != nil {
				categoryName = item.Category.Name
			}

			customerMap[order.CustomerID].OrderItems = append(customerMap[order.CustomerID].OrderItems, OrderItemResponse{
				ID:           item.ID,
				CustomerID:   order.CustomerID,
				TableNumber:  customer.TableNumber,
				MenuID:       item.MenuID,
				Name:         item.Name,
				Image:        item.Image,
				CategoryID:   item.CategoryID,
				CategoryName: categoryName,
				Price:        item.Price,
				Quantity:     item.Quantity,
			})
		}
	}

	var responses []OrderResponse
	for _, response := range customerMap {
		responses = append(responses, *response)
	}
	if len(responses) == 0 {
		return []OrderResponse{}, nil
	}

	return responses, nil
}

func NewOrderService(db *gorm.DB,
	orderRepository *repositories.OrderRePoSitroy,
	customerRepoSitory *repositories.CustomerRepository,
) OrderService {
	return &orderService{
		orderRepository:    *orderRepository,
		customerRepoSitory: *customerRepoSitory,
		db:                 db,
	}
}
