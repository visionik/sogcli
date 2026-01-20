# sog — AI Help

Standard Open (mail) CLI using IMAP4 and SMTP protocols.
A gog-compatible alternative using standard protocols instead of Google APIs.

## Quick Start

```bash
# Add account with auto-discovery
sog auth add you@gmail.com --discover --password "your-app-password"

# List messages
sog mail list

# Send a message
sog mail send --to recipient@example.com --subject "Hi" --body "Hello!"
```

## Authentication

```bash
sog auth add <email> [flags]
  --imap-host     IMAP server hostname
  --imap-port     IMAP port (default: 993)
  --smtp-host     SMTP server hostname  
  --smtp-port     SMTP port (default: 587)
  --password      Password (will prompt if not provided)
  --discover      Auto-discover servers from MX/SRV records
  --insecure      Skip TLS certificate verification
  --no-tls        Disable TLS (plain text connection)

sog auth list           # List configured accounts
sog auth test [email]   # Test IMAP/SMTP connection
sog auth remove <email> # Remove account
```

## Reading Mail

```bash
sog mail list [folder]
  --max N         Maximum messages (default: 20)
  --unseen        Only unread messages
  
sog mail get <uid>
  --headers       Headers only
  --raw           Raw RFC822 format
  --folder        Folder containing message (default: INBOX)

sog mail search <query>
  IMAP SEARCH syntax: FROM, TO, SUBJECT, SINCE, BEFORE, TEXT, etc.
  Examples:
    sog mail search 'FROM viz'
    sog mail search 'SUBJECT meeting SINCE 1-Jan-2026'
    sog mail search 'UNSEEN'
    sog mail search 'ALL'
```

## Sending Mail

```bash
sog mail send
  --to            Recipients (comma-separated, required)
  --cc            CC recipients (comma-separated)
  --bcc           BCC recipients (comma-separated)
  --subject       Subject line (required)
  --body          Message body
  --body-file     Read body from file (- for stdin)

sog mail reply <uid> --body "Reply text"
  --all           Reply to all recipients
  --folder        Folder containing original (default: INBOX)

sog mail forward <uid> --to recipient@example.com
  --body          Additional message to prepend
  --folder        Folder containing original (default: INBOX)
```

## Organizing Mail

```bash
sog mail move <uid> <folder>    # Move message
sog mail copy <uid> <folder>    # Copy message
sog mail flag <uid> <flag>      # Add flag (seen, flagged, answered, deleted, draft)
sog mail unflag <uid> <flag>    # Remove flag
sog mail delete <uid>           # Delete message (mark deleted + expunge)

sog folders list                # List all folders
sog folders create <name>       # Create folder
sog folders delete <name>       # Delete folder
sog folders rename <old> <new>  # Rename folder
```

## Drafts

```bash
sog drafts list                 # List drafts
sog drafts create               # Create draft
  --to            Recipients
  --subject       Subject
  --body          Body text
sog drafts send <uid>           # Send a draft
sog drafts delete <uid>         # Delete a draft
```

## IMAP IDLE (Watch for New Mail)

```bash
sog idle
  --folder        Folder to watch (default: INBOX)
  --exec          Command to run on new mail
```

## Output Formats

```bash
--json          JSON output (for scripting)
--plain         Plain TSV output (no colors)
(default)       Human-readable
```

## Environment Variables

```
SOG_ACCOUNT              Default account email
SOG_PASSWORD_<email>     Password for account (e.g., SOG_PASSWORD_USER_EXAMPLE_COM)
```

## Configuration

Config file: `~/.config/sog/config.json`
Credentials: OS keyring (macOS Keychain, etc.)

## Examples for AI Agents

```bash
# Check for unread mail
sog mail list --unseen --json

# Send notification email
sog mail send --to alerts@example.com --subject "Alert" --body "Something happened"

# Search recent messages from a sender
sog mail search 'FROM boss@company.com SINCE 7-Jan-2026' --json

# Watch for new mail and run a command
sog idle --exec 'echo "New mail arrived"'

# Forward a message
sog mail forward 12345 --to backup@example.com --body "Archiving this"
```
