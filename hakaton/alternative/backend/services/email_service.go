package services

import (
	"backend/config"
	"fmt"
	"net/smtp"
)

type EmailService struct {
	config *config.SMTPConfig
}

func NewEmailService(cfg *config.SMTPConfig) *EmailService {
	return &EmailService{config: cfg}
}

func (e *EmailService) SendPasswordResetEmail(to, token string) error {
	subject := "Запрос на восстановление пароля"
	body := fmt.Sprintf(`
Здравствуйте!

Вы запросили восстановление пароля. Пожалуйста, перейдите по ссылке ниже для сброса пароля:

http://localhost:8080/reset-password?token=%s

Эта ссылка действительна в течение 1 часа.

Если вы не запрашивали восстановление пароля, проигнорируйте это письмо.

С уважением,
Команда Alternative
`, token)

	return e.sendEmail(to, subject, body)
}

func (e *EmailService) sendEmail(to, subject, body string) error {
	// Create the message with proper headers
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n", e.config.From, to, subject)
	message := []byte(headers + body)

	// Set up authentication
	auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.Host)

	// Send email using smtp.SendMail (handles STARTTLS automatically)
	addr := fmt.Sprintf("%s:%d", e.config.Host, e.config.Port)
	err := smtp.SendMail(addr, auth, e.config.From, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
