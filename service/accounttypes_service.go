package service

import (
	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/model"
	"gorm.io/gorm"
)

type AccountTypeService interface {
	Get() ([]model.AccountType, error)
}

type accounttypeservice struct {
	db *gorm.DB
}

func NewAccountTypeService() AccountTypeService {
	return &accounttypeservice{
		db: config.DB,
	}
}

func (s *accounttypeservice) Get() ([]model.AccountType, error) {
	var accounttype []model.AccountType
	if err := s.db.Find(&accounttype).Error; err != nil {
		return nil, err
	}
	return accounttype, nil
}
