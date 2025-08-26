package user_service

import (
	"laotop_final/models"
	"laotop_final/utils"
)

func MapUserToResponseInfos(req *models.User) *ResponseInfos {
	var rolePermissions []string
	for _, p := range req.Role.Permission {
		rolePermissions = append(rolePermissions, p.Keyword)
	}
	return &ResponseInfos{
		ID:           req.ID,
		Name:         req.Name,
		Username:     req.Username,
		AccessStatus: utils.PointerBoolToBool(req.AccessStatus),
		Role: Role{
			Name:       req.Role.Name,
			Permission: rolePermissions,
		},
		Profile:   "",
		LoginAt:   req.LoginAt.Format("2006-01-02 15:04:05"),
		CreatedAt: req.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdateAt:  req.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
