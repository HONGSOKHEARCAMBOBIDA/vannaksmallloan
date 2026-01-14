package service

import (
	"time"

	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/helper"
	"github.com/hearbong/smallloanbackend/model"
	"github.com/hearbong/smallloanbackend/response"
	"github.com/hearbong/smallloanbackend/utils"
	"gorm.io/gorm"
)

type CashierSessionService interface {
	Create(userID int) error
	Get(userID int) ([]response.CashierSessionResponse, error)
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
		TotalReceipts:  0.00,
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

func (s *cashiersessionservice) Get(userID int) ([]response.CashierSessionResponse, error) {
	var cashiersession []response.CashierSessionResponse
	db := s.db.Table("cashier_sessions cs").Select(`
		cs.id AS id,
		cs.session_number AS session_number,
		u.id AS user_id,
		u.name AS user_name,
		cs.session_date AS session_date,
		cs.start_time AS start_time,
		cs.end_time AS end_time,
		cs.opening_balance AS opening_balance,
		cs.closing_balance AS closing_balance,
		cs.total_receipts AS total_receipts,
		cs.difference AS difference,
		cs.status AS status,
		cs.notes AS status,
		uv.id AS verified_by,
		uv.name AS verified_by_name,
		cs.verified_at AS verified_at
	`).
		Joins("LEFT JOIN users u ON u.id = cs.user_id").
		Joins("LEFT JOIN users uv ON uv.id = cs.verified_by")
	db = db.Where("user_id =? AND status =?", userID, "OPEN")
	if err := db.Scan(&cashiersession).Error; err != nil {
		return nil, err
	}
	for i := range cashiersession {
		cashiersession[i].SessionDate = helper.FormatDate(cashiersession[i].SessionDate)
		cashiersession[i].StartTime = helper.FormatTime(cashiersession[i].StartTime)
	}
	return cashiersession, nil
}
