package service

import (
	"fmt"
	"time"

	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/helper"
	"github.com/hearbong/smallloanbackend/model"
	"github.com/hearbong/smallloanbackend/response"
	"gorm.io/gorm"
)

type PaymentScheduleService interface {
	RemovePenalty(id int) error
	GetFullPaymentSchedule(loanID int) (response.FullPaymentResponse, error)
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

func (s *paymentschedulesservice) GetFullPaymentSchedule(loanID int) (response.FullPaymentResponse, error) {
	var result response.FullPaymentResponse

	baseQuery := s.db.Table("loans l").
		Joins("LEFT JOIN clients c ON c.id = l.client_id").
		Joins("LEFT JOIN users co ON co.id = l.co_id").
		Where("l.id = ?", loanID).
		Where("l.status = ?", 3)

	// Since we're getting by loanID, we expect only one loan
	var loan struct {
		LoanID       int     `gorm:"column:loan_id"`
		ClientID     int     `gorm:"column:client_id"`
		ClientName   string  `gorm:"column:client_name"`
		ClientGender string  `gorm:"column:client_gender"`
		ClientPhone  string  `gorm:"column:client_phone"`
		LoanAmount   float64 `gorm:"column:loan_amount"`
		ProcessFee   float64 `gorm:"column:process_fee"`
		ApproveDate  *string `gorm:"column:approve_date"`
		Purpose      string  `gorm:"column:purpose"`
		Duration     int     `gorm:"column:duration"`
		CoID         int     `gorm:"column:co_id"`
		CoName       string  `gorm:"column:co_name"`
		CoPhone      string  `gorm:"column:co_phone"`
	}

	err := baseQuery.Select(`
        l.id AS loan_id,
        c.id AS client_id,
        c.name AS client_name,
        CASE 
            WHEN c.gender = 1 THEN 'ប្រុស'
            WHEN c.gender = 2 THEN 'ស្រី'
            ELSE 'មិនបានកំណត់'
        END AS client_gender,
        c.phone AS client_phone,
        l.loan_amount AS loan_amount,
        l.process_fee AS process_fee,
        l.approve_date AS approve_date,
        l.purpose AS purpose,
        l.duration AS duration,
        l.co_id AS co_id,
        co.name AS co_name,
        co.phone AS co_phone
    `).Scan(&loan).Error


	if err != nil {
		return result, err
	}


	// Check if loan exists
	if loan.LoanID == 0 {
		return result, fmt.Errorf("loan not found")
	}

	result = response.FullPaymentResponse{
		LoanID:       loan.LoanID,
		ClientID:     loan.ClientID,
		ClientName:   loan.ClientName,
		ClientGender: loan.ClientGender, // This should be string, not int
		ClientPhone:  loan.ClientPhone,
		LoanAmount:   loan.LoanAmount,
		ProcessFee:   loan.ProcessFee,
		ApproveDate:  loan.ApproveDate,
		Purpose:      loan.Purpose,
		Duration:     loan.Duration,
		CoID:         loan.CoID,
		CoName:       loan.CoName,
		CoPhone:      loan.CoPhone,
		Schedule:     []response.Schedule{},
	}

	
if result.ApproveDate != nil {
	formatted := helper.FormatDate(*result.ApproveDate)
	result.ApproveDate = &formatted
}



	// Get payment schedules
	var schedules []model.PaymentSchedule
	if err := s.db.Where("loan_id = ?", loan.LoanID).
		Order("schedule_number ASC").
		Find(&schedules).Error; err != nil {
		return result, err
	}
	for i := range schedules {
		schedules[i].PaymentDate = helper.FormatDate(schedules[i].PaymentDate)
		
	}

	for _, ps := range schedules {
		dueAmount := ps.DueAmount

		penaltyAmount := 0.0

		// Calculate penalty - only if payment is overdue AND not fully paid
		currentTime := time.Now().Format("2006-01-02")
		if ps.PaymentDate < currentTime {
			// Check if not fully paid
			var paidAmount float64
			if ps.PaidAmount != nil {
				paidAmount = *ps.PaidAmount
			}

			if paidAmount < dueAmount {
		
				penaltyAmount = ps.PenaltyAmount 
			}
		}

		var paidAmount float64
		if ps.PaidAmount != nil {
			paidAmount = *ps.PaidAmount
		}

		var penaltyPaid float64
		if ps.PenaltyPaid != nil {
			penaltyPaid = *ps.PenaltyPaid
		}

		totalDue := dueAmount + penaltyAmount 
		remainingPenalty := penaltyAmount - penaltyPaid

		totalOwe := totalDue - paidAmount 
		// Determine status
		status := "មិនទាន់បង់"
		if ps.PaidDate != nil {
			if paidAmount >= totalDue {
				status = "បានបង់"
			} else if paidAmount > 0 {
				status = "PARTIAL"
			}
		} else if ps.PaymentDate < time.Now().Format("2006-01-02") {
			status = "យឺត"
		}

		schedule := response.Schedule{
			ID:             ps.ID,
			ScheduleNumber: ps.ScheduleNumber,
			PaymentDate:    ps.PaymentDate,
			DueAmount:      dueAmount,
			Penalty:        remainingPenalty,
			Total:          totalDue,
			PaidAmount:     paidAmount,
			Totalowe:       totalOwe,
			Status:         status,
			PenaltyPaid: penaltyPaid,
		}

		result.Schedule = append(result.Schedule, schedule)

	}

	return result, nil
}
