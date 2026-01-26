package service

import (
	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/model"
	"gorm.io/gorm"
)

type CommunceService interface {
	GetCommunce(id int) ([]model.Communce, error)
}

type communceservice struct {
	db *gorm.DB
}

func NewCommunceService() CommunceService {
	return &communceservice{
		db: config.DB,
	}
}

func (s communceservice) GetCommunce(id int) ([]model.Communce, error) {
	var communce []model.Communce
	if err := s.db.Where("district_id =?", id).Find(&communce).Error; err != nil {
		return nil, err
	}
	return communce, nil
}
