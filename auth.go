package handlers

import (
	"auth-system/config"
	"auth-system/models"
	"auth-system/services"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db           *gorm.DB
	config       *config.Config
	smsService   *services.SMSService
	emailService *services.EmailService
}

func NewAuthHandler(cfg *config.Config, sms *services.SMSService, email *services.EmailService) *AuthHandler {
	return &AuthHandler{
		db:           database.GetDB(),
		config:       cfg,
		smsService:   sms,
		emailService: email,
	}
}

func (h *AuthHandler) SendSMS(c *gin.Context) {
	var req models.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Генерация 6-значного кода
	code := generateCode(6)

	// Сохранение кода в БД
	verificationCode := models.VerificationCode{
		Code:      code,
		Type:      "sms",
		Target:    req.Target,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := h.db.Create(&verificationCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save code"})
		return
	}

	// Отправка SMS
	if err := h.smsService.SendSMS(req.Target, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send SMS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMS sent successfully"})
}

func (h *AuthHandler) SendEmailCode(c *gin.Context) {
	var req models.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code := generateCode(6)

	verificationCode := models.VerificationCode{
		Code:      code,
		Type:      "email",
		Target:    req.Target,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := h.db.Create(&verificationCode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save code"})
		return
	}

	if err := h.emailService.SendVerificationCode(req.Target, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

func (h *AuthHandler) VerifySMS(c *gin.Context) {
	var req models.VerifyCodeRequest
	if err := c.ShouldBindJSON(&err); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var code models.VerificationCode
	err := h.db.Where("target = ? AND code = ? AND type = ? AND is_used = ? AND expires_at > ?",
		req.Target, req.Code, "sms", false, time.Now()).First(&code).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired code"})
		return
	}

	// Помечаем код как использованный
	code.IsUsed = true
	h.db.Save(&code)

	c.JSON(http.StatusOK, gin.H{"message": "SMS verified successfully"})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req models.VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var code models.VerificationCode
	err := h.db.Where("target = ? AND code = ? AND type = ? AND is_used = ? AND expires_at > ?",
		req.Target, req.Code, "email", false, time.Now()).First(&code).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired code"})
		return
	}

	code.IsUsed = true
	h.db.Save(&code)

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка кода подтверждения
	var code models.VerificationCode
	err := h.db.Where("target = ? AND code = ? AND type = ? AND is_used = ?",
		req.Email, req.Code, "email", false).First(&code).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Создание пользователя
	user := models.User{
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		PasswordHash: string(hashedPassword),
		IsVerified:   true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Помечаем код как использованный
	code.IsUsed = true
	h.db.Save(&code)

	// Генерация токенов
	tokens, err := h.generateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusCreated, tokens)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	
	// Поиск по email или телефону
	if err := h.db.Where("email = ? OR phone_number = ?", req.Credential, req.Credential).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Генерация токенов
	tokens, err := h.generateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"is_verified":  user.IsVerified,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" {
		// Здесь можно добавить токен в черный список через Redis
		// Для простоты просто возвращаем успех
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) generateTokens(user models.User) (map[string]string, error) {
	// Access token (15 минут)
	accessClaims := &models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Phone:  user.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Refresh token (7 дней)
	refreshClaims := &models.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Сохранение сессии в БД
	session := models.Session{
		ID:           generateSessionID(),
		UserID:       user.ID,
		Token:        accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}
	
	if err := h.db.Create(&session).Error; err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
	}, nil
}

func generateCode(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	for i := 0; i < length; i++ {
		bytes[i] = 48 + (bytes[i] % 10)
	}
	return string(bytes)
}

func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
