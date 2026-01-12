package service

import (
	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/model"
	"gorm.io/gorm"
)

type LoanProductService interface {
	GetAll() ([]model.LoanProduct, error)
}

type loanproduct struct {
	db *gorm.DB
}

func NewLoanProductService() LoanProductService {
	return &loanproduct{
		db: config.DB,
	}
}

func (s *loanproduct) GetAll() ([]model.LoanProduct, error) {
	var loanproduct []model.LoanProduct
	if err := s.db.Find(&loanproduct).Error; err != nil {
		return nil, err
	}
	return loanproduct, nil
}
