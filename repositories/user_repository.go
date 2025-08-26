package repositories

import (
	"errors"
	"gorm.io/gorm"
	"laotop_final/models"
	"time"
)

type UserRepository interface {
	GetPermission() ([]models.Permission, error)
	GetUserRole(roleID string) ([]models.User, error)
	GetOneUserRole(roleID string) (*models.User, error)
	GetRole() ([]models.Role, error)
	CreateRole(req models.Role) error
	UpdateRole(roleID int, permissionID []int) error
	DeleteRole(req models.Role) error
	GetUser() ([]models.User, error)
	CheckUsername(req models.User) (*models.User, error)
	GetOneUser(req models.User) (*models.User, error)
	CreateUser(model models.User) error
	UpdateUser(req models.User) error
	GetRoleByName(roleName string) (models.Role, error)
	DeleteUser(req models.User) error
}

type userRepository struct{ db *gorm.DB }

func (r *userRepository) GetPermission() (permission []models.Permission, err error) {
	err = r.db.Order("sort ASC").Find(&permission).Error
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (r *userRepository) GetUserRole(roleID string) (user []models.User, err error) {
	err = r.db.Where("role_id=?", roleID).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetOneUserRole(roleID string) (user *models.User, err error) {
	err = r.db.Where("role_id=?", roleID).Preload("Role", func(db *gorm.DB) *gorm.DB {
		return db
	}).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetRole() (role []models.Role, err error) {
	err = r.db.
		Preload("Permission", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,name").Order("sort ASC")
		}).
		Find(&role).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *userRepository) CreateRole(req models.Role) error {
	err := r.db.Create(&req).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) UpdateRole(roleID int, permissionIDs []int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("role ID not found")
		}
		var rolePermissions []models.RolePermission
		for _, pid := range permissionIDs {
			rolePermissions = append(rolePermissions, models.RolePermission{
				RoleID:       roleID,
				PermissionID: pid,
			})
		}
		if len(rolePermissions) > 0 {
			if err := tx.Create(&rolePermissions).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *userRepository) DeleteRole(req models.Role) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("role_id = ?", req.ID).Delete(&models.RolePermission{})
		if res.Error != nil {
			return res.Error
		}

		res = tx.Where("id = ?", req.ID).Delete(&models.Role{})
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return errors.New("User ID not found")
		}

		return nil
	})
}

func (r *userRepository) DeleteUser(req models.User) error {
	tx := r.db.Where("id = ?", req.ID).Delete(&models.User{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("ID not found")
	}
	return nil
}

func (r *userRepository) GetUser() (user []models.User, err error) {
	err = r.db.Preload("Role", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,name")
	}).Order("created_at DESC").Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *userRepository) CheckUsername(req models.User) (user *models.User, err error) {
	err = r.db.Preload("Role.Permission", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, keyword")
	}).Where("username=?", req.Username).First(&user).Error
	if err != nil {
		return nil, err
	}
	err = r.db.Model(&user).Where("username=?", req.Username).Update("login_at", time.Now()).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetOneUser(req models.User) (user *models.User, err error) {
	err = r.db.
		Select("id,name,username,access_status,password,role_id,login_at,created_at,updated_at").
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,name")
		}).
		Preload("Role.Permission", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,name,keyword")
		}).
		Where("id=?", req.ID).
		First(&user).
		Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) CreateUser(model models.User) error {
	err := r.db.Create(&model).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *userRepository) GetRoleByName(roleName string) (models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", roleName).First(&role).Error
	if err != nil {
		return models.Role{}, err
	}
	return role, nil
}

func (r *userRepository) UpdateUser(req models.User) error {
	tx := r.db.Model(&models.User{}).Where("id = ?", req.ID).Updates(req)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("user ID not found")
	}
	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
