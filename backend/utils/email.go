package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/itsXanadium/singa-riset-backend/config"
)

type SMTPService struct {
	config config.Config
}

func NewSMTPService(config config.Config) *SMTPService {
	return &SMTPService{config: config}
}

func (s *SMTPService) SendVerificationCode(to string, code string) error {
	body := fmt.Sprintf("<h1>Your Verification Code is: %s </h1>", code)
	return s.send(to, "Email Verification Code", body)
}

func (s *SMTPService) SendNotification(to, subject, body string) error {
	return s.send(to, subject, buildNotificationHTML(subject, body))
}

func (s *SMTPService) send(to, subject, htmlBody string) error {
	if s.config.MailMailer == "log" {
		log.Printf("[mail:log] to=%s subject=%q\n%s", to, subject, htmlBody)
		return nil
	}

	from := mail.Address{Name: s.config.MailFromName, Address: s.config.MailFromAddress}
	headers := []string{
		"From: " + from.String(),
		"To: " + to,
		"Subject: " + mime.QEncoding.Encode("UTF-8", subject),
		"Date: " + time.Now().Format(time.RFC1123Z),
		"MIME-Version: 1.0",
		`Content-Type: text/html; charset="UTF-8"`,
	}
	message := []byte(strings.Join(headers, "\r\n") + "\r\n\r\n" + htmlBody)

	client, err := s.dial()
	if err != nil {
		return err
	}
	defer client.Close()

	if s.config.SMTPUser != "" {
		auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}
	if err := client.Mail(from.Address); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write(message); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}
	return client.Quit()
}

// dial connects with implicit TLS for MAIL_ENCRYPTION=tls/ssl; otherwise it
// upgrades with STARTTLS whenever the server offers it.
func (s *SMTPService) dial() (*smtp.Client, error) {
	addr := net.JoinHostPort(s.config.SMTPHost, strconv.Itoa(s.config.SMTPPort))
	tlsConfig := &tls.Config{ServerName: s.config.SMTPHost}

	switch s.config.MailEncryption {
	case "tls", "ssl", "smtps":
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 15 * time.Second}, "tcp", addr, tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("smtp tls dial: %w", err)
		}
		return smtp.NewClient(conn, s.config.SMTPHost)
	}

	conn, err := net.DialTimeout("tcp", addr, 15*time.Second)
	if err != nil {
		return nil, fmt.Errorf("smtp dial: %w", err)
	}
	client, err := smtp.NewClient(conn, s.config.SMTPHost)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("smtp handshake: %w", err)
	}
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(tlsConfig); err != nil {
			client.Close()
			return nil, fmt.Errorf("smtp starttls: %w", err)
		}
	} else if s.config.MailEncryption == "starttls" {
		client.Close()
		return nil, fmt.Errorf("smtp server %s does not support STARTTLS", s.config.SMTPHost)
	}
	return client, nil
}
func buildNotificationHTML(subject, body string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="margin:0;padding:0;background:#f4f5f7;font-family:Arial,Helvetica,sans-serif;">
    <table width="100%%" cellpadding="0" cellspacing="0" style="background:#f4f5f7;padding:24px 0;">
        <tr>
            <td align="center">
                <table width="600" cellpadding="0" cellspacing="0" style="background:#ffffff;border-radius:8px;overflow:hidden;box-shadow:0 1px 3px rgba(0,0,0,0.08);">
                    <!-- Header -->
                    <tr>
                        <td style="background:#1e40af;color:#ffffff;padding:20px 32px;">
                            <h1 style="margin:0;font-size:18px;font-weight:600;">Singa Riset</h1>
                        </td>
                    </tr>
                    <!-- Subject -->
                    <tr>
                        <td style="padding:32px 32px 8px 32px;">
                            <h2 style="margin:0;font-size:20px;color:#111827;">%s</h2>
                        </td>
                    </tr>
                    <!-- Body -->
                    <tr>
                        <td style="padding:8px 32px 32px 32px;color:#374151;font-size:15px;line-height:1.6;">
                            %s
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="background:#f9fafb;padding:20px 32px;border-top:1px solid #e5e7eb;">
                            <p style="margin:0;font-size:12px;color:#6b7280;">
                                This is an automated message from the Singa Riset Platform.<br/>
                                Please do not reply to this email.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, subject, body)
}
