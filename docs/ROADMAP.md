# sog Roadmap — Standard Protocol Suite

Extending sog to cover all gog/mog functionality using open standards.

**Legend (RFC 2119):** `!` = MUST, `~` = SHOULD, `≉` = SHOULD NOT, `⊗` = MUST NOT, `?` = MAY

---

## Protocol Mapping

| gog/mog Feature | Standard Protocol | Format | sog Command |
|-----------------|-------------------|--------|-------------|
| Gmail | IMAP4/SMTP | RFC 5322 | `sog mail` ✅ |
| Calendar | CalDAV | iCalendar (.ics) | `sog cal` |
| Contacts | CardDAV | vCard (.vcf) | `sog contacts` |
| Tasks/To-Do | CalDAV | VTODO | `sog tasks` |
| Drive | WebDAV | — | `sog drive` |
| Meet invites | iTIP/iMIP | iCalendar | `sog cal invite` |

---

## Phase 1: Calendar (CalDAV + iCalendar)

### Protocol Stack
- ! Use **CalDAV** (RFC 4791) for calendar access
- ! Use **iCalendar** (RFC 5545) for event format
- ! Use **iTIP** (RFC 5546) for scheduling (invites, responses)
- ! Use **iMIP** (RFC 6047) for email-based invites
- ~ Support **VTIMEZONE** for proper timezone handling
- ? Support **VFREEBUSY** for availability queries

### Commands

```
sog cal list [calendar]           # List events
sog cal get <uid>                 # Get event details
sog cal create                    # Create event
sog cal update <uid>              # Update event
sog cal delete <uid>              # Delete event
sog cal today                     # Today's events
sog cal week                      # This week's events
sog cal search <query>            # Search events
sog cal calendars                 # List calendars
sog cal invite <uid>              # Send meeting invite (iMIP)
sog cal respond <uid> <accept|decline|tentative>
sog cal import <file.ics>         # Import ICS file
sog cal export [calendar]         # Export to ICS
```

### Configuration

```json
{
  "caldav": {
    "url": "https://caldav.example.com/",
    "principal": "/principals/users/you@example.com"
  }
}
```

### Auto-Discovery
- ! Support `.well-known/caldav` (RFC 6764)
- ~ Support SRV records `_caldavs._tcp` and `_caldav._tcp`
- ~ Detect Google Calendar (`https://apidata.googleusercontent.com/caldav/v2/`)
- ~ Detect iCloud (`https://caldav.icloud.com/`)
- ~ Detect Fastmail, Nextcloud, etc.

### Implementation Requirements
- ! Parse/generate valid iCalendar (VEVENT, VTODO, VJOURNAL)
- ! Handle RRULE for recurring events
- ! Handle VTIMEZONE correctly
- ! Support ETAG for sync/conflict detection
- ~ Support calendar-multiget for batch fetch
- ~ Support sync-collection for efficient sync
- ? Support calendar-query for server-side filtering

### Dependencies
```go
require (
    github.com/emersion/go-webdav  // CalDAV client
    github.com/arran4/golang-ical  // iCalendar parser
)
```

---

## Phase 2: Contacts (CardDAV + vCard)

### Protocol Stack
- ! Use **CardDAV** (RFC 6352) for contact access
- ! Use **vCard 4.0** (RFC 6350) for contact format
- ~ Support **vCard 3.0** for legacy compatibility
- ? Support **jCard** (RFC 7095) for JSON output

### Commands

```
sog contacts list                 # List contacts
sog contacts get <uid>            # Get contact details
sog contacts create               # Create contact
sog contacts update <uid>         # Update contact
sog contacts delete <uid>         # Delete contact
sog contacts search <query>       # Search contacts
sog contacts groups               # List groups
sog contacts import <file.vcf>    # Import VCF file
sog contacts export               # Export to VCF
sog contacts merge <uid1> <uid2>  # Merge duplicates
```

### vCard Fields (! = required support)
- ! FN (formatted name)
- ! N (structured name)
- ! EMAIL
- ! TEL
- ~ ADR (address)
- ~ ORG (organization)
- ~ TITLE
- ~ BDAY
- ~ PHOTO (URL or base64)
- ~ NOTE
- ? IMPP (instant messaging)
- ? RELATED
- ? CATEGORIES/GROUPS

### Auto-Discovery
- ! Support `.well-known/carddav` (RFC 6764)
- ~ Support SRV records `_carddavs._tcp` and `_carddav._tcp`
- ~ Detect Google Contacts (via People API or CardDAV)
- ~ Detect iCloud (`https://contacts.icloud.com/`)

### Dependencies
```go
require (
    github.com/emersion/go-webdav  // CardDAV client
    github.com/emersion/go-vcard   // vCard parser
)
```

---

## Phase 3: Tasks (CalDAV + VTODO)

### Protocol Stack
- ! Use **CalDAV** (RFC 4791) with VTODO support
- ! Use **iCalendar VTODO** (RFC 5545) for task format
- ~ Support task lists as separate calendars

### Commands

```
sog tasks list [list]             # List tasks
sog tasks get <uid>               # Get task details
sog tasks add <title>             # Add task
sog tasks done <uid>              # Mark complete
sog tasks undone <uid>            # Mark incomplete
sog tasks delete <uid>            # Delete task
sog tasks due <date>              # Tasks due by date
sog tasks overdue                 # Overdue tasks
sog tasks lists                   # List task lists
sog tasks move <uid> <list>       # Move to list
```

### VTODO Fields (! = required support)
- ! SUMMARY (title)
- ! STATUS (NEEDS-ACTION, IN-PROCESS, COMPLETED, CANCELLED)
- ! UID
- ~ DUE
- ~ DTSTART
- ~ PRIORITY (1-9)
- ~ PERCENT-COMPLETE
- ~ CATEGORIES
- ~ DESCRIPTION
- ? RELATED-TO (subtasks)
- ? RRULE (recurring tasks)

### Status Mapping
```
gog/mog         VTODO STATUS        sog
--------        ------------        ---
needsAction     NEEDS-ACTION        pending
inProgress      IN-PROCESS          active
completed       COMPLETED           done
cancelled       CANCELLED           cancelled
```

---

## Phase 4: Files (WebDAV)

### Protocol Stack
- ! Use **WebDAV** (RFC 4918) for file access
- ~ Support **PROPFIND** for metadata
- ~ Support **COPY**, **MOVE** for file operations
- ? Support **LOCK** for collaborative editing
- ? Support **SEARCH** (RFC 5323) if available

### Commands

```
sog drive list [path]             # List files/folders
sog drive get <path>              # Download file
sog drive put <local> <remote>    # Upload file
sog drive mkdir <path>            # Create folder
sog drive rm <path>               # Delete file/folder
sog drive mv <src> <dst>          # Move/rename
sog drive cp <src> <dst>          # Copy
sog drive info <path>             # File metadata
sog drive search <query>          # Search (if supported)
sog drive share <path>            # Get share link (provider-specific)
```

### Auto-Discovery
- ~ Support `.well-known/webdav`
- ~ Detect Nextcloud, ownCloud
- ~ Detect Box, Dropbox (if WebDAV supported)

### Dependencies
```go
require (
    github.com/emersion/go-webdav  // WebDAV client
)
```

---

## Phase 5: Meeting Invites (iTIP/iMIP)

### Protocol Stack
- ! Use **iTIP** (RFC 5546) for scheduling semantics
- ! Use **iMIP** (RFC 6047) for email transport
- ! Generate valid METHOD:REQUEST, REPLY, CANCEL
- ~ Parse incoming invites from email

### Workflow

```
# Send invite
sog cal create --attendees "a@ex.com,b@ex.com" --send-invite

# Accept invite from email
sog cal respond <uid> accept

# Decline
sog cal respond <uid> decline --message "Sorry, conflict"

# Propose new time
sog cal respond <uid> counter --time "2026-01-25T15:00:00"

# Cancel meeting (organizer)
sog cal cancel <uid> --message "Meeting cancelled"
```

### iMIP Email Format
- ! Include `text/calendar` MIME part
- ! Set `method=REQUEST|REPLY|CANCEL`
- ~ Include human-readable body
- ~ Support `RSVP=TRUE` in ATTENDEE

---

## Architecture

### Unified Auth

```
sog auth add <email> [flags]
  --discover          Auto-discover all protocols
  --imap-host         IMAP host
  --smtp-host         SMTP host
  --caldav-url        CalDAV URL
  --carddav-url       CardDAV URL
  --webdav-url        WebDAV URL
```

### Config Structure

```json
{
  "accounts": {
    "you@example.com": {
      "imap": { "host": "...", "port": 993 },
      "smtp": { "host": "...", "port": 587 },
      "caldav": { "url": "https://caldav.example.com/" },
      "carddav": { "url": "https://carddav.example.com/" },
      "webdav": { "url": "https://webdav.example.com/" }
    }
  }
}
```

### Package Structure

```
sog/
├── cmd/sog/
├── internal/
│   ├── cli/
│   │   ├── auth.go
│   │   ├── mail.go      ✅
│   │   ├── cal.go       (Phase 1)
│   │   ├── contacts.go  (Phase 2)
│   │   ├── tasks.go     (Phase 3)
│   │   └── drive.go     (Phase 4)
│   ├── config/          ✅
│   ├── discover/        ✅ (extend)
│   ├── imap/            ✅
│   ├── smtp/            ✅
│   ├── caldav/          (Phase 1)
│   ├── carddav/         (Phase 2)
│   ├── webdav/          (Phase 4)
│   ├── ical/            (Phase 1, 3)
│   └── vcard/           (Phase 2)
```

---

## Implementation Priority

### Phase 1: Calendar (HIGH)
- Most requested feature after email
- CalDAV widely supported (Google, iCloud, Fastmail, Nextcloud)
- Foundation for Tasks (VTODO)

### Phase 2: Contacts (MEDIUM)
- CardDAV well-supported
- Useful for address book integration
- Enables contact autocomplete in mail

### Phase 3: Tasks (MEDIUM)
- Reuses CalDAV infrastructure
- VTODO less universally supported
- Some providers use proprietary APIs

### Phase 4: Files (LOW)
- WebDAV support varies
- Many providers prefer proprietary APIs
- Large file handling complexity

### Phase 5: Invites (LOW)
- Depends on Phase 1 (Calendar)
- iTIP/iMIP complex to implement correctly
- Email transport adds complexity

---

## Testing Strategy

### Unit Tests
- ! ≥85% coverage on parsers (iCalendar, vCard)
- ! ≥85% coverage on protocol logic
- ~ Mock WebDAV/CalDAV/CardDAV servers

### Integration Tests
- ! Test against Nextcloud (Docker)
- ~ Test against Radicale (Docker)
- ~ Test against Baikal (Docker)
- ? Test against Google (manual, rate-limited)
- ? Test against iCloud (manual)

### Test Servers (Docker)

```bash
# Radicale (CalDAV/CardDAV)
docker run -d -p 5232:5232 tomsquest/docker-radicale

# Nextcloud (CalDAV/CardDAV/WebDAV)
docker run -d -p 8080:80 nextcloud

# Baikal (CalDAV/CardDAV)
docker run -d -p 80:80 ckulka/baikal
```

---

## Timeline Estimate

| Phase | Feature | Effort | Dependencies |
|-------|---------|--------|--------------|
| 1 | Calendar | 2-3 days | go-webdav, golang-ical |
| 2 | Contacts | 1-2 days | go-webdav, go-vcard |
| 3 | Tasks | 1 day | Phase 1 |
| 4 | Files | 1-2 days | go-webdav |
| 5 | Invites | 2 days | Phase 1, SMTP |

**Total: ~8-10 days for full gog parity**

---

## Open Questions

1. **Google Calendar API vs CalDAV?**
   - Google's CalDAV is limited; API is richer
   - ~ Prefer CalDAV for universality
   - ? Add `--use-api` flag for provider-specific features

2. **Offline sync?**
   - ~ Support local cache with ETAG-based sync
   - ? Full offline mode with conflict resolution

3. **Shared calendars/address books?**
   - ~ Support listing shared resources
   - ? Support delegation (CalDAV ACL)

4. **Attachments?**
   - ~ Support ATTACH in VEVENT
   - ? Support inline vs URL attachments

---

## References

- [RFC 4791](https://tools.ietf.org/html/rfc4791) — CalDAV
- [RFC 5545](https://tools.ietf.org/html/rfc5545) — iCalendar
- [RFC 5546](https://tools.ietf.org/html/rfc5546) — iTIP
- [RFC 6047](https://tools.ietf.org/html/rfc6047) — iMIP
- [RFC 6350](https://tools.ietf.org/html/rfc6350) — vCard 4.0
- [RFC 6352](https://tools.ietf.org/html/rfc6352) — CardDAV
- [RFC 4918](https://tools.ietf.org/html/rfc4918) — WebDAV
- [RFC 6764](https://tools.ietf.org/html/rfc6764) — CalDAV/CardDAV SRV Discovery
