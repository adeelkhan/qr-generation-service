package qr

import (
	"errors"

	"github.com/adeelkhan/qr-service/internal/models"
	"github.com/google/uuid"
	goqr "github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

var (
	ErrNotFound  = errors.New("qr code not found")
	ErrForbidden = errors.New("access denied")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Generate(userID uuid.UUID, inputText, title string) (*models.QRCode, error) {
	png, err := goqr.Encode(inputText, goqr.Medium, 256)
	if err != nil {
		return nil, err
	}
	code := &models.QRCode{
		UserID:    userID,
		InputText: inputText,
		ImageData: png,
		Title:     title,
	}
	if err := s.db.Create(code).Error; err != nil {
		return nil, err
	}
	return code, nil
}

func (s *Service) List(userID uuid.UUID) ([]models.QRCode, error) {
	var codes []models.QRCode
	if err := s.db.Where("user_id = ?", userID).
		Select("id, user_id, input_text, title, created_at").
		Order("created_at DESC").
		Find(&codes).Error; err != nil {
		return nil, err
	}
	return codes, nil
}

func (s *Service) Get(userID, codeID uuid.UUID) (*models.QRCode, error) {
	var code models.QRCode
	if err := s.db.First(&code, "id = ?", codeID).Error; err != nil {
		return nil, ErrNotFound
	}
	if code.UserID != userID {
		return nil, ErrForbidden
	}
	return &code, nil
}

func (s *Service) Delete(userID, codeID uuid.UUID) error {
	var code models.QRCode
	if err := s.db.First(&code, "id = ?", codeID).Error; err != nil {
		return ErrNotFound
	}
	if code.UserID != userID {
		return ErrForbidden
	}
	return s.db.Delete(&code).Error
}
