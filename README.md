# sog

**Standards Ops Gadget** — CLI for IMAP/SMTP/CalDAV/CardDAV/WebDAV

The open-standards alternative to [gog](https://github.com/steipete/gog) (Google) and mog (Microsoft).

```bash
# Mail
sog mail list
sog mail get 1234
sog mail send --to user@example.com --subject "Hello" --body "Hi there"

# Calendar
sog cal today
sog cal create "Team Meeting" --start "2024-01-25T14:00" --duration 1h

# Contacts
sog contacts search "John"
sog contacts create "John Doe" -e john@example.com -p 555-1234

# Tasks
sog tasks list
sog tasks add "Review PR" --due 2024-01-26 -p 1

# Files
sog drive ls
sog drive upload report.pdf /documents/
sog drive download /documents/report.pdf

# Meeting Invites
sog invite send "Project Kickoff" user@example.com --start "2024-01-25T10:00"
```

## Features

| Feature | Protocol | Commands |
|---------|----------|----------|
| **Mail** | IMAP/SMTP | `list`, `get`, `search`, `send`, `reply`, `forward`, `move`, `delete` |
| **Calendar** | CalDAV | `list`, `get`, `search`, `today`, `week`, `create`, `update`, `delete` |
| **Contacts** | CardDAV | `list`, `get`, `search`, `create`, `update`, `delete` |
| **Tasks** | CalDAV VTODO | `list`, `add`, `update`, `done`, `undo`, `delete`, `clear`, `due`, `overdue` |
| **Files** | WebDAV | `ls`, `get`, `download`, `upload`, `mkdir`, `delete`, `move`, `copy`, `cat` |
| **Invites** | iTIP/iMIP | `send`, `reply`, `cancel`, `parse`, `preview` |

## Installation

### From Source

```bash
go install github.com/visionik/sogcli/cmd/sog@latest
```

### Build from Source

```bash
git clone https://github.com/visionik/sogcli.git
cd sogcli
go build -o sog ./cmd/sog
```

## Quick Start

### 1. Add an Account

```bash
# Auto-discover servers from DNS
sog auth add you@fastmail.com --discover

# Or specify manually
sog auth add you@example.com \
  --imap-host imap.example.com \
  --smtp-host smtp.example.com \
  --caldav-url https://caldav.example.com/ \
  --carddav-url https://carddav.example.com/ \
  --webdav-url https://webdav.example.com/
```

### 2. Test Connection

```bash
sog auth test
```

### 3. Start Using

```bash
sog mail list
sog cal today
sog contacts list
```

## Configuration

Config file: `~/.config/sog/config.json`

Passwords are stored securely in your system keychain.

### Environment Variables

- `SOG_ACCOUNT` — Default account email

### Global Flags

```
--ai-help        Detailed help for AI/LLM agents
--account, -a    Account to use
--json           JSON output (for scripting)
--plain          TSV output (parseable)
--force          Skip confirmations
--no-input       Never prompt (CI mode)
--verbose, -v    Debug logging
```

## Tested Providers

| Provider | IMAP/SMTP | CalDAV | CardDAV | WebDAV |
|----------|-----------|--------|---------|--------|
| **Fastmail** | ✅ | ✅ | ✅ | ✅ |

*Other standards-compliant providers should work but have not been tested yet.*

## Comparison with gog/mog

| | sog | gog | mog |
|---|-----|-----|-----|
| **Provider** | Any (open standards) | Google only | Microsoft only |
| **Auth** | App passwords | OAuth 2.0 | OAuth 2.0 |
| **Mail** | ✅ IMAP/SMTP | ✅ Gmail API | ✅ Graph API |
| **Calendar** | ✅ CalDAV | ✅ Calendar API | ✅ Graph API |
| **Contacts** | ✅ CardDAV | ✅ People API | ✅ Graph API |
| **Tasks** | ✅ CalDAV VTODO | ✅ Tasks API | ✅ To Do API |
| **Files** | ✅ WebDAV | ✅ Drive API | ✅ OneDrive API |
| **Sheets** | — | ✅ | — |
| **Docs** | — | ✅ | — |

## AI/LLM Integration

sog is designed for AI agents and automation:

```bash
# Detailed help for AI agents
sog --ai-help

# Structured output
sog mail list --json
sog cal today --json
```

## License

MIT

## Contributing

Contributions welcome! Please open an issue or PR.

## Credits

Created by [visionik@pobox.com](mailto:visionik@pobox.com) and Vinston ([Clawdbot](https://github.com/clawdbot/clawdbot)) using the visionik.md coding skill.

Part of the **Ops Gadget** family:
- **gog** — Google Ops Gadget
- **mog** — Microsoft Ops Gadget  
- **sog** — Standards Ops Gadget
