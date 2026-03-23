package services

import (
	"auth-system/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SMSService struct {
	config *config.Config
}

func NewSMSService(cfg *config.Config) *SMSService {
	return &SMSService{config: cfg}
}

// Пример для Twilio (можно заменить на любого провайдера)
func (s *SMSService) SendSMS(phoneNumber, code string) error {
	// Пример для Twilio API
	url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.config.SMSAPIKey)
	
	payload := map[string]string{
		"To":   phoneNumber,
		"From": "+1234567890", // Ваш номер Twilio
		"Body": fmt.Sprintf("Your verification code is: %s", code),
	}
	
	jsonPayload, _ := json.Marshal(payload)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	
	req.SetBasicAuth(s.config.SMSAPIKey, s.config.SMSAPISecret)
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send SMS: %s", resp.Status)
	}
	
	return nil
}

// Альтернативный вариант для локальной разработки
func (s *SMSService) SendSMSMock(phoneNumber, code string) error {
	fmt.Printf("[MOCK SMS] To: %s, Code: %s\n", phoneNumber, code)
	return nil
}
