package service

import (
	"errors"
	"strings"

	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/model"
	"github.com/hearbong/smallloanbackend/request"
	"github.com/hearbong/smallloanbackend/response"
	"github.com/hearbong/smallloanbackend/utils"
	"gorm.io/gorm"
)

type UserService interface {
	Register(input request.UserRequestCreate) error
	GetUser(filters map[string]string) ([]response.UserResponse, error)
	Update(id int, input request.UserRequestUpdate) error
	ChangeStatusUser(id int) error
	ResetPassword(id int, input request.UserRequestResetPassword) error
}

type userservice struct {
	db *gorm.DB
}

func NewUserService() UserService {
	return &userservice{
		db: config.DB,
	}
}

func (s *userservice) Register(input request.UserRequestCreate) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	hashpassword := utils.HashPassword(input.Password)
	username := strings.ToLower(
		strings.Join(strings.Fields(input.Name), ""))
	email := username + "168@vlaon.com"
	newuser := model.User{
		Name:     input.Name,
		Username: username,
		Password: hashpassword,
		RoleId:   input.RoleId,
		Email:    email,
		Phone:    input.Phone,
		Isactive: true,
	}
	if err := tx.Create(&newuser).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *userservice) GetUser(filters map[string]string) ([]response.UserResponse, error) {
	var user []response.UserResponse
	db := s.db.Table("users").Select(`
		users.id AS id,
		users.name AS name,
		users.username AS username,
		r.id AS role_id,
		r.name AS role_name,
		users.email AS email,
		users.phone AS phone,
		users.isactive AS is_active
	`).
		Joins("INNER JOIN roles r ON r.id = users.role_id")
	if v, ok := filters["name"]; ok && v != "" {
		db = db.Where("users.name LIKE ? OR users.username LIKE ?", "%"+v+"%", "%"+v+"%")
	}
	if v, ok := filters["role_id"]; ok && v != "" {
		db = db.Where("r.id =?", v)
	}
	if v, ok := filters["is_active"]; ok && v != "" {
		db = db.Where("users.isactive =?", v)
	}
	if err := db.Order("users.id DESC").Scan(&user).Error; err != nil {
		return nil, err
	}
	return user, nil

}

func (s *userservice) Update(id int, input request.UserRequestUpdate) error {
	updates := map[string]interface{}{}
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.UserName != nil {
		updates["username"] = *input.UserName
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.Phone != nil {
		updates["phone"] = *input.Phone
	}
	if input.RoleId != nil {
		updates["role_id"] = *input.RoleId
	}
	if len(updates) == 0 {
		return errors.New("no fields to update")
	}
	result := s.db.Model(&model.User{}).Where("id =?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *userservice) ChangeStatusUser(id int) error {
	result := s.db.Model(&model.User{}).Where("id =?", id).Update("isactive", gorm.Expr("!isactive"))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *userservice) ResetPassword(id int, input request.UserRequestResetPassword) error {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return errors.New("user not found")
	}
	hashedPassword := utils.HashPassword(input.NewPassword)
	if err := s.db.Model(&model.User{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error; err != nil {
		return err
	}
	return nil
}
