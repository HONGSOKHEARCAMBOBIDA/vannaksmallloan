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
}

type journalservice struct {
	db *gorm.DB
}

func NewJournalService() JournalService {
	return &journalservice{
		db: config.DB,
	}
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
		c.name AS chart_account_name,
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
	if err := db.Offset(offset).Limit(pagination.PageSize).Order("j.id DESC").Scan(&journal).Error; err != nil {
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
