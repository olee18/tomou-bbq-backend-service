package customer_service

type CreateDataCustomerReq struct {
	Adults      int    `json:"adults" validate:"required"`
	Children    int    `json:"children"`
	TableNumber string `json:"table_number" validate:"required"`
}

type UpdateDataCustomerReq struct {
	ID          uint   `json:"id" validate:"required"`
	TableNumber string `json:"table_number" `
	Adults      int    `json:"adults"`
	Children    int    `json:"children"`
}

type DeleteDataCustomerReq struct {
	ID uint `json:"id" validate:"required"`
}
type DataCustomerByIdReq struct {
	ID int `json:"id" validate:"required"`
}

type DataCustomerResponse struct {
	ID            uint   `json:"id"`
	TableNumber   string `json:"table_number"`
	Adults        int    `json:"adults"`
	Children      int    `json:"children"`
	TotalCustomer int    `json:"total_customer"`
	Status        bool   `json:"status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	CanOrder      bool   `json:"can_order"`
}

type CustomerResponseClient struct {
	TableNumber   string `json:"table_number"`
	Adults        int    `json:"adults"`
	Children      int    `json:"children"`
	TotalCustomer int    `json:"total_customer"`
}
