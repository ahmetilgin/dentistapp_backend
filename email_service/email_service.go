package email_service

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

type EmailService struct {
	smtpServer string
	smtpPort   string
	from       string
	password   string
}

func NewEmailService(smtpServer, smtpPort, from, password string) *EmailService {
	return &EmailService{
		smtpServer: smtpServer,
		smtpPort:   smtpPort,
		from:       from,
		password:   password,
	}
}

func (es *EmailService) SendEmail(to, subject, body string) error {
	// Set up authentication information.
	auth := smtp.PlainAuth("", es.from, es.password, es.smtpServer)

	// Connect to the server
	conn, err := net.Dial("tcp", es.smtpServer+":"+es.smtpPort)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	// Create a new SMTP client
	client, err := smtp.NewClient(conn, es.smtpServer)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	// Start TLS
	if err = client.StartTLS(&tls.Config{ServerName: es.smtpServer}); err != nil {
		return fmt.Errorf("failed to start TLS: %v", err)
	}

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	// Set the sender and recipient
	if err = client.Mail(es.from); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// Send the email body.
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to create data writer: %v", err)
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", es.from, to, subject, body)
	_, err = writer.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write email body: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %v", err)
	}

	return client.Quit()
}