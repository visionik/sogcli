# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-01-24

### Added

#### Mail (IMAP/SMTP)
- `sog mail list` — List messages in a folder
- `sog mail get` — Get message by UID
- `sog mail search` — Search messages
- `sog mail send` — Send a message
- `sog mail reply` — Reply to a message
- `sog mail forward` — Forward a message
- `sog mail move` — Move message to folder
- `sog mail copy` — Copy message to folder
- `sog mail flag` — Set message flag
- `sog mail unflag` — Remove message flag
- `sog mail delete` — Delete a message
- `sog folders list/create/delete/rename` — Folder management
- `sog drafts list/create/send/delete` — Draft management
- `sog idle` — Watch for new mail (IMAP IDLE)

#### Calendar (CalDAV)
- `sog cal list` — List events
- `sog cal get` — Get event details
- `sog cal search` — Search events
- `sog cal today` — Today's events
- `sog cal week` — This week's events
- `sog cal create` — Create an event
- `sog cal update` — Update an event
- `sog cal delete` — Delete an event
- `sog cal calendars` — List calendars

#### Contacts (CardDAV)
- `sog contacts list` — List contacts
- `sog contacts get` — Get contact details
- `sog contacts search` — Search contacts
- `sog contacts create` — Create a contact
- `sog contacts update` — Update a contact
- `sog contacts delete` — Delete a contact
- `sog contacts books` — List address books

#### Tasks (CalDAV VTODO)
- `sog tasks list` — List tasks
- `sog tasks add` — Add a task
- `sog tasks get` — Get task details
- `sog tasks update` — Update a task
- `sog tasks done` — Mark task complete
- `sog tasks undo` — Mark task incomplete
- `sog tasks delete` — Delete a task
- `sog tasks clear` — Clear completed tasks
- `sog tasks due` — Tasks due by date
- `sog tasks overdue` — Overdue tasks
- `sog tasks lists` — List task lists

#### Files (WebDAV)
- `sog drive ls` — List files and folders
- `sog drive get` — Get file metadata
- `sog drive download` — Download a file
- `sog drive upload` — Upload a file
- `sog drive mkdir` — Create a directory
- `sog drive delete` — Delete file or directory
- `sog drive move` — Move/rename file
- `sog drive copy` — Copy file
- `sog drive cat` — Output file contents

#### Meeting Invites (iTIP/iMIP)
- `sog invite send` — Send meeting invitation
- `sog invite reply` — Reply to invitation (accept/decline/tentative)
- `sog invite cancel` — Cancel a meeting
- `sog invite parse` — Parse .ics file
- `sog invite preview` — Preview invite without sending

#### Authentication
- `sog auth add` — Add account with auto-discovery
- `sog auth list` — List configured accounts
- `sog auth test` — Test account connection
- `sog auth remove` — Remove an account
- `sog auth password` — Set protocol-specific passwords
- Secure keychain storage for passwords
- Support for separate passwords per protocol (IMAP, SMTP, CalDAV, CardDAV, WebDAV)

#### Output Formats
- `--json` — JSON output for scripting
- `--plain` — TSV output for parsing
- Colored terminal output (auto-detected)

#### AI/LLM Support
- `sog ai-help` — Detailed help for AI agents
- Structured output modes for automation

### Tested Providers
- Fastmail (full support)
- Works with any standards-compliant provider

[0.1.0]: https://github.com/visionik/sogcli/releases/tag/v0.1.0
