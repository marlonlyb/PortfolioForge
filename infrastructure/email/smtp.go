package email

import (
	"fmt"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/marlonlyb/portfolioforge/model"
)

type SMTPConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	FromAddress string
	FromName    string
}

type SMTPMailer struct {
	config SMTPConfig
}

type NoopMailer struct{}

func NewSMTPMailer(config SMTPConfig) *SMTPMailer {
	return &SMTPMailer{config: config}
}

func NewNoopMailer() *NoopMailer {
	return &NoopMailer{}
}

func (m *SMTPMailer) SendEmailVerificationOTP(message model.EmailVerificationMessage) error {
	fromHeader := m.config.FromAddress
	if strings.TrimSpace(m.config.FromName) != "" {
		fromHeader = fmt.Sprintf("%s <%s>", m.config.FromName, m.config.FromAddress)
	}

	subject := "PortfolioForge email verification"
	body := fmt.Sprintf(
		"Hello,\r\n\r\nYour PortfolioForge verification code is %s. It expires in %d minutes.\r\n\r\nIf you did not request this code, you can ignore this email.\r\n",
		message.OTPCode,
		message.ExpiresInMinute,
	)

	rawMessage := strings.Join([]string{
		fmt.Sprintf("From: %s", fromHeader),
		fmt.Sprintf("To: %s", message.ToEmail),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	address := fmt.Sprintf("%s:%d", m.config.Host, m.config.Port)
	var auth smtp.Auth
	if strings.TrimSpace(m.config.Username) != "" || strings.TrimSpace(m.config.Password) != "" {
		auth = smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
	}

	return smtp.SendMail(address, auth, m.config.FromAddress, []string{message.ToEmail}, []byte(rawMessage))
}

func (m *NoopMailer) SendEmailVerificationOTP(model.EmailVerificationMessage) error {
	return nil
}

func ParsePort(raw string) (int, error) {
	port, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0, err
	}
	if port <= 0 {
		return 0, fmt.Errorf("smtp port must be positive")
	}
	return port, nil
}
