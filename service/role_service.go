package service

import (
	"errors"

	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/model"
	"github.com/hearbong/smallloanbackend/request"
	"gorm.io/gorm"
)

type RoleService interface {
	GetRole() ([]model.Role, error)
	CreateRole(input request.RoleRequestCreate) error
	UpdateRole(id int, input request.RoleRequestUpdate) error
	ChangeStatusRole(id int) error
}

type roleservice struct {
	db *gorm.DB
}

func NewRoleService() RoleService {
	return &roleservice{
		db: config.DB,
	}
}

func (s *roleservice) GetRole() ([]model.Role, error) {
	var role []model.Role
	if err := s.db.Find(&role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (s *roleservice) CreateRole(input request.RoleRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	newrole := model.Role{
		Name:        input.Name,
		DisPlayName: input.DisPlayName,
		IsActive:    true,
	}
	if err := tx.Create(&newrole).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *roleservice) UpdateRole(id int, input request.RoleRequestUpdate) error {
	updates := map[string]interface{}{}

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.DisPlayName != nil {
		updates["display_name"] = *input.DisPlayName
	}
	if len(updates) == 0 {
		return errors.New("no fields to update")
	}
	result := s.db.Model(&model.Role{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("role not found")
	}

	return nil
}

func (s *roleservice) ChangeStatusRole(id int) error {
	result := s.db.Model(&model.Role{}).Where("id =?", id).Update("is_active", gorm.Expr("!is_active"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("role not found or not change")
	}
	return nil
}
