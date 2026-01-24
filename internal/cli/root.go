// Package cli defines the command-line interface for sog.
package cli

import (
	"fmt"
	"os"
)

// Root is the top-level CLI structure.
type Root struct {
	// Global flags (matching gog patterns)
	Account string      `help:"Account email to use" env:"SOG_ACCOUNT" short:"a"`
	JSON    bool        `help:"Output JSON to stdout (best for scripting)" xor:"format"`
	Plain   bool        `help:"Output stable, parseable text to stdout (TSV; no colors)" xor:"format"`
	Color   string      `help:"Color output: auto|always|never" default:"auto" enum:"auto,always,never"`
	Force   bool        `help:"Skip confirmations for destructive commands"`
	NoInput bool        `help:"Never prompt; fail instead (useful for CI)" name:"no-input"`
	Verbose bool        `help:"Enable verbose logging" short:"v"`
	Version VersionFlag `name:"version" help:"Print version and exit"`

	// Subcommands
	Auth     AuthCmd     `cmd:"" help:"Manage accounts"`
	Mail     MailCmd     `cmd:"" aliases:"m" help:"Read and send mail"`
	Cal      CalCmd      `cmd:"" aliases:"c" help:"Calendar operations (CalDAV)"`
	Contacts ContactsCmd `cmd:"" aliases:"con" help:"Contact operations (CardDAV)"`
	Tasks    TasksCmd    `cmd:"" aliases:"t" help:"Task operations (CalDAV VTODO)"`
	Drive    DriveCmd    `cmd:"" aliases:"files" help:"File operations (WebDAV)"`
	Invite   InviteCmd   `cmd:"" aliases:"inv" help:"Meeting invitations (iTIP/iMIP)"`
	Folders  FoldersCmd  `cmd:"" aliases:"f" help:"Manage folders"`
	Drafts   DraftsCmd   `cmd:"" aliases:"d" help:"Manage drafts"`
	Idle     IdleCmd     `cmd:"" help:"Watch for new mail (IMAP IDLE)"`
	AIHelp   AIHelpCmd   `cmd:"" name:"ai-help" help:"Show detailed help for AI/LLM agents"`
}

// VersionFlag handles --version.
type VersionFlag string

// BeforeApply prints version and exits.
func (v VersionFlag) BeforeApply() error {
	fmt.Println(v)
	os.Exit(0)
	return nil
}

// AIHelpCmd prints detailed AI-friendly documentation.
type AIHelpCmd struct{}

// Run executes the ai-help command.
func (c *AIHelpCmd) Run(root *Root) error {
	content, err := os.ReadFile("docs/ai-help.md")
	if err != nil {
		// Fallback to embedded help if file not found
		fmt.Println(aiHelpText)
		return nil
	}
	fmt.Println(string(content))
	return nil
}

var aiHelpText = `# sog — AI Help

Standard Open (mail) CLI using IMAP4 and SMTP protocols.

## Authentication

sog auth add <email> [flags]
  --imap-host     IMAP server hostname
  --imap-port     IMAP port (default: 993)
  --smtp-host     SMTP server hostname  
  --smtp-port     SMTP port (default: 587)
  --password      Password (or use keyring)
  --discover      Auto-discover from MX/SRV records

sog auth list
sog auth test [email]
sog auth remove <email>

## Reading Mail

sog mail list [folder]
  --max N         Maximum messages (default: 20)
  --unseen        Only unread messages
  
sog mail get <uid>
  --headers       Headers only
  --raw           Raw RFC822 format

sog mail search <query>
  IMAP SEARCH syntax: FROM, TO, SUBJECT, SINCE, BEFORE, etc.
  Example: sog mail search 'FROM viz SINCE 1-Jan-2026'

## Sending Mail

sog mail send
  --to            Recipient(s)
  --cc            CC recipient(s)
  --bcc           BCC recipient(s)
  --subject       Subject line
  --body          Message body
  --body-file     Read body from file (- for stdin)

## Folders

sog folders list
sog folders create <name>
sog folders delete <name>
sog folders rename <old> <new>

## Output Formats

--json          JSON output
--plain         Plain TSV output
(default)       Human-readable

## Environment Variables

SOG_ACCOUNT     Default account
SOG_PASSWORD_*  Password for account (e.g., SOG_PASSWORD_user_example_com)
`
