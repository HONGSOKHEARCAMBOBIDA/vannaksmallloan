package service

import (
	"errors"

	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/model"
	"github.com/hearbong/smallloanbackend/request"
	"github.com/hearbong/smallloanbackend/response"
	"gorm.io/gorm"
)

type ChartAccountService interface {
	Create(input request.ChartAccountRequestCreate) error
	Get() ([]response.ChartAccountResponse, error)
	Update(id int, input request.ChartAccountRequestUpdate) error
	ChangeStatusChartAccount(id int) error
}

type chartaccountservice struct {
	db *gorm.DB
}

func NewChartAccountService() ChartAccountService {
	return &chartaccountservice{
		db: config.DB,
	}
}

func (s *chartaccountservice) Create(input request.ChartAccountRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	newchartaccount := model.ChartAccount{
		Code:          input.Code,
		Name:          input.Name,
		Description:   input.Description,
		AccountTypeID: input.AccountTypeID,
		IsActive:      true,
	}
	if err := tx.Create(&newchartaccount).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *chartaccountservice) Get() ([]response.ChartAccountResponse, error) {
	var chartaccount []response.ChartAccountResponse
	db := s.db.Table("chart_accounts").Select(`
		chart_accounts.id AS id,
		chart_accounts.code AS code,
		chart_accounts.name AS name,
		chart_accounts.description AS description,
		chart_accounts.is_active AS is_active,
		a.id AS account_type_id,
		a.name AS account_type_name
	`).
		Joins("LEFT JOIN account_types a ON a.id = chart_accounts.account_type_id")
	if err := db.Scan(&chartaccount).Error; err != nil {
		return nil, err
	}
	return chartaccount, nil
}

func (s *chartaccountservice) Update(id int, input request.ChartAccountRequestUpdate) error {
	updates := map[string]interface{}{}
	if input.Code != nil {
		updates["code"] = *input.Code
	}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.AccountTypeID != nil {
		updates["account_type_id"] = *input.AccountTypeID
	}
	if len(updates) == 0 {
		return errors.New("No Field to update")
	}
	result := s.db.Model(&model.ChartAccount{}).Where("id =?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("chartaccocunt not update")
	}
	return nil
}

func (s *chartaccountservice) ChangeStatusChartAccount(id int) error {
	result := s.db.Model(&model.ChartAccount{}).Where("id =?", id).Update("is_active", gorm.Expr("!is_active"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("chartaccount not found")
	}
	return nil
}
