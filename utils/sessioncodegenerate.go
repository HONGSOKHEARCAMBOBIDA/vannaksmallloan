package utils

import (
	"fmt"
	"time"

	"github.com/hearbong/smallloanbackend/model"
	"gorm.io/gorm"
)

func GenerateSessionNumber(tx *gorm.DB) (string, error) {
	today := time.Now().Format("20060102")

	var count int64
	err := tx.Model(&model.CashierSession{}).
		Where("session_number LIKE ?", "CS-"+today+"%").
		Count(&count).Error

	if err != nil {
		return "", err
	}

	sequence := count + 1
	sessionNumber := fmt.Sprintf("CS-%s-%04d", today, sequence)

	return sessionNumber, nil
}
