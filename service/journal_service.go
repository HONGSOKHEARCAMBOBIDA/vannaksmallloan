package service

import (
	"errors"
	"math"
	"strings"

	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/helper"
	"github.com/hearbong/smallloanbackend/model"
	"github.com/hearbong/smallloanbackend/request"
	"github.com/hearbong/smallloanbackend/response"
	"github.com/hearbong/smallloanbackend/utils"
	"gorm.io/gorm"
)

type JournalService interface {
	Create(userID int, input request.JournalRequestCreate) error
	Get(filters map[string]string, pagination request.Pagination) ([]response.JournalResponse, *model.PaginationMetadata, error)
	Update(id int, input request.JournalRequestUpdate) error
	Delete(id int) error
	GenerateBalanceSheet(asofDate string) (*response.BalanceSheetResponse, error)
	Incomestatment(asofDate string)(response.Incomestatment,error)
}

type journalservice struct {
	db *gorm.DB
}

func NewJournalService() JournalService {
	return &journalservice{
		db: config.DB,
	}
}

func (s *journalservice) processAssets(accounts []AccountBalanceResult) response.BalanceSheetSection {
	section := response.BalanceSheetSection{
		Title: "ASSETS",
	}
	for _, acc := range accounts {
		if acc.AccountType == "ASSET" {
			section.Accounts = append(section.Accounts, response.AccountBalance{
				AccountCode: acc.AccountCode,
				AccountName: acc.AccountName,
				Balance:     acc.Balance,
			})
			section.Total += acc.Balance
		}
	}
	return section
}
func (s *journalservice) processLiabilities(accounts []AccountBalanceResult) response.BalanceSheetSection {
	section := response.BalanceSheetSection{
		Title: "LIABILITIES",
	}

	for _, acc := range accounts {
		if acc.AccountType == "LIABILITY" {
			section.Accounts = append(section.Accounts, response.AccountBalance{
				AccountCode: acc.AccountCode,
				AccountName: acc.AccountName,
				Balance:     acc.Balance,
			})
			section.Total += acc.Balance
		}
	}

	return section
}

func (s *journalservice) processEquity(accounts []AccountBalanceResult) response.BalanceSheetSection {
	section := response.BalanceSheetSection{
		Title: "EQUITY",
	}

	for _, acc := range accounts {
		if acc.AccountType == "EQUITY" || acc.AccountType == "INCOME" {
			section.Accounts = append(section.Accounts, response.AccountBalance{
				AccountCode: acc.AccountCode,
				AccountName: acc.AccountName,
				Balance:     acc.Balance,
			})
			section.Total += acc.Balance
		}
	}

	return section
}

func (s *journalservice) calculateTotals(res *response.BalanceSheetResponse) response.BalanceSheetTotals {
	totals := response.BalanceSheetTotals{
		TotalAssets:            res.Assets.Total,
		TotalLiabilities:       res.Liabilities.Total,
		TotalEquity:            res.Equity.Total,
		TotalLiabilitiesEquity: res.Liabilities.Total + res.Equity.Total,
	}

	totals.Difference = totals.TotalAssets - totals.TotalLiabilitiesEquity

	return totals
}

func (s *journalservice) checkIfBalanced(totals response.BalanceSheetTotals) bool {
	// Allow small rounding differences
	diff := totals.TotalAssets - totals.TotalLiabilitiesEquity
	return diff >= -0.01 && diff <= 0.01
}

type AccountBalanceResult struct {
	AccountCode string  `gorm:"column:account_code"`
	AccountName string  `gorm:"column:account_name"`
	AccountType string  `gorm:"column:account_type_name"`
	Balance     float64 `gorm:"column:balance"`
}

func (s *journalservice) GenerateBalanceSheet(asOfDate string) (*response.BalanceSheetResponse, error) {
	response := &response.BalanceSheetResponse{
		ReportTitle: "BALANCE SHEET",
		ReportDate:  asOfDate,
	}

	// Get account balances
	var accountBalances []AccountBalanceResult

	// Query to get all account balances
	query := `
        SELECT 
            ca.code AS account_code,
            ca.description AS account_name,
            at.name AS account_type_name,
            CASE 
                WHEN at.name IN ('ASSET', 'EXPENSE') 
                THEN COALESCE(SUM(j.debit_amount - j.credit_amount), 0)
                ELSE COALESCE(SUM(j.credit_amount - j.debit_amount), 0)
            END AS balance
        FROM chart_accounts ca
        LEFT JOIN account_types at ON ca.account_type_id = at.id
        LEFT JOIN journals j ON ca.id = j.chart_account_id 
            AND j.transaction_date <= ?
        GROUP BY ca.id, ca.code, ca.description, at.name
        HAVING balance != 0
        ORDER BY ca.code
    `

	if err := s.db.Raw(query, asOfDate).Scan(&accountBalances).Error; err != nil {
		return nil, err
	}

	// Process accounts into sections
	response.Assets = s.processAssets(accountBalances)
	response.Liabilities = s.processLiabilities(accountBalances)
	response.Equity = s.processEquity(accountBalances)

	// Calculate totals
	response.Totals = s.calculateTotals(response)

	// Check if balance sheet is balanced
	response.Totals.IsBalanced = s.checkIfBalanced(response.Totals)
	response.Message = "Balance sheet generate successfully"

	return response, nil
}

func (s *journalservice) Create(userID int, input request.JournalRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	journalCode := utils.GenerateJournalCode()

	debitJournal := model.Journal{
		TransactionDate: input.TransactionDate,
		ChartAccountID:  input.DebitAccountID,
		DebitAmount:     input.Amount,
		CreditAmount:    0,
		Description:     input.Description,
		ReferenceCode:   journalCode,
		CreatedBy:       userID,
	}
	creditJournal := model.Journal{
		TransactionDate: input.TransactionDate,
		ChartAccountID:  input.CreditAccountID,
		DebitAmount:     0,
		CreditAmount:    input.Amount,
		Description:     input.Description,
		ReferenceCode:   journalCode,
		CreatedBy:       userID,
	}
	if err := tx.Create(&[]model.Journal{debitJournal, creditJournal}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *journalservice) Incomestatment(asofDate string) (response.Incomestatment, error) {
	var res response.Incomestatment

	query := `
		SELECT
			COALESCE(SUM(
				CASE 
					WHEN at.name = 'INCOME' 
					THEN j.credit_amount - j.debit_amount 
					ELSE 0 
				END
			), 0) AS total_income,

			COALESCE(SUM(
				CASE 
					WHEN at.name = 'EXPENSE' 
					THEN j.debit_amount - j.credit_amount 
					ELSE 0 
				END
			), 0) AS total_expense
		FROM chart_accounts ca
		LEFT JOIN account_types at ON at.id = ca.account_type_id
		LEFT JOIN journals j 
			ON ca.id = j.chart_account_id
			AND j.transaction_date <= ?
	`

	if err := s.db.Raw(query, asofDate).Scan(&res).Error; err != nil {
		return response.Incomestatment{}, err
	}

	return res, nil
}


func (s *journalservice) Update(id int, input request.JournalRequestUpdate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updates := map[string]interface{}{}

	if input.TransactionDate != nil {
		updates["transaction_date"] = *input.TransactionDate
	}
	if input.ChartAccountID != nil {
		updates["chart_account_id"] = *input.ChartAccountID
	}
	if input.DebitAmount != nil {
		updates["debit_amount"] = *input.DebitAmount
	}
	if input.CreditAmount != nil {
		updates["credit_amount"] = *input.CreditAmount
	}
	if input.Description != nil {
		updates["description"] = *input.Description
	}

	if len(updates) == 0 {
		return errors.New("no field to update")
	}

	result := tx.Model(&model.Journal{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("journal not found")
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (s *journalservice) Get(filters map[string]string, pagination request.Pagination) ([]response.JournalResponse, *model.PaginationMetadata, error) {
	var journal []response.JournalResponse
	var totalCount int64
	offset := (pagination.Page - 1) * pagination.PageSize
	db := s.db.Table("journals j").Select(`
		j.id AS id,
		j.transaction_date AS transaction_date,
		c.id AS chart_account_id,
		c.code AS chart_account_code,
		c.description AS chart_account_name,
		j.debit_amount AS debit_amount,
		j.credit_amount AS credit_amount,
		j.description AS description,
		j.reference_id AS reference_id,
		j.reference_code AS reference_code,
		u.id AS created_by,
		u.name AS created_by_name
	`).
		Joins("LEFT JOIN chart_accounts c ON c.id = j.chart_account_id").
		Joins("LEFT JOIN users u ON u.id = j.created_by")
	if v, ok := filters["reference_code"]; ok && v != "" {
		db = db.Where("j.reference_code LIKE ?", "%"+v+"%")
	}
	if v, ok := filters["between"]; ok && v != "" {
		dates := strings.Split(v, ",")
		if len(dates) == 2 {
			db = db.Where("j.transaction_date BETWEEN ? AND ?", dates[0], dates[1])
		}
	}
	//http://localhost:8080/api/journals?between=2024-01-01,2024-01-31&page=1&pageSize=10
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}
	if err := db.Offset(offset).Limit(pagination.PageSize).Scan(&journal).Error; err != nil {
		return nil, nil, err
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PageSize)))
	for i := range journal {
		journal[i].TransactionDate = helper.FormatDate(journal[i].TransactionDate)
	}
	return journal, &model.PaginationMetadata{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalCount: int(totalCount),
		TotalPages: totalPages,
		HasNext:    pagination.Page < totalPages,
		HasPrev:    pagination.Page > 1,
	}, nil
}

func (s *journalservice) Delete(id int) error {
	result := s.db.Delete(&model.Journal{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("journal not found")
	}

	return nil
}
