package user_service

// ===================== Request =======================//

type ResetPassword struct {
	ID int `json:"id" validate:"required"`
}

//type Create struct {
//	Name         string `json:"name" validate:"required"`
//	Username     string `json:"username" validate:"required"`
//	AccessStatus bool   `json:"access_status"`
//	Password     string `json:"password" validate:"required"`
//	RoleID       int    `json:"role_id" validate:"required"`
//}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	RoleName string `json:"role_name" validate:"required"`
}

type SignIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type ChangePassword struct {
	ID              int    `json:"id"`
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}
type Update struct {
	ID           int   `json:"id" validate:"required"`
	RoleID       int   `json:"role_id" validate:"required"`
	AccessStatus *bool `json:"access_status"`
}
type PermissionCreate struct {
	Name    string `json:"name"  validate:"required"`
	Keyword string `json:"keyword"  validate:"required"`
	Sort    string `json:"sort"  validate:"required"`
}

type RoleCreate struct {
	Name       string `json:"name"  validate:"required"`
	Permission []int  `json:"permission" validate:"required"`
}
type RoleUpdate struct {
	ID         int   `json:"id" validate:"required"`
	Permission []int `json:"permission" validate:"required"`
}
type RoleDelete struct {
	ID int `json:"id" validate:"required"`
}
type UserDelect struct {
	ID int `json:"id" validate:"required"`
}
type IssuerMerchantReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//===================== Responses =======================//

type LoginResponse struct {
	Token string `json:"token"`
}

type ResponseLogin struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	LoginAt  string `json:"login_at"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}

type Response struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	AccessStatus bool   `json:"access_status"`
	Role         string `json:"role"`
	RoleID       int    `json:"role_id"`
	Profile      string `json:"profile"`
	LoginAt      string `json:"login_at"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ResponseInfos struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	AccessStatus bool   `json:"access_status"`
	Role         Role   `json:"role"`
	Profile      string `json:"profile"`
	LoginAt      string `json:"login_at"`
	CreatedAt    string `json:"created_at"`
	UpdateAt     string `json:"update_at"`
}

type Role struct {
	Name       string   `json:"name"`
	Permission []string `json:"permission"`
}

type RoleResponse struct {
	ID         int             `json:"id"`
	Name       string          `json:"name"`
	Permission []PermissionRes `json:"permission"`
	CreatedAt  string          `json:"created_at"`
	UpdateAt   string          `json:"update_at"`
}
type PermissionRes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type ResponsesResetPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type SupplyResponseInfo struct {
	Name         string `json:"name"`
	Username     string `json:"username"`
	AccessStatus bool   `json:"access_status"`
	AccessToken  string `json:"access_token"`
}
