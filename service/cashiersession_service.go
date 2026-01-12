package service

import (
	"time"

	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/model"
	"github.com/hearbong/smallloanbackend/utils"
	"gorm.io/gorm"
)

type CashierSessionService interface {
	Create(userID int) error
}

type cashiersessionservice struct {
	db *gorm.DB
}

func NewCashierSessionService() CashierSessionService {
	return &cashiersessionservice{
		db: config.DB,
	}
}

func (s *cashiersessionservice) Create(userID int) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	sessionNumber, err := utils.GenerateSessionNumber(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	cashiersession := model.CashierSession{
		SessionNumber:  sessionNumber,
		UserID:         userID,
		SessionDate:    time.Now().Format("2006-01-02"),
		StartTime:      time.Now(),
		EndTime:        nil,
		OpeningBalance: 0.00,
		ClosingBalance: nil,
		TotalReceipts:  nil,
		Difference:     nil,
		Notes:          nil,
		VerifiedBy:     nil,
		VerifiedAt:     nil,
		Status:         model.OPEN,
	}

	if err := tx.Create(&cashiersession).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
