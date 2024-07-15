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

func (es *EmailService) SendEmail(to, subject, htmlBody, textBody string) error {
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

	boundary := "nhK2RUJRHs5MZ9TH"
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: multipart/alternative; boundary=%s\r\n"+
		"\r\n"+
		"--%s\r\n"+
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n"+
		"\r\n"+
		"--%s\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n"+
		"\r\n"+
		"--%s--",
		es.from, to, subject, boundary, boundary, textBody, boundary, htmlBody, boundary)

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

func (es *EmailService) SendPasswordResetEmail(to, host, resetToken string) error {
	subject := "Password Reset"
	htmlBody := fmt.Sprintf(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Password Reset</title>
			</head>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
					<h1 style="color: #444;">Password Reset</h1>
					<p>You have requested to reset your password. Click the link below to reset your password:</p>
					<p>
						<a href="%s/reset-password/%s" style="display: inline-block; padding: 10px 20px; background-color: #007bff; color: #ffffff; text-decoration: none; border-radius: 5px;">Reset Password</a>
					</p>
					<p>If you did not request a password reset, please ignore this email.</p>
					<p>This link will expire in 1 hour for security reasons.</p>
				</div>
			</body>
			</html>
		`, host, resetToken)

	textBody := fmt.Sprintf(`
	Password Reset

	You have requested to reset your password. Use the link below to reset your password:

	%s/reset-password/%s

	If you did not request a password reset, please ignore this email.
	This link will expire in 1 hour for security reasons.
		`, host, resetToken)

	return es.SendEmail(to, subject, htmlBody, textBody)
}
