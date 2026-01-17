package service

import (
	"math"
	"time"

	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/model"
	"github.com/hearbong/smallloanbackend/request"
	"github.com/hearbong/smallloanbackend/response"
	"github.com/hearbong/smallloanbackend/utils"
	"gorm.io/gorm"
)

type ReceiptService interface {
	Collectfromgoodloan(filters map[string]string, pagination request.Pagination) ([]response.CollectfromgoodloanResponse, *model.PaginationMetadata, error)
	CreateReceipt(id int, userID int, input request.ReceiptRequest) error
}

type receiptservice struct {
	db *gorm.DB
}

func NewReceiptService() ReceiptService {
	return &receiptservice{
		db: config.DB,
	}
}

func (s *receiptservice) Collectfromgoodloan(filters map[string]string, pagination request.Pagination) ([]response.CollectfromgoodloanResponse, *model.PaginationMetadata, error) {

	var result []response.CollectfromgoodloanResponse
	var totalCount int64

	offset := (pagination.Page - 1) * pagination.PageSize

	baseQuery := s.db.Table("loans l").
		Joins("LEFT JOIN clients c ON c.id = l.client_id").
		Joins("LEFT JOIN villages v ON v.id = c.village_id").
		Joins("LEFT JOIN users u ON u.id = l.co_id").
		Where("l.status = ?", 3)

	if v := filters["client_name"]; v != "" {
		baseQuery = baseQuery.Where("c.name LIKE ?", "%"+v+"%")
	}

	if v := filters["village_name"]; v != "" {
		baseQuery = baseQuery.Where("v.name LIKE ?", "%"+v+"%")
	}

	if err := baseQuery.Session(&gorm.Session{}).Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}

	err := baseQuery.Select(`
		l.id AS id,
		c.id AS client_id,
		c.name AS client_name,
		u.id AS user_id,
		u.name AS user_name,
		v.id AS village_id,
		v.name AS village_name,
		(
			SELECT COALESCE(SUM(
				CASE 
					WHEN COALESCE(due_amount,0) != COALESCE(paid_amount, 0) AND DATE(payment_date) <= CURRENT_DATE
					THEN (COALESCE(due_amount,0) - COALESCE(paid_amount, 0))
					ELSE 0 
				END
			), 0)
			FROM payment_schedules ps
			WHERE ps.loan_id = l.id
		) AS total_collect,
		(
			SELECT COALESCE(SUM(
				CASE 
					WHEN COALESCE(due_amount,0) != COALESCE(paid_amount, 0)
					AND DATE(payment_date) < CURRENT_DATE 
					AND COALESCE(penalty_amount,0) != COALESCE(penalty_paid,0)
					THEN (COALESCE(penalty_amount,0) - COALESCE(penalty_paid, 0))
					ELSE 0
				END
			), 0)
			FROM payment_schedules ps
			WHERE ps.loan_id = l.id
		) AS total_penalty
	`).
		Order("l.id DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Scan(&result).Error

	if err != nil {
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PageSize)))

	return result, &model.PaginationMetadata{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalCount: int(totalCount),
		TotalPages: totalPages,
		HasNext:    pagination.Page < totalPages,
		HasPrev:    pagination.Page > 1,
	}, nil
}

func (s *receiptservice) CreateReceipt(id int, userID int, input request.ReceiptRequest) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var loan model.Loan
	if err := tx.First(&loan, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	var cs model.CashierSession
	if err := tx.Where("user_id = ? AND status = ?", userID, 1).
		Order("id DESC").
		First(&cs).Error; err != nil {
		tx.Rollback()
		return err
	}

	newReceipt := model.Receipt{
		ReceiptNumber:    utils.GenerateReceiptNumber(),
		LoanID:           id,
		ReceiptDate:      time.Now().Format("2006-01-02"),
		TotalAmount:      input.TotalReceipt,
		CashierSessionID: cs.ID,
		ReceiveBy:        loan.DisbursedBy,
	}

	if err := tx.Create(&newReceipt).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.CashierSession{}).
		Where("id = ?", cs.ID).
		Update("total_receipts", cs.TotalReceipts+input.TotalReceipt).Error; err != nil {
		tx.Rollback()
		return err
	}

	remaining := input.TotalReceipt

	for remaining > 0 {
		var ps model.PaymentSchedule
		err := tx.
			Where("loan_id = ? AND status = ?", id, model.PENDING).
			Order("schedule_number ASC").
			First(&ps).Error

		if err != nil {
			tx.Rollback()

			break
		}

		pp, ip, penp := 0.0, 0.0, 0.0

		if ps.PrincipalPaid != nil {
			pp = *ps.PrincipalPaid
		}
		if ps.InterestPaid != nil {
			ip = *ps.InterestPaid
		}
		if ps.PenaltyPaid != nil {
			penp = float64(*ps.PenaltyPaid)
		}

		penaltyDue := 0.0
		today := time.Now().Format("2006-01-02")

		if ps.PaymentDate < today && float64(ps.PenaltyAmount) > penp {
			penaltyDue = float64(ps.PenaltyAmount) - penp
		}

		principalDue := ps.PrincipalAmount - pp
		interestDue := ps.InterestAmount - ip

		payPenalty := math.Min(remaining, penaltyDue)
		remaining -= payPenalty

		payPrincipal := math.Min(remaining, principalDue)
		remaining -= payPrincipal

		payInterest := math.Min(remaining, interestDue)
		remaining -= payInterest

		newPP := pp + payPrincipal
		newIP := ip + payInterest
		newPen := penp + payPenalty

		ps.PrincipalPaid = &newPP
		ps.InterestPaid = &newIP
		ps.PenaltyPaid = func() *float32 {
			v := float32(newPen)
			return &v
		}()

		paidAmount := newPP + newIP + newPen
		ps.PaidAmount = &paidAmount

		if paidAmount >= ps.DueAmount {
			now := time.Now().Format("2006-01-02")
			ps.PaidDate = &now
			ps.Status = model.PAID
		}

		updates := map[string]interface{}{
			"principal_paid": ps.PrincipalPaid,
			"interest_paid":  ps.InterestPaid,
			"penalty_paid":   ps.PenaltyPaid,
			"paid_amount":    ps.PaidAmount,
			"status":         ps.Status,
		}
		if ps.PaidDate != nil {
			updates["paid_date"] = ps.PaidDate
		}
		if err := tx.Model(&model.PaymentSchedule{}).Where("id = ?", ps.ID).Updates(updates).Error; err != nil {
			tx.Rollback()
			return err
		}

		allocation := model.ReceiptAllocation{
			ReceiptID:       newReceipt.ID,
			ScheduleID:      ps.ID,
			PrincipalAmount: payPrincipal,
			InterestAmount:  payInterest,
			PenaltyAmount:   payPenalty,
			AllocationDate:  time.Now().Format("2006-01-02"),
		}

		if err := tx.Create(&allocation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
