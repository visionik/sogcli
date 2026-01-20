package cli

import (
	"fmt"

	"github.com/visionik/sog/internal/config"
	"github.com/visionik/sog/internal/discover"
	"github.com/visionik/sog/internal/imap"
	"github.com/visionik/sog/internal/smtp"
)

// AuthCmd handles account management.
type AuthCmd struct {
	Add    AuthAddCmd    `cmd:"" help:"Add an IMAP/SMTP account"`
	List   AuthListCmd   `cmd:"" help:"List configured accounts"`
	Test   AuthTestCmd   `cmd:"" help:"Test account connection"`
	Remove AuthRemoveCmd `cmd:"" help:"Remove an account"`
}

// AuthAddCmd adds a new account.
type AuthAddCmd struct {
	Email    string `arg:"" help:"Email address for the account"`
	IMAPHost string `help:"IMAP server hostname" name:"imap-host"`
	IMAPPort int    `help:"IMAP server port" name:"imap-port" default:"993"`
	SMTPHost string `help:"SMTP server hostname" name:"smtp-host"`
	SMTPPort int    `help:"SMTP server port" name:"smtp-port" default:"587"`
	Password string `help:"Password (will prompt if not provided)"`
	Discover bool   `help:"Auto-discover servers from DNS"`
	Insecure bool   `help:"Skip TLS certificate verification"`
	NoTLS    bool   `help:"Disable TLS (plain text connection)" name:"no-tls"`
}

// Run executes the auth add command.
func (c *AuthAddCmd) Run(root *Root) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Auto-discover if --discover flag set
	if c.Discover {
		fmt.Printf("Auto-discovering servers for %s...\n", c.Email)
		result, err := discover.Discover(c.Email)
		if err != nil {
			return fmt.Errorf("auto-discover failed: %w", err)
		}
		if result.IMAP != nil && c.IMAPHost == "" {
			c.IMAPHost = result.IMAP.Host
			c.IMAPPort = result.IMAP.Port
			fmt.Printf("  IMAP: %s:%d\n", c.IMAPHost, c.IMAPPort)
		}
		if result.SMTP != nil && c.SMTPHost == "" {
			c.SMTPHost = result.SMTP.Host
			c.SMTPPort = result.SMTP.Port
			fmt.Printf("  SMTP: %s:%d\n", c.SMTPHost, c.SMTPPort)
		}
	}

	// Validate required fields
	if c.IMAPHost == "" {
		return fmt.Errorf("--imap-host is required (or use --discover)")
	}
	if c.SMTPHost == "" {
		return fmt.Errorf("--smtp-host is required (or use --discover)")
	}

	// TODO: Prompt for password if not provided
	if c.Password == "" {
		return fmt.Errorf("--password is required (keyring integration coming soon)")
	}

	acct := config.Account{
		Email: c.Email,
		IMAP: config.ServerConfig{
			Host:     c.IMAPHost,
			Port:     c.IMAPPort,
			TLS:      !c.NoTLS,
			Insecure: c.Insecure,
			NoTLS:    c.NoTLS,
		},
		SMTP: config.ServerConfig{
			Host:     c.SMTPHost,
			Port:     c.SMTPPort,
			TLS:      !c.NoTLS,
			StartTLS: !c.NoTLS,
			Insecure: c.Insecure,
			NoTLS:    c.NoTLS,
		},
	}

	if err := cfg.AddAccount(acct, c.Password); err != nil {
		return fmt.Errorf("failed to add account: %w", err)
	}

	fmt.Printf("Added account: %s\n", c.Email)
	return nil
}

// AuthListCmd lists configured accounts.
type AuthListCmd struct{}

// Run executes the auth list command.
func (c *AuthListCmd) Run(root *Root) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	accounts := cfg.ListAccounts()
	if len(accounts) == 0 {
		fmt.Println("No accounts configured. Use 'sog auth add' to add one.")
		return nil
	}

	for _, acct := range accounts {
		marker := "  "
		if acct.Email == cfg.DefaultAccount {
			marker = "* "
		}
		fmt.Printf("%s%s (IMAP: %s:%d, SMTP: %s:%d)\n",
			marker, acct.Email,
			acct.IMAP.Host, acct.IMAP.Port,
			acct.SMTP.Host, acct.SMTP.Port)
	}

	return nil
}

// AuthTestCmd tests account connection.
type AuthTestCmd struct {
	Email string `arg:"" optional:"" help:"Account to test (default: default account)"`
}

// Run executes the auth test command.
func (c *AuthTestCmd) Run(root *Root) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	email := c.Email
	if email == "" {
		email = root.Account
	}
	if email == "" {
		email = cfg.DefaultAccount
	}
	if email == "" {
		return fmt.Errorf("no account specified and no default set")
	}

	acct, err := cfg.GetAccount(email)
	if err != nil {
		return err
	}

	password, err := cfg.GetPassword(email)
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	fmt.Printf("Testing %s...\n", email)

	// Test IMAP
	fmt.Printf("  IMAP %s:%d... ", acct.IMAP.Host, acct.IMAP.Port)
	imapClient, err := imap.Connect(imap.Config{
		Host:     acct.IMAP.Host,
		Port:     acct.IMAP.Port,
		TLS:      acct.IMAP.TLS,
		Insecure: acct.IMAP.Insecure,
		NoTLS:    acct.IMAP.NoTLS,
		Email:    email,
		Password: password,
	})
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("OK")
		imapClient.Close()
	}

	// Test SMTP
	fmt.Printf("  SMTP %s:%d... ", acct.SMTP.Host, acct.SMTP.Port)
	smtpClient := smtp.NewClient(smtp.Config{
		Host:     acct.SMTP.Host,
		Port:     acct.SMTP.Port,
		TLS:      acct.SMTP.TLS,
		StartTLS: acct.SMTP.StartTLS,
		Insecure: acct.SMTP.Insecure,
		NoTLS:    acct.SMTP.NoTLS,
		Email:    email,
		Password: password,
	})
	if err := smtpClient.TestConnection(); err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("OK")
	}

	return nil
}

// AuthRemoveCmd removes an account.
type AuthRemoveCmd struct {
	Email string `arg:"" help:"Account to remove"`
	// Note: Uses global --force flag for confirmation skip
}

// Run executes the auth remove command.
func (c *AuthRemoveCmd) Run(root *Root) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// TODO: Confirmation prompt if not --force

	if err := cfg.RemoveAccount(c.Email); err != nil {
		return err
	}

	fmt.Printf("Removed account: %s\n", c.Email)
	return nil
}
