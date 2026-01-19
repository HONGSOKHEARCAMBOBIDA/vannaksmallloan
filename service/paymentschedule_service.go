package service

import (
	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/model"
	"gorm.io/gorm"
)

type PaymentScheduleService interface {
	RemovePenalty(id int) error
}

type paymentschedulesservice struct {
	db *gorm.DB
}

func NewPaymentScheduleService() PaymentScheduleService {
	return &paymentschedulesservice{
		db: config.DB,
	}
}

func (s *paymentschedulesservice) RemovePenalty(id int) error {
	result := s.db.Model(&model.PaymentSchedule{}).Where("loan_id =?", id).Update("penalty_amount", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return result.Error
	}
	return nil
}
