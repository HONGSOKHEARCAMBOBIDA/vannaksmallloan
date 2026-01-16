package service

import (
	"fmt"
	"math"
	"time"

	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/helper"
	"github.com/hearbong/smallloanbackend/model"
	"github.com/hearbong/smallloanbackend/request"
	"github.com/hearbong/smallloanbackend/response"
	"github.com/hearbong/smallloanbackend/utils"
	"gorm.io/gorm"
)

type LoanService interface {
	Create(userID int, input request.LoanRequest) error
	GetLoanForCheck(userID int) ([]response.LoanResponse, error)
	GetLoanForApprove(userID int) ([]response.LoanResponse, error)
	CheckLoan(id int) error
	ApproveLoan(id int) error
	DeleteLoan(id int) error
}

type loanservice struct {
	db *gorm.DB
}

func NewLoanService() LoanService {
	return &loanservice{
		db: config.DB,
	}
}

func (s *loanservice) Create(userID int, input request.LoanRequest) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var loanproduct model.LoanProduct
	if err := tx.Where("id =?", input.LoanProductID).First(&loanproduct).Error; err != nil {
		tx.Rollback()
		return err
	}
	var cashiersession model.CashierSession
	if err := tx.Where("user_id =? AND status LIKE ?", userID, 1).First(&cashiersession).Error; err != nil {
		tx.Rollback()
		return err
	}
	interesAmount := math.Ceil(float64(input.LoanAmount)*float64(loanproduct.InterestRate))/100 + 0.00
	processfeeamount := math.Ceil(float64(input.LoanAmount)*float64(loanproduct.ProcessFeeRate))/100 + 0.00
	dailyprinciple := math.Ceil(float64(input.LoanAmount) / float64(loanproduct.TermDay))
	// Round up to nearest 100
	dailypaymentamount := math.Ceil((dailyprinciple+math.Ceil(interesAmount/float64(loanproduct.TermDay)))/100) * 100
	now := time.Now().Format("2006-01-02")
	loan := model.Loan{
		ClientID:           input.ClientID,
		CoID:               userID,
		LoanProductID:      input.LoanProductID,
		CashierSessionID:   cashiersession.ID,
		LoanAmount:         input.LoanAmount,
		InterestRate:       float32(interesAmount),
		ProcessFee:         float64(processfeeamount),
		DisbursedDate:      now,
		DisbursedBy:        userID,
		DailyPaymentAmount: float32(dailypaymentamount),
		Purpose:            input.Purpose,
		Duration:           loanproduct.TermDay,
		Status:             model.Pending,
		DocumentTypeID:     input.DocumentTypeID,
		CheckByID:          input.CheckByID,
		ApprovedByID:       input.ApprovedByID,
		ApproveDate:        nil,
		LoanStartDate:      nil,
		LoanEndDate:        nil,
		ClosedDate:         nil,
	}
	if err := tx.Create(&loan).Error; err != nil {
		tx.Rollback()
		return err
	}
	if len(input.Guarantors) > 0 {
		var guarantors []model.LoanGuarantor
		for _, g := range input.Guarantors {
			guarantor := model.LoanGuarantor{
				LoanID:       loan.ID,
				ClientID:     g.ClientID,
				Relationship: g.Relationship,
				SignedDate:   now,
				Notes:        g.Notes,
			}
			guarantors = append(guarantors, guarantor)
		}

		if err := tx.Create(&guarantors).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil

}

func (s *loanservice) GetLoanForCheck(userID int) ([]response.LoanResponse, error) {
	var loan []response.LoanResponse
	db := s.db.Table("loans l").
		Select(`
	l.id AS id,
	c.id AS client_id,
	c.name AS client_name,
	c.gender AS client_gender,
	c.marital_status AS client_marital_status,
	c.date_of_birth AS client_date_of_birth,
	c.occupation AS client_occupation,
	c.phone AS client_phone,

	p.id AS province_id,
	p.name AS province_name,
	d.id AS district_id,
	d.name AS district_name,
	cm.id AS communce_id,
	cm.name AS communce_name,
	v.id AS village_id,
	v.name AS village_name,

	c.latitude AS latitude,
	c.longitude AS longitude,

	u.id AS co_id,
	u.name AS co_name,

	lp.id AS loan_product_id,
	lp.name AS loan_product_name,

	l.loan_amount AS loan_amount,
	l.interest_rate AS interest_rate,
	l.process_fee AS process_fee,
	l.approve_date AS approve_date,
	l.loan_start_date AS loan_start_date,
	l.loan_end_date AS loan_end_date,
	l.disbursed_date AS disbursed_date,

	u.id AS disbursed_by,
	u.name AS disburse_by_name,

	l.daily_payment_amount AS daily_payment_amount,
	l.purpose AS purpose,
	l.duration AS duration,
	l.status AS status,

	dc.id AS document_type_id,
	dc.name AS document_type_name,

	uc.id AS check_by_id,
	uc.name AS check_by_name,

	up.id AS approve_by_id,
	up.name AS approve_by_name,

	l.closed_date AS close_date,
	l.closed_reason AS close_reason
`).
		Joins("LEFT JOIN clients c ON c.id = l.client_id").
		Joins("LEFT JOIN villages v ON v.id = c.village_id").
		Joins("LEFT JOIN communces cm ON cm.id = v.communce_id").
		Joins("LEFT JOIN districts d ON d.id = cm.district_id").
		Joins("LEFT JOIN provinces p ON p.id = d.province_id").
		Joins("LEFT JOIN users u ON u.id = l.co_id").
		Joins("LEFT JOIN loan_products lp ON lp.id = l.loan_product_id").
		Joins("LEFT JOIN document_types dc ON dc.id = l.document_type_id").
		Joins("LEFT JOIN users uc ON uc.id = l.check_by_id").
		Joins("LEFT JOIN users up ON up.id = l.approved_by_id")

	db = db.Where("check_by_id = ? AND status =?", userID, model.Pending)
	if err := db.Order("l.id DESC").Scan(&loan).Error; err != nil {
		return nil, err
	}
	for i := range loan {
		loan[i].ClientDoB = helper.FormatDate(loan[i].ClientDoB)
		loan[i].DisbursedDate = helper.FormatDate(loan[i].DisbursedDate)
		loan[i].ApproveDate = helper.FormatDate(loan[i].ApproveDate)
		loan[i].LoanStartDate = helper.FormatDate(loan[i].LoanStartDate)
		loan[i].LoanEndDate = helper.FormatDate(loan[i].LoanEndDate)
	}
	return loan, nil
}

func (s *loanservice) CheckLoan(id int) error {
	result := s.db.
		Model(&model.Loan{}).
		Where("id = ?", id).
		Update("status", model.Checked)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}

func (s *loanservice) GetLoanForApprove(userID int) ([]response.LoanResponse, error) {
	var loan []response.LoanResponse
	db := s.db.Table("loans l").
		Select(`
	l.id AS id,
	c.id AS client_id,
	c.name AS client_name,
	c.gender AS client_gender,
	c.marital_status AS client_marital_status,
	c.date_of_birth AS client_date_of_birth,
	c.occupation AS client_occupation,
	c.phone AS client_phone,

	p.id AS province_id,
	p.name AS province_name,
	d.id AS district_id,
	d.name AS district_name,
	cm.id AS communce_id,
	cm.name AS communce_name,
	v.id AS village_id,
	v.name AS village_name,

	c.latitude AS latitude,
	c.longitude AS longitude,

	u.id AS co_id,
	u.name AS co_name,

	lp.id AS loan_product_id,
	lp.name AS loan_product_name,

	l.loan_amount AS loan_amount,
	l.interest_rate AS interest_rate,
	l.process_fee AS process_fee,
	l.approve_date AS approve_date,
	l.loan_start_date AS loan_start_date,
	l.loan_end_date AS loan_end_date,
	l.disbursed_date AS disbursed_date,

	u.id AS disbursed_by,
	u.name AS disburse_by_name,

	l.daily_payment_amount AS daily_payment_amount,
	l.purpose AS purpose,
	l.duration AS duration,
	l.status AS status,

	dc.id AS document_type_id,
	dc.name AS document_type_name,

	uc.id AS check_by_id,
	uc.name AS check_by_name,

	up.id AS approve_by_id,
	up.name AS approve_by_name,

	l.closed_date AS close_date,
	l.closed_reason AS close_reason
`).
		Joins("LEFT JOIN clients c ON c.id = l.client_id").
		Joins("LEFT JOIN villages v ON v.id = c.village_id").
		Joins("LEFT JOIN communces cm ON cm.id = v.communce_id").
		Joins("LEFT JOIN districts d ON d.id = cm.district_id").
		Joins("LEFT JOIN provinces p ON p.id = d.province_id").
		Joins("LEFT JOIN users u ON u.id = l.co_id").
		Joins("LEFT JOIN loan_products lp ON lp.id = l.loan_product_id").
		Joins("LEFT JOIN document_types dc ON dc.id = l.document_type_id").
		Joins("LEFT JOIN users uc ON uc.id = l.check_by_id").
		Joins("LEFT JOIN users up ON up.id = l.approved_by_id")

	db = db.Where("approved_by_id = ? AND status =?", userID, model.Checked)
	if err := db.Order("l.id DESC").Scan(&loan).Error; err != nil {
		return nil, err
	}
	for i := range loan {
		loan[i].ClientDoB = helper.FormatDate(loan[i].ClientDoB)
		loan[i].DisbursedDate = helper.FormatDate(loan[i].DisbursedDate)
		loan[i].ApproveDate = helper.FormatDate(loan[i].ApproveDate)
		loan[i].LoanStartDate = helper.FormatDate(loan[i].LoanStartDate)
		loan[i].LoanEndDate = helper.FormatDate(loan[i].LoanEndDate)
	}
	return loan, nil
}

func (s *loanservice) ApproveLoan(id int) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	approvedate := time.Now().Format("2006-01-02")
	result := tx.Model(&model.Loan{}).Where("id = ?", id).Updates(&model.Loan{
		Status:      model.Approved,
		ApproveDate: &approvedate,
	})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("loan not found")
	}

	var loan model.Loan
	if err := tx.Where("id = ?", id).First(&loan).Error; err != nil {
		tx.Rollback()
		return err
	}

	var loanproduct model.LoanProduct
	if err := tx.Where("id = ?", loan.LoanProductID).First(&loanproduct).Error; err != nil {
		tx.Rollback()
		return err
	}

	receiptNumber := utils.GenerateReceiptNumber()
	newreceipt := model.Receipt{
		ReceiptNumber:    receiptNumber,
		LoanID:           loan.ID,
		ReceiptDate:      time.Now().Format("2006-01-02"),
		TotalAmount:      float64(loan.ProcessFee),
		CashierSessionID: loan.CashierSessionID,
		ReceiveBy:        loan.DisbursedBy,
		Notes:            "ទទួលបានពីសេវាកម្ចី",
	}

	if err := tx.Create(&newreceipt).Error; err != nil {
		tx.Rollback()
		return err
	}

	var cashiersession model.CashierSession
	if err := tx.Where("id = ?", loan.CashierSessionID).First(&cashiersession).Error; err != nil {
		tx.Rollback()
		return err
	}

	addTotalReceipts := cashiersession.TotalReceipts + loan.ProcessFee
	if err := tx.Model(&model.CashierSession{}).Where("id = ?", loan.CashierSessionID).
		Update("total_receipts", addTotalReceipts).Error; err != nil {
		tx.Rollback()
		return err
	}

	currentDate := time.Now().Format("2006-01-02")
	if err := tx.Model(&model.Loan{}).Where("id = ?", id).
		Update("loan_start_date", currentDate).Error; err != nil {
		tx.Rollback()
		return err
	}

	if loan.Duration > 0 {
		principal := math.Ceil(float64(loan.LoanAmount) / float64(loan.Duration))
		interest := math.Ceil(float64(loan.InterestRate) / float64(loan.Duration))
		dueamount := principal + interest

		// Calculate first payment date - TOMORROW (skip weekends if enabled)
		firstPaymentDate := utils.GetNextBusinessDay(currentDate, loanproduct.SkipWeeken)

		for i := 1; i <= loan.Duration; i++ {
			var scheduleDate string
			if i == 1 {
				scheduleDate = firstPaymentDate
			} else {
				daysToAdd := 0
				switch loanproduct.PaymentFrequency {
				case "បង់ប្រចាំថ្ងៃ":
					daysToAdd = 1
				case "WEEKLY":
					daysToAdd = 7
				case "MONTHLY":
					daysToAdd = 30
				default:
					daysToAdd = 1
				}

				scheduleDate = utils.CalculateNextScheduleDate(firstPaymentDate, i-1, daysToAdd, loanproduct.SkipWeeken)
			}

			newschedule := model.PaymentSchedule{
				LoanID:          loan.ID,
				ScheduleNumber:  i,
				PaymentDate:     scheduleDate,
				DueAmount:       dueamount,
				PrincipalAmount: principal,
				InterestAmount:  interest,
				PaidDate:        nil,
				PrincipalPaid:   nil,
				InterestPaid:    nil,
				PenaltyAmount:   loanproduct.LatePenaltyFixed,
				PenaltyPaid:     nil,
				PaidAmount:      nil,
				DayLate:         nil,
				Status:          model.ScheduelStatus(model.PENDING),
			}

			if err := tx.Create(&newschedule).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		// Calculate and update loan end date
		if loan.Duration > 0 {
			lastScheduleDate := utils.CalculateNextScheduleDate(firstPaymentDate, loan.Duration-1, 1, loanproduct.SkipWeeken)
			if err := tx.Model(&model.Loan{}).Where("id = ?", id).
				Update("loan_end_date", lastScheduleDate).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *loanservice) DeleteLoan(id int) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var loan model.Loan
	if err := s.db.Where("id =?", id).First(&loan).Error; err != nil {
		return err
	}

	if err := tx.Where("loan_id = ?", id).Delete(&model.Receipt{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	var cashiersession model.CashierSession
	if err := tx.Where("id =?", loan.CashierSessionID).First(&cashiersession).Error; err != nil {
		return err
	}

	removeTotalReceipts := cashiersession.TotalReceipts - loan.ProcessFee
	if err := tx.Model(&model.CashierSession{}).Where("id =?", loan.CashierSessionID).
		Update("total_receipts", removeTotalReceipts).Error; err != nil {
		tx.Rollback()
		return err

	}

	if err := tx.Where("loan_id =?", id).Delete(&model.PaymentSchedule{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id =?", id).Delete(&model.Loan{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil

}
