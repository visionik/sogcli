# sog — Standard Open (mail) CLI

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A **gog-compatible** CLI for email using standard IMAP4 and SMTP protocols. No vendor lock-in, works with any mail provider.

## Why sog?

| Feature | gog (Google) | mog (Microsoft) | **sog (Standard)** |
|---------|--------------|-----------------|---------------------|
| Protocol | Google APIs | MS Graph | IMAP4/SMTP |
| Providers | Gmail only | Outlook only | **Any provider** |
| Auth | OAuth | OAuth | App passwords |
| Vendor lock-in | Yes | Yes | **No** |

## Install

```bash
go install github.com/visionik/sogcli/cmd/sog@latest
```

Or build from source:
```bash
git clone https://github.com/visionik/sogcli.git
cd sogcli
go build -o sog ./cmd/sog
```

## Quick Start

```bash
# Add account (auto-discovers Gmail, Outlook, etc.)
sog auth add you@gmail.com --discover --password "your-app-password"

# List messages
sog mail list

# Send a message
sog mail send --to friend@example.com --subject "Hello" --body "Hi from sog!"

# Reply to a message
sog mail reply 12345 --body "Thanks!"

# Watch for new mail
sog idle
```

## Features

### Authentication
```bash
sog auth add <email> [flags]    # Add account (supports --discover)
sog auth list                   # List accounts
sog auth test [email]           # Test connection
sog auth remove <email>         # Remove account
```

### Mail Operations
```bash
sog mail list [folder]          # List messages
sog mail get <uid>              # Get message by UID
sog mail search <query>         # Search (IMAP syntax)
sog mail send                   # Send message
sog mail reply <uid>            # Reply to message
sog mail forward <uid>          # Forward message
sog mail move <uid> <folder>    # Move message
sog mail copy <uid> <folder>    # Copy message
sog mail flag <uid> <flag>      # Set flag (seen, flagged, etc.)
sog mail unflag <uid> <flag>    # Remove flag
sog mail delete <uid>           # Delete message
```

### Folders
```bash
sog folders list                # List folders
sog folders create <name>       # Create folder
sog folders delete <name>       # Delete folder
sog folders rename <old> <new>  # Rename folder
```

### Drafts
```bash
sog drafts list                 # List drafts
sog drafts create               # Create draft
sog drafts send <uid>           # Send draft
sog drafts delete <uid>         # Delete draft
```

### IMAP IDLE
```bash
sog idle                        # Watch for new mail
sog idle --exec 'notify-send "New mail!"'
```

## Search Syntax

sog uses IMAP SEARCH syntax:

```bash
sog mail search 'FROM sender@example.com'
sog mail search 'SUBJECT meeting'
sog mail search 'SINCE 1-Jan-2026'
sog mail search 'FROM boss SUBJECT urgent UNSEEN'
sog mail search 'ALL'
```

## Global Flags

```
-a, --account     Account to use (or set SOG_ACCOUNT)
    --json        JSON output
    --plain       TSV output (no colors)
    --force       Skip confirmations
    --no-input    Never prompt (CI mode)
-v, --verbose     Verbose output
```

## Configuration

- **Config file:** `~/.config/sog/config.json`
- **Credentials:** OS keyring (macOS Keychain, Windows Credential Manager, etc.)
- **Env vars:** `SOG_ACCOUNT`, `SOG_PASSWORD_<email>`

## Provider Setup

### Gmail / Google Workspace
1. Enable 2-Step Verification: https://myaccount.google.com/security
2. Generate App Password: https://myaccount.google.com/apppasswords
3. Add account:
```bash
sog auth add you@gmail.com --discover --password "xxxx xxxx xxxx xxxx"
```

### Outlook / Microsoft 365
```bash
sog auth add you@outlook.com --discover --password "your-app-password"
```

### Other Providers
```bash
sog auth add you@example.com \
  --imap-host imap.example.com --imap-port 993 \
  --smtp-host smtp.example.com --smtp-port 587 \
  --password "your-password"
```

## For AI Agents

```bash
sog ai-help    # Detailed documentation for LLM integration
```

## Development

```bash
task build          # Build binary
task test           # Run unit tests
task test:integration  # Run integration tests (requires Docker)
task test:all       # All tests with coverage
task check          # Pre-commit checks
```

## License

MIT

## Credits

Built with:
- [emersion/go-imap](https://github.com/emersion/go-imap) — IMAP client
- [emersion/go-smtp](https://github.com/emersion/go-smtp) — SMTP client
- [alecthomas/kong](https://github.com/alecthomas/kong) — CLI framework
- [zalando/go-keyring](https://github.com/zalando/go-keyring) — Credential storage
