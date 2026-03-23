package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PhoneNumber  string         `gorm:"uniqueIndex" json:"phone_number"`
	Email        string         `gorm:"uniqueIndex" json:"email"`
	PasswordHash string         `json:"-"`
	IsVerified   bool           `json:"is_verified"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type VerificationCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	Code      string    `json:"code"`
	Type      string    `json:"type"` // sms, email
	Target    string    `json:"target"` // phone number or email
	ExpiresAt time.Time `json:"expires_at"`
	IsUsed    bool      `json:"is_used"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	UserID       uint      `json:"user_id"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Credential string `json:"credential" binding:"required"` // email or phone
	Password   string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	Code        string `json:"code" binding:"required"`
}

type VerifyCodeRequest struct {
	Target string `json:"target" binding:"required"` // phone or email
	Code   string `json:"code" binding:"required"`
	Type   string `json:"type" binding:"required"` // sms or email
}

type SendCodeRequest struct {
	Target string `json:"target" binding:"required"`
	Type   string `json:"type" binding:"required"`
}
