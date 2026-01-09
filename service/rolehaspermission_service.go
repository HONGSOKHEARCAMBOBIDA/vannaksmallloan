package service

import (
	"errors"

	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/model"
	"github.com/hearbong/smallloanbackend/request"
	"github.com/hearbong/smallloanbackend/response"
	"gorm.io/gorm"
)

type RolePermissionService interface {
	CreateRolePermissions(input request.CreateRolePermissionInput) error
	DeleteRolePermission(input request.DeleteRolePermissionsInput) error
	GetRolePermission(roleID int) (map[string][]response.PermissionWithAssignedRole, error)
}

type rolepermissionservice struct {
	db *gorm.DB
}

func NewRolePermissionService() RolePermissionService {
	return &rolepermissionservice{
		db: config.DB,
	}
}

func (s *rolepermissionservice) CreateRolePermissions(input request.CreateRolePermissionInput) error {
	if len(input.PermissionIDs) == 0 {
		return errors.New("permission_ids cannot be empty")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var rolePermissions []model.RolePermission
	for _, permissionID := range input.PermissionIDs {
		rolePermissions = append(rolePermissions, model.RolePermission{
			RoleID:       uint(input.RoleID),
			PermissionID: uint(permissionID),
		})
	}

	if err := tx.Create(&rolePermissions).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *rolepermissionservice) DeleteRolePermission(input request.DeleteRolePermissionsInput) error {
	if len(input.PermissionIDs) == 0 {
		return errors.New("permission_ids cannot be empty")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.
		Where("role_id = ? AND permission_id IN ?", input.RoleID, input.PermissionIDs).
		Delete(&model.RolePermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *rolepermissionservice) GetRolePermission(roleID int) (map[string][]response.PermissionWithAssignedRole, error) {
	var permissions []response.PermissionWithAssignedRole

	err := s.db.
		Table("permissions").
		Select(`
			permissions.id,
			permissions.name,
			permissions.display_name,
			permissions.permission_group AS permission_group,
			permissions.sort_order AS sort_order,
			CASE 
				WHEN role_has_permissions.permission_id IS NOT NULL THEN true
				ELSE false
			END AS assigned
		`).
		Joins(`
			LEFT JOIN role_has_permissions 
			ON permissions.id = role_has_permissions.permission_id 
			AND role_has_permissions.role_id = ?
		`, roleID).
		Order("permissions.permission_group, permissions.sort_order ASC").
		Scan(&permissions).Error

	if err != nil {
		return nil, err
	}
	groupedPermissions := make(map[string][]response.PermissionWithAssignedRole)
	for _, perm := range permissions {
		groupedPermissions[perm.GroupName] = append(groupedPermissions[perm.GroupName], perm)
	}

	return groupedPermissions, nil
}
