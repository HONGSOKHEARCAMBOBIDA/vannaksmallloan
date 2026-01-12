package service

import (
	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/model"
	"gorm.io/gorm"
)

type DocumentTypeService interface {
	GetAll() ([]model.DocumentType, error)
}

type documenttypeservice struct {
	db *gorm.DB
}

func NewDocumentTypeService() DocumentTypeService {
	return &documenttypeservice{
		db: config.DB,
	}
}

func (s *documenttypeservice) GetAll() ([]model.DocumentType, error) {
	var documenttype []model.DocumentType
	if err := s.db.Find(&documenttype).Error; err != nil {
		return nil, err
	}
	return documenttype, nil
}
