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
	Verify(userID int, id int) error
	RollbackVerify(id int) error
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

func (s *cashiersessionservice) Verify(userID int, id int) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var cs model.CashierSession
	if err := tx.Where("id = ?", id).First(&cs).Error; err != nil {
		tx.Rollback()
		return err
	}

	journalCode := utils.GenerateJournalCode()
	var totalLoan float64
	tx.Model(&model.Loan{}).Where("cashier_sessions_id = ?", id).Select("COALESCE(SUM(loan_amount), 0)").Scan(&totalLoan)
	if err := utils.CreateJournalPair(tx, userID, id, journalCode, 2, 1, totalLoan, "លុយទម្លាក់កម្ចី"); err != nil {
		tx.Rollback()
		return err
	}
	var totalLoanFee float64
	tx.Model(&model.Receipt{}).
		Where("cashier_session_id = ? AND notes LIKE ?", id, "ទទួលបានពីសេវាកម្ចី").
		Select("COALESCE(SUM(total_amount), 0)").Scan(&totalLoanFee)
	if err := utils.CreateJournalPair(tx, userID, id, journalCode, 1, 8, totalLoanFee, "ទទួលបានពីសេវាកម្ចី"); err != nil {
		tx.Rollback()
		return err
	}
	var receipts []model.Receipt
	tx.Where("cashier_session_id = ?", id).Find(&receipts)
	receiptIDs := make([]int, len(receipts))
	for i, r := range receipts {
		receiptIDs[i] = r.ID
	}
	var totalPrincipal, totalInterest, totalPenalty float64
	if len(receiptIDs) > 0 {
		tx.Model(&model.ReceiptAllocation{}).Where("receipt_id IN ?", receiptIDs).Select("COALESCE(SUM(principal_amount), 0)").Scan(&totalPrincipal)
		tx.Model(&model.ReceiptAllocation{}).Where("receipt_id IN ?", receiptIDs).Select("COALESCE(SUM(interest_amount), 0)").Scan(&totalInterest)
		tx.Model(&model.ReceiptAllocation{}).Where("receipt_id IN ?", receiptIDs).Select("COALESCE(SUM(penalty_amount), 0)").Scan(&totalPenalty)
	}

	if err := utils.CreateJournalPair(tx, userID, id, journalCode, 1, 2, totalPrincipal, "ប្រាក់ដេីមអតិថិជនបង់ត្រឡប់"); err != nil {
		tx.Rollback()
		return err
	}
	if err := utils.CreateJournalPair(tx, userID, id, journalCode, 1, 7, totalInterest, "ចំណូលពីការប្រាក់"); err != nil {
		tx.Rollback()
		return err
	}
	if err := utils.CreateJournalPair(tx, userID, id, journalCode, 1, 9, totalPenalty, "ចំណូលពីការពិន័យ"); err != nil {
		tx.Rollback()
		return err
	}
	closingBalance := cs.TotalReceipts
	difference := closingBalance - cs.OpeningBalance
	if err := tx.Model(&model.CashierSession{}).Where("id = ?", id).Updates(map[string]interface{}{
		"closing_balance": closingBalance,
		"end_time":        time.Now(),
		"difference":      difference,
		"status":          model.CLOSED,
		"verified_by":     userID,
		"notes":           "បានផ្ទៀងផ្ទាត់",
		"verified_at":     time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *cashiersessionservice) RollbackVerify(id int) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var cs model.CashierSession
	if err := tx.Where("id = ?", id).First(&cs).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("reference_id =?", id).Delete(&model.Journal{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	closingBalance := 0.00
	difference := 0.00
	if err := tx.Model(&model.CashierSession{}).Where("id = ?", id).Updates(map[string]interface{}{
		"closing_balance": closingBalance,
		"end_time":        nil,
		"difference":      difference,
		"status":          model.OPEN,
		"verified_by":     nil,
		"notes":           nil,
		"verified_at":     nil,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error

}
