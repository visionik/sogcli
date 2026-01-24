// Package caldav provides a CalDAV client for calendar operations.
package caldav

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
)

// Client wraps a CalDAV client with convenience methods.
type Client struct {
	client   *caldav.Client
	email    string
	calURL   string
	calendar *caldav.Calendar
}

// Config holds CalDAV connection configuration.
type Config struct {
	URL      string // CalDAV server URL
	Email    string // Account email (for auth)
	Password string // Account password
}

// Event represents a calendar event.
type Event struct {
	UID         string    `json:"uid"`
	Summary     string    `json:"summary"`
	Description string    `json:"description,omitempty"`
	Location    string    `json:"location,omitempty"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	AllDay      bool      `json:"all_day,omitempty"`
	Organizer   string    `json:"organizer,omitempty"`
	Attendees   []string  `json:"attendees,omitempty"`
	Status      string    `json:"status,omitempty"`
	URL         string    `json:"url,omitempty"`
	ETag        string    `json:"etag,omitempty"`
}

// Calendar represents a calendar.
type Calendar struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

// Connect establishes a connection to a CalDAV server.
func Connect(cfg Config) (*Client, error) {
	httpClient := webdav.HTTPClientWithBasicAuth(http.DefaultClient, cfg.Email, cfg.Password)

	client, err := caldav.NewClient(httpClient, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create CalDAV client: %w", err)
	}

	return &Client{
		client: client,
		email:  cfg.Email,
		calURL: cfg.URL,
	}, nil
}

// Close closes the client connection.
func (c *Client) Close() error {
	// CalDAV is stateless HTTP, nothing to close
	return nil
}

// FindCalendars discovers available calendars.
func (c *Client) FindCalendars(ctx context.Context) ([]Calendar, error) {
	// First find the calendar home set
	principal, err := c.client.FindCurrentUserPrincipal(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find principal: %w", err)
	}

	homeSet, err := c.client.FindCalendarHomeSet(ctx, principal)
	if err != nil {
		return nil, fmt.Errorf("failed to find calendar home set: %w", err)
	}

	cals, err := c.client.FindCalendars(ctx, homeSet)
	if err != nil {
		return nil, fmt.Errorf("failed to find calendars: %w", err)
	}

	result := make([]Calendar, 0, len(cals))
	for _, cal := range cals {
		result = append(result, Calendar{
			Path:        cal.Path,
			Name:        cal.Name,
			Description: cal.Description,
		})
	}
	return result, nil
}

// ListEvents retrieves events from a calendar within a time range.
func (c *Client) ListEvents(ctx context.Context, calPath string, start, end time.Time) ([]Event, error) {
	query := &caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name: "VCALENDAR",
			Props: []string{
				"VERSION",
			},
			Comps: []caldav.CalendarCompRequest{{
				Name: "VEVENT",
				Props: []string{
					"UID",
					"SUMMARY",
					"DESCRIPTION",
					"LOCATION",
					"DTSTART",
					"DTEND",
					"ORGANIZER",
					"ATTENDEE",
					"STATUS",
					"URL",
				},
			}},
		},
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{{
				Name:  "VEVENT",
				Start: start,
				End:   end,
			}},
		},
	}

	objects, err := c.client.QueryCalendar(ctx, calPath, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query calendar: %w", err)
	}

	events := make([]Event, 0, len(objects))
	for _, obj := range objects {
		event, err := parseICalEvent(obj.Data)
		if err != nil {
			continue // Skip malformed events
		}
		event.ETag = obj.ETag
		events = append(events, *event)
	}

	return events, nil
}

// GetEvent retrieves a single event by UID.
func (c *Client) GetEvent(ctx context.Context, calPath, uid string) (*Event, error) {
	// Query for the specific UID
	query := &caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name: "VCALENDAR",
			Comps: []caldav.CalendarCompRequest{{
				Name: "VEVENT",
			}},
		},
		CompFilter: caldav.CompFilter{
			Name: "VCALENDAR",
			Comps: []caldav.CompFilter{{
				Name: "VEVENT",
				Props: []caldav.PropFilter{{
					Name:      "UID",
					TextMatch: &caldav.TextMatch{Text: uid},
				}},
			}},
		},
	}

	objects, err := c.client.QueryCalendar(ctx, calPath, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query calendar: %w", err)
	}

	if len(objects) == 0 {
		return nil, fmt.Errorf("event not found: %s", uid)
	}

	event, err := parseICalEvent(objects[0].Data)
	if err != nil {
		return nil, err
	}
	event.ETag = objects[0].ETag
	return event, nil
}

// CreateEvent creates a new event.
func (c *Client) CreateEvent(ctx context.Context, calPath string, event *Event) error {
	cal := createICalEvent(event)
	_, err := c.client.PutCalendarObject(ctx, calPath+"/"+event.UID+".ics", cal)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

// UpdateEvent updates an existing event.
func (c *Client) UpdateEvent(ctx context.Context, calPath string, event *Event) error {
	cal := createICalEvent(event)
	_, err := c.client.PutCalendarObject(ctx, calPath+"/"+event.UID+".ics", cal)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	return nil
}

// DeleteEvent deletes an event.
func (c *Client) DeleteEvent(ctx context.Context, calPath, uid string) error {
	err := c.client.RemoveAll(ctx, calPath+"/"+uid+".ics")
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}

// parseICalEvent parses an iCalendar VEVENT into an Event.
func parseICalEvent(cal *ical.Calendar) (*Event, error) {
	if cal == nil {
		return nil, fmt.Errorf("nil calendar data")
	}

	for _, child := range cal.Children {
		if child.Name != ical.CompEvent {
			continue
		}

		event := &Event{}

		// UID
		if prop := child.Props.Get(ical.PropUID); prop != nil {
			event.UID = prop.Value
		}

		// Summary
		if prop := child.Props.Get(ical.PropSummary); prop != nil {
			event.Summary = prop.Value
		}

		// Description
		if prop := child.Props.Get(ical.PropDescription); prop != nil {
			event.Description = prop.Value
		}

		// Location
		if prop := child.Props.Get(ical.PropLocation); prop != nil {
			event.Location = prop.Value
		}

		// Organizer
		if prop := child.Props.Get(ical.PropOrganizer); prop != nil {
			event.Organizer = strings.TrimPrefix(prop.Value, "mailto:")
		}

		// Status
		if prop := child.Props.Get(ical.PropStatus); prop != nil {
			event.Status = prop.Value
		}

		// URL
		if prop := child.Props.Get(ical.PropURL); prop != nil {
			event.URL = prop.Value
		}

		// Start time
		if prop := child.Props.Get(ical.PropDateTimeStart); prop != nil {
			t, err := prop.DateTime(time.Local)
			if err == nil {
				event.Start = t
			}
			// Check for all-day event (VALUE=DATE)
			if param := prop.Params.Get(ical.ParamValue); param == "DATE" {
				event.AllDay = true
			}
		}

		// End time
		if prop := child.Props.Get(ical.PropDateTimeEnd); prop != nil {
			t, err := prop.DateTime(time.Local)
			if err == nil {
				event.End = t
			}
		} else if prop := child.Props.Get(ical.PropDuration); prop != nil {
			// Handle DURATION property instead of DTEND
			dur, err := prop.Duration()
			if err == nil && !event.Start.IsZero() {
				event.End = event.Start.Add(dur)
			}
		}

		// Attendees
		for _, prop := range child.Props[ical.PropAttendee] {
			attendee := strings.TrimPrefix(prop.Value, "mailto:")
			event.Attendees = append(event.Attendees, attendee)
		}

		return event, nil
	}

	return nil, fmt.Errorf("no VEVENT found in calendar data")
}

// createICalEvent creates an iCalendar from an Event.
func createICalEvent(event *Event) *ical.Calendar {
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//sog//CalDAV Client//EN")

	vevent := ical.NewComponent(ical.CompEvent)
	vevent.Props.SetText(ical.PropUID, event.UID)
	vevent.Props.SetText(ical.PropSummary, event.Summary)

	if event.Description != "" {
		vevent.Props.SetText(ical.PropDescription, event.Description)
	}
	if event.Location != "" {
		vevent.Props.SetText(ical.PropLocation, event.Location)
	}
	if event.URL != "" {
		vevent.Props.SetText(ical.PropURL, event.URL)
	}
	if event.Status != "" {
		vevent.Props.SetText(ical.PropStatus, event.Status)
	}

	// Start and end times
	if event.AllDay {
		startProp := ical.NewProp(ical.PropDateTimeStart)
		startProp.SetDate(event.Start)
		vevent.Props.Set(startProp)

		endProp := ical.NewProp(ical.PropDateTimeEnd)
		endProp.SetDate(event.End)
		vevent.Props.Set(endProp)
	} else {
		startProp := ical.NewProp(ical.PropDateTimeStart)
		startProp.SetDateTime(event.Start)
		vevent.Props.Set(startProp)

		endProp := ical.NewProp(ical.PropDateTimeEnd)
		endProp.SetDateTime(event.End)
		vevent.Props.Set(endProp)
	}

	if event.Organizer != "" {
		vevent.Props.SetText(ical.PropOrganizer, "mailto:"+event.Organizer)
	}

	for _, attendee := range event.Attendees {
		prop := ical.NewProp(ical.PropAttendee)
		prop.Value = "mailto:" + attendee
		vevent.Props.Add(prop)
	}

	// DTSTAMP is required
	dtstamp := ical.NewProp(ical.PropDateTimeStamp)
	dtstamp.SetDateTime(time.Now().UTC())
	vevent.Props.Set(dtstamp)

	cal.Children = append(cal.Children, vevent)
	return cal
}
