package utils

import (
	"time"

	"github.com/hearbong/smallloanbackend/model"
	"gorm.io/gorm"
)

func CreateJournalPair(tx *gorm.DB, userID, referenceID int, journalCode string, chartDebit int, chartCredit int, amount float64, description string) error {
	journals := []model.Journal{
		{
			TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
			ChartAccountID:  chartDebit,
			DebitAmount:     amount,
			CreditAmount:    0,
			Description:     description,
			ReferenceID:     referenceID,
			ReferenceCode:   journalCode,
			CreatedBy:       userID,
		},
		{
			TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
			ChartAccountID:  chartCredit,
			DebitAmount:     0,
			CreditAmount:    amount,
			Description:     description,
			ReferenceID:     referenceID,
			ReferenceCode:   journalCode,
			CreatedBy:       userID,
		},
	}

	return tx.Create(&journals).Error
}
