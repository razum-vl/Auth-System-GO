package services

import (
	"auth-system/config"
	"fmt"
	"net/smtp"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{config: cfg}
}

func (e *EmailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth(
		"",
		e.config.EmailUser,
		e.config.EmailPassword,
		e.config.EmailHost,
	)
	
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)
	
	addr := fmt.Sprintf("%s:%s", e.config.EmailHost, e.config.EmailPort)
	err := smtp.SendMail(addr, auth, e.config.EmailUser, []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	
	return nil
}

func (e *EmailService) SendVerificationCode(to, code string) error {
	subject := "Verification Code"
	body := fmt.Sprintf("Your verification code is: %s\n\nThis code will expire in 10 minutes.", code)
	
	return e.SendEmail(to, subject, body)
}
