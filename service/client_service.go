package service

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hearbong/smallloanbackend/config"
	"github.com/hearbong/smallloanbackend/helper" // assuming you have a helper package
	"github.com/hearbong/smallloanbackend/model"
	"github.com/hearbong/smallloanbackend/request"
	"github.com/hearbong/smallloanbackend/response"
	"gorm.io/gorm"
)

type ClientService interface {
	GetAll() ([]model.Client, error)
	GetList(filters map[string]string, pagination request.Pagination) ([]response.ClientResponse, *model.PaginationMetadata, error)
	Create(input request.ClientRequestCreate, c *gin.Context, userID int) error
	Update(id int, input request.ClientRequestUpdate, c *gin.Context, userID int) error
	ChangeStatusClient(id int) error
}

type clientservice struct {
	db *gorm.DB
}

func NewClientService() ClientService {
	return &clientservice{
		db: config.DB,
	}
}

func (s *clientservice) Create(input request.ClientRequestCreate, c *gin.Context, userID int) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var imagePath string
	file, err := c.FormFile("clientimage")
	if err == nil && file != nil {
		if !helper.ProtectImage(file) {
			tx.Rollback()
			return fmt.Errorf("invalid image file")
		}
		clientimageDir := "public/clientimage"
		if _, err := os.Stat(clientimageDir); os.IsNotExist(err) {
			if err := os.MkdirAll(clientimageDir, os.ModePerm); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create directory: %v", err)
			}
		}
		extension := filepath.Ext(file.Filename)
		clientimageName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
		clientimagePath := filepath.Join(clientimageDir, clientimageName)
		if err := c.SaveUploadedFile(file, clientimagePath); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to save image: %v", err)
		}
		imagePath = clientimagePath
	}
	gender := model.Gender(input.Gender)
	client := model.Client{
		Name:          input.Name,
		Gender:        gender,
		MaritatStatus: input.MaritatStatus,
		DateOfBirth:   input.DateofBirth,
		Occupation:    input.Occupation,
		IdCardNumber:  input.IdCardNumber,
		Phone:         input.Phone,
		VillageID:     input.VillageID,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
		ImagePath:     imagePath,
		Note:          input.Note,
		IsActive:      true,
		CreateBy:      userID,
	}
	if err := tx.Create(&client).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (s *clientservice) Update(id int, input request.ClientRequestUpdate, c *gin.Context, userID int) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find existing client
	var client model.Client
	if err := tx.First(&client, id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("client not found")
		}
		return err
	}

	// Handle file upload if exists
	file, err := c.FormFile("clientimage")
	if err == nil && file != nil {
		// Validate image
		if !helper.ProtectImage(file) {
			tx.Rollback()
			return fmt.Errorf("invalid image file")
		}
		clientimageDir := "public/clientimage"
		if _, err := os.Stat(clientimageDir); os.IsNotExist(err) {
			if err := os.MkdirAll(clientimageDir, os.ModePerm); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create directory: %v", err)
			}
		}
		extension := filepath.Ext(file.Filename)
		clientimageName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
		clientimagePath := filepath.Join(clientimageDir, clientimageName)
		if err := c.SaveUploadedFile(file, clientimagePath); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to save image: %v", err)
		}
		if client.ImagePath != "" {
			if err := os.Remove(client.ImagePath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Warning: Failed to delete old image: %v\n", err)
			}
		}
		client.ImagePath = clientimagePath
	}

	// Update client fields
	gender := model.Gender(input.Gender)
	client.Name = input.Name
	client.Gender = gender
	client.MaritatStatus = input.MaritatStatus
	client.DateOfBirth = input.DateofBirth
	client.Occupation = input.Occupation
	client.IdCardNumber = input.IdCardNumber
	client.Phone = input.Phone
	client.VillageID = input.VillageID
	client.Latitude = input.Latitude
	client.Longitude = input.Longitude
	client.Note = input.Note
	client.CreateBy = userID
	if err := tx.Save(&client).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (s *clientservice) GetAll() ([]model.Client, error) {

	var clients []model.Client

	if err := s.db.Joins("LEFT JOIN loans ON clients.id = loans.client_id").Where("loans.client_id IS NULL").Find(&clients).Error; err != nil {
		return nil, err
	}

	return clients, nil
}

func (s *clientservice) GetList(filters map[string]string, pagination request.Pagination) ([]response.ClientResponse, *model.PaginationMetadata, error) {
	var client []response.ClientResponse
	var totalCount int64
	offset := (pagination.Page - 1) * pagination.PageSize
	db := s.db.Table("clients").Select(`
		clients.id AS id,
		clients.name AS name,
		clients.gender AS gender,
		clients.marital_status AS marital_status,
		clients.date_of_birth AS date_of_birth,
		clients.occupation AS occupation,
		clients.id_card_number AS id_card_number,
		clients.phone AS phone,
		clients.latitude AS latitude,
		clients.longitude AS longitude,
		clients.image_path AS image_path,
		clients.notes AS notes,
		clients.is_active AS is_active,
		u.id AS created_by,
		u.name AS create_by_name,
		p.id AS province_id,
		p.name AS province_name,
		d.id AS district_id,
		d.name AS district_name,
		c.id AS communce_id,
		c.name AS communce_name,
		v.id AS village_id,
		v.name AS village_name
	`).
		Joins("LEFT JOIN users u ON u.id = clients.created_by").
		Joins("LEFT JOIN villages v ON v.id = clients.village_id").
		Joins("LEFT JOIN communces c ON c.id = v.communce_id").
		Joins("LEFT JOIN districts d ON d.id = c.district_id").
		Joins("LEFT JOIN provinces p ON p.id = d.province_id")

	if v, ok := filters["name"]; ok && v != "" {
		db = db.Where("clients.name LIKE ?", "%"+v+"%")
	}
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, nil, err
	}
	if err := db.Offset(offset).Limit(pagination.PageSize).Order("clients.id DESC").Scan(&client).Error; err != nil {
		return nil, nil, err
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PageSize)))
	for i := range client {
		client[i].DateofBirth = helper.FormatDate(client[i].DateofBirth)
	}
	return client, &model.PaginationMetadata{
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalCount: int(totalCount),
		TotalPages: totalPages,
		HasNext:    pagination.Page < totalPages,
		HasPrev:    pagination.Page > 1,
	}, nil
}

func (s *clientservice) ChangeStatusClient(id int) error {
	var client model.Client
	if err := s.db.First(&client, id).Error; err != nil {
		return err
	}

	client.IsActive = !client.IsActive
	return s.db.Save(&client).Error
}
