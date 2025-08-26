package user_service

import (
	"errors"
	"fmt"
	"laotop_final/middlewares"
	"laotop_final/models"
	"laotop_final/repositories"
	"laotop_final/utils"
	"strconv"
	"time"
)

type UserService interface {
	GetPermission() ([]models.Permission, error)

	GetRole() ([]RoleResponse, error)
	CreateRole(req RoleCreate) error
	UpdateRole(req RoleUpdate) error
	DeleteRole(req RoleDelete) error
	DeleteUser(req UserDelect) error

	SignIn(req SignIn) (*ResponseLogin, string, error)
	GetUserInfos(req RoleDelete) (*ResponseInfos, error)
	ChangePassword(req ChangePassword) (string, error)
	SignOut(keyStr, key string) (string, error)

	GetUser() ([]Response, error)
	CreateUser(req CreateUserRequest) (string, error)
	UpdateUser(req Update) (string, error)
	ResetPassword(req ResetPassword) (*ResponsesResetPassword, error)
}

type userService struct {
	repositoryUser  repositories.UserRepository
	repositoryRedis repositories.RedisRepository
}

func (s *userService) GetPermission() (permission []models.Permission, err error) {
	permission, err = s.repositoryUser.GetPermission()
	if err != nil {
		return nil, err
	}
	return permission, nil
}
func (s *userService) GetRole() (role []RoleResponse, err error) {
	roles, err := s.repositoryUser.GetRole()
	if err != nil {
		return nil, err
	}
	for _, item := range roles {
		var mappedPermissions []PermissionRes
		for _, p := range item.Permission {
			mappedPermissions = append(mappedPermissions, PermissionRes{
				ID:   p.ID,
				Name: p.Name,
			})
		}
		role = append(role, RoleResponse{
			ID:         item.ID,
			Name:       item.Name,
			Permission: mappedPermissions,
			CreatedAt:  item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdateAt:   item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	if len(role) <= 0 {
		role = []RoleResponse{}
	}
	return role, nil
}
func (s *userService) CreateRole(req RoleCreate) error {
	permissions := make([]models.Permission, len(req.Permission))
	for i, permID := range req.Permission {
		permissions[i] = models.Permission{ID: permID}
	}
	err := s.repositoryUser.CreateRole(models.Role{
		Name:       req.Name,
		Permission: permissions,
		CreatedAt:  time.Now(),
	})
	if err != nil {
		return errors.New("name exist")
	}
	return nil
}
func (s *userService) UpdateRole(req RoleUpdate) error {
	if len(req.Permission) <= 0 {
		return errors.New("permission is empty")
	}
	err := s.repositoryUser.UpdateRole(req.ID, req.Permission)
	if err != nil {
		return err
	}
	user, err := s.repositoryUser.GetOneUserRole(strconv.Itoa(req.ID))
	if err != nil {
		return err
	}
	_ = s.repositoryRedis.DelDataRedis(strconv.Itoa(user.Role.ID))
	return nil
}
func (s *userService) DeleteRole(req RoleDelete) error {
	err := s.repositoryUser.DeleteRole(models.Role{ID: req.ID})
	if err != nil {
		return err
	}
	return nil
}
func (s *userService) DeleteUser(req UserDelect) error {
	err := s.repositoryUser.DeleteUser(models.User{ID: req.ID})
	if err != nil {
		return err
	}
	return nil
}
func (s *userService) SignIn(req SignIn) (user *ResponseLogin, token string, err error) {
	//get username
	users, err := s.repositoryUser.CheckUsername(models.User{
		Username: req.Username,
	})
	if err != nil {
		fmt.Println(err)
		return nil, "", errors.New("username not found")
	}
	//check password
	if !utils.CheckPasswordHash(req.Password, users.Password) {
		return nil, "", errors.New("password wrong")
	}
	//response
	user = &ResponseLogin{
		Name:     users.Name,
		Username: users.Username,
		LoginAt:  users.LoginAt.Format("2006-01-02 15:04:05"),
		CreateAt: users.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdateAt: users.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	//map permission to me []
	var permission []string
	for _, p := range users.Role.Permission {
		permission = append(permission, p.Keyword)
	}
	//generate token
	token, err = middlewares.GenerateTokenWeb(strconv.Itoa(users.ID), strconv.Itoa(users.Role.ID), permission)
	if err != nil {
		return nil, "", err
	}
	//set token to redis_service
	err = s.repositoryRedis.SetHasRedis(strconv.Itoa(users.Role.ID), strconv.Itoa(users.ID), token)
	if err != nil {
		return nil, "", err
	}
	if utils.PointerBoolToBool(users.AccessStatus) == false {
		return nil, "", errors.New("user is not access")
	}

	return user, token, nil
}

func (s *userService) SignOut(keyStr, key string) (string, error) {
	err := s.repositoryRedis.SetHasRedis(keyStr, key, "")
	if err != nil {
		return "", err
	}
	return "success", nil
}
func (s *userService) GetUserInfos(req RoleDelete) (user *ResponseInfos, err error) {
	getUser, err := s.repositoryUser.GetOneUser(models.User{ID: req.ID})
	if err != nil {
		return nil, err
	}
	userResponse := MapUserToResponseInfos(getUser)
	return userResponse, nil
}
func (s *userService) ChangePassword(req ChangePassword) (string, error) {
	if req.NewPassword != req.ConfirmPassword {
		return "", errors.New("new password not match")
	}
	if req.OldPassword == req.NewPassword {
		return "", errors.New("change new password")
	}
	users, err := s.repositoryUser.GetOneUser(models.User{
		ID: req.ID,
	})
	if err != nil {
		return "", errors.New("id not found")
	}
	//check password
	if !utils.CheckPasswordHash(req.OldPassword, users.Password) {
		return "", errors.New("password wrong")
	}
	password, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return "", err
	}
	err = s.repositoryUser.UpdateUser(models.User{
		ID:       req.ID,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	err = s.repositoryRedis.SetHasRedis(strconv.Itoa(users.Role.ID), strconv.Itoa(users.ID), "")
	if err != nil {
		return "", err
	}
	return "success", nil
}

func (s *userService) GetUser() (user []Response, err error) {
	users, err := s.repositoryUser.GetUser()
	for _, item := range users {
		user = append(user, Response{
			ID:           item.ID,
			Name:         item.Name,
			Username:     item.Username,
			AccessStatus: utils.PointerBoolToBool(item.AccessStatus),
			Role:         item.Role.Name,
			RoleID:       item.Role.ID,
			Profile:      "",
			LoginAt:      item.LoginAt.Format("2006-01-02 15:04:05"),
			CreatedAt:    item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) CreateUser(req CreateUserRequest) (string, error) {
	role, err := s.repositoryUser.GetRoleByName(req.RoleName)

	if err != nil {
		return "", fmt.Errorf("role not found: %v", err)
	}

	password, err := utils.HashPassword(req.Password)
	if err != nil {
		return "", err
	}

	accessStatus := true
	user := models.User{
		Name:         req.Name,
		Username:     req.Username,
		AccessStatus: &accessStatus,
		Password:     password,
		RoleID:       role.ID,
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.repositoryUser.CreateUser(user); err != nil {
		return "", err
	}

	return "success", nil
}

func (s *userService) UpdateUser(req Update) (string, error) {
	err := s.repositoryUser.UpdateUser(models.User{
		ID:           req.ID,
		RoleID:       req.RoleID,
		AccessStatus: req.AccessStatus,
	})
	if err != nil {
		return "", err
	}
	return "success", nil
}
func (s *userService) ResetPassword(req ResetPassword) (myRes *ResponsesResetPassword, err error) {
	getUser, err := s.repositoryUser.GetOneUser(models.User{ID: req.ID})
	if err != nil {
		return nil, err
	}
	ranPassword := utils.RandomNumber(8)

	password, err := utils.HashPassword(ranPassword)
	if err != nil {
		return nil, err
	}
	err = s.repositoryUser.UpdateUser(models.User{
		ID:       req.ID,
		Password: password,
	})

	if err != nil {
		return nil, err
	}
	_ = s.repositoryRedis.SetHasRedis(strconv.Itoa(getUser.Role.ID), strconv.Itoa(getUser.ID), "")
	myRes = &ResponsesResetPassword{
		Username: getUser.Username,
		Password: ranPassword,
	}
	return myRes, nil
}
func NewUserService(
	repositoryUser *repositories.UserRepository,
	repositoryRedis *repositories.RedisRepository,
	// repo
) UserService {
	return &userService{
		repositoryUser:  *repositoryUser,
		repositoryRedis: *repositoryRedis,
		//repo
	}
}
