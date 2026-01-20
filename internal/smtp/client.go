// Package smtp provides SMTP client functionality.
package smtp

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// Client wraps SMTP configuration.
type Client struct {
	host     string
	port     int
	tls      bool
	startTLS bool
	insecure bool
	noTLS    bool
	email    string
	password string
}

// Config holds SMTP connection configuration.
type Config struct {
	Host     string
	Port     int
	TLS      bool
	StartTLS bool
	Insecure bool // Skip TLS cert verification
	NoTLS    bool // Disable TLS entirely
	Email    string
	Password string
}

// NewClient creates a new SMTP client.
func NewClient(cfg Config) *Client {
	return &Client{
		host:     cfg.Host,
		port:     cfg.Port,
		tls:      cfg.TLS,
		startTLS: cfg.StartTLS,
		insecure: cfg.Insecure,
		noTLS:    cfg.NoTLS,
		email:    cfg.Email,
		password: cfg.Password,
	}
}

// Message represents an email to send.
type Message struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
}

// Send sends an email message.
func (c *Client) Send(msg *Message) error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	// Build recipients list
	recipients := make([]string, 0, len(msg.To)+len(msg.Cc)+len(msg.Bcc))
	recipients = append(recipients, msg.To...)
	recipients = append(recipients, msg.Cc...)
	recipients = append(recipients, msg.Bcc...)

	// Build email content
	var content strings.Builder
	content.WriteString(fmt.Sprintf("From: %s\r\n", msg.From))
	content.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ", ")))
	if len(msg.Cc) > 0 {
		content.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(msg.Cc, ", ")))
	}
	content.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))
	content.WriteString("MIME-Version: 1.0\r\n")
	content.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	content.WriteString("\r\n")
	content.WriteString(msg.Body)

	tlsConfig := &tls.Config{
		ServerName:         c.host,
		InsecureSkipVerify: c.insecure,
	}

	var client *smtp.Client
	var err error

	if c.noTLS {
		client, err = smtp.Dial(addr)
	} else if c.tls {
		client, err = smtp.DialTLS(addr, tlsConfig)
	} else if c.startTLS {
		client, err = smtp.DialStartTLS(addr, tlsConfig)
	} else {
		client, err = smtp.Dial(addr)
	}
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	// Authenticate
	auth := sasl.NewPlainClient("", c.email, c.password)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender
	if err := client.Mail(msg.From, nil); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, rcpt := range recipients {
		if err := client.Rcpt(rcpt, nil); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", rcpt, err)
		}
	}

	// Send data
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data: %w", err)
	}

	if _, err := wc.Write([]byte(content.String())); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close data: %w", err)
	}

	return client.Quit()
}

// TestConnection tests the SMTP connection.
func (c *Client) TestConnection() error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	var client *smtp.Client
	var err error

	tlsConfig := &tls.Config{
		ServerName:         c.host,
		InsecureSkipVerify: c.insecure,
	}

	if c.noTLS {
		// Plain text connection
		client, err = smtp.Dial(addr)
	} else if c.tls {
		// Direct TLS (SMTPS, port 465)
		client, err = smtp.DialTLS(addr, tlsConfig)
	} else if c.startTLS {
		// STARTTLS (port 587)
		client, err = smtp.DialStartTLS(addr, tlsConfig)
	} else {
		client, err = smtp.Dial(addr)
	}
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	// Authenticate
	auth := sasl.NewPlainClient("", c.email, c.password)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	return client.Quit()
}
