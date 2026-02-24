package service

import (
	"fmt"
	"net/smtp"

	"github.com/0DayMonxrch/project-management-system/internal/config"
)

type emailService struct {
	cfg config.SMTPConfig
}

func NewEmailService(cfg config.SMTPConfig) *emailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) SendVerificationEmail(to, token string) error {
	subject := "Verify your email"
	body := fmt.Sprintf("Click the link to verify your email: http://localhost:8080/api/v1/auth/verify-email/%s", token)
	return s.send(to, subject, body)
}

func (s *emailService) SendPasswordResetEmail(to, token string) error {
	subject := "Reset your password"
	body := fmt.Sprintf("Click the link to reset your password: http://localhost:8080/api/v1/auth/reset-password/%s", token)
	return s.send(to, subject, body)
}

func (s *emailService) send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", s.cfg.From, to, subject, body)
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, []byte(msg))
}