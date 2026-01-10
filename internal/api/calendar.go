package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// CalendarEvent represents a Canvas calendar event
type CalendarEvent struct {
	ID                   int64           `json:"id"`
	Title                string          `json:"title"`
	StartAt              *time.Time      `json:"start_at,omitempty"`
	EndAt                *time.Time      `json:"end_at,omitempty"`
	Description          string          `json:"description,omitempty"`
	LocationName         string          `json:"location_name,omitempty"`
	LocationAddress      string          `json:"location_address,omitempty"`
	ContextCode          string          `json:"context_code"`
	EffectiveContextCode string          `json:"effective_context_code,omitempty"`
	ContextName          string          `json:"context_name,omitempty"`
	AllContextCodes      string          `json:"all_context_codes,omitempty"`
	WorkflowState        string          `json:"workflow_state"`
	Hidden               bool            `json:"hidden"`
	ParentEventID        *int64          `json:"parent_event_id,omitempty"`
	ChildEventsCount     int             `json:"child_events_count"`
	ChildEvents          []CalendarEvent `json:"child_events,omitempty"`
	URL                  string          `json:"url,omitempty"`
	HTMLURL              string          `json:"html_url,omitempty"`
	AllDayDate           string          `json:"all_day_date,omitempty"`
	AllDay               bool            `json:"all_day"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
	AppointmentGroupID   *int64          `json:"appointment_group_id,omitempty"`
	AppointmentGroupURL  string          `json:"appointment_group_url,omitempty"`
	OwnReservation       bool            `json:"own_reservation"`
	ReserveURL           string          `json:"reserve_url,omitempty"`
	Reserved             bool            `json:"reserved"`
	ParticipantType      string          `json:"participant_type,omitempty"`
	ParticipantsPerAppt  *int            `json:"participants_per_appointment,omitempty"`
	AvailableSlots       *int            `json:"available_slots,omitempty"`
	User                 *User           `json:"user,omitempty"`
	Group                interface{}     `json:"group,omitempty"`
	ImportantDates       bool            `json:"important_dates"`
	SeriesUUID           string          `json:"series_uuid,omitempty"`
	RRule                string          `json:"rrule,omitempty"`
	SeriesHead           *bool           `json:"series_head,omitempty"`
	SeriesNaturalLang    string          `json:"series_natural_language,omitempty"`
	BlackoutDate         bool            `json:"blackout_date"`
}

// CalendarService handles calendar-related API calls
type CalendarService struct {
	client *Client
}

// NewCalendarService creates a new calendar service
func NewCalendarService(client *Client) *CalendarService {
	return &CalendarService{client: client}
}

// ListCalendarEventsOptions holds options for listing calendar events
type ListCalendarEventsOptions struct {
	Type           string   // event, assignment, sub_assignment
	StartDate      string   // yyyy-mm-dd or ISO 8601
	EndDate        string   // yyyy-mm-dd or ISO 8601
	Undated        bool
	AllEvents      bool
	ContextCodes   []string // course_123, user_456, etc.
	Excludes       []string // description, child_events, assignment
	Includes       []string // web_conference, series_natural_language
	ImportantDates bool
	BlackoutDate   bool
	Page           int
	PerPage        int
}

// List retrieves calendar events
func (s *CalendarService) List(ctx context.Context, opts *ListCalendarEventsOptions) ([]CalendarEvent, error) {
	path := "/api/v1/calendar_events"

	if opts != nil {
		query := url.Values{}

		if opts.Type != "" {
			query.Add("type", opts.Type)
		}

		if opts.StartDate != "" {
			query.Add("start_date", opts.StartDate)
		}

		if opts.EndDate != "" {
			query.Add("end_date", opts.EndDate)
		}

		if opts.Undated {
			query.Add("undated", "true")
		}

		if opts.AllEvents {
			query.Add("all_events", "true")
		}

		for _, code := range opts.ContextCodes {
			query.Add("context_codes[]", code)
		}

		for _, exc := range opts.Excludes {
			query.Add("excludes[]", exc)
		}

		for _, inc := range opts.Includes {
			query.Add("includes[]", inc)
		}

		if opts.ImportantDates {
			query.Add("important_dates", "true")
		}

		if opts.BlackoutDate {
			query.Add("blackout_date", "true")
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var events []CalendarEvent
	if err := s.client.GetAllPages(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// ListForUser retrieves calendar events for a specific user
func (s *CalendarService) ListForUser(ctx context.Context, userID int64, opts *ListCalendarEventsOptions) ([]CalendarEvent, error) {
	path := fmt.Sprintf("/api/v1/users/%d/calendar_events", userID)

	if opts != nil {
		query := url.Values{}

		if opts.Type != "" {
			query.Add("type", opts.Type)
		}

		if opts.StartDate != "" {
			query.Add("start_date", opts.StartDate)
		}

		if opts.EndDate != "" {
			query.Add("end_date", opts.EndDate)
		}

		if opts.Undated {
			query.Add("undated", "true")
		}

		if opts.AllEvents {
			query.Add("all_events", "true")
		}

		for _, code := range opts.ContextCodes {
			query.Add("context_codes[]", code)
		}

		for _, exc := range opts.Excludes {
			query.Add("excludes[]", exc)
		}

		for _, inc := range opts.Includes {
			query.Add("includes[]", inc)
		}

		if opts.ImportantDates {
			query.Add("important_dates", "true")
		}

		if opts.BlackoutDate {
			query.Add("blackout_date", "true")
		}

		if opts.Page > 0 {
			query.Add("page", strconv.Itoa(opts.Page))
		}

		if opts.PerPage > 0 {
			query.Add("per_page", strconv.Itoa(opts.PerPage))
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	var events []CalendarEvent
	if err := s.client.GetAllPages(ctx, path, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// Get retrieves a single calendar event
func (s *CalendarService) Get(ctx context.Context, eventID int64) (*CalendarEvent, error) {
	path := fmt.Sprintf("/api/v1/calendar_events/%d", eventID)

	var event CalendarEvent
	if err := s.client.GetJSON(ctx, path, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

// CreateCalendarEventParams holds parameters for creating a calendar event
type CreateCalendarEventParams struct {
	ContextCode     string // Required: course_123, user_456, etc.
	Title           string
	Description     string
	StartAt         string // ISO 8601
	EndAt           string // ISO 8601
	LocationName    string
	LocationAddress string
	TimeZoneEdited  string
	AllDay          bool
	RRule           string
	BlackoutDate    bool
}

// Create creates a new calendar event
func (s *CalendarService) Create(ctx context.Context, params *CreateCalendarEventParams) (*CalendarEvent, error) {
	path := "/api/v1/calendar_events"

	body := map[string]interface{}{
		"calendar_event": make(map[string]interface{}),
	}

	eventData := body["calendar_event"].(map[string]interface{})
	eventData["context_code"] = params.ContextCode

	if params.Title != "" {
		eventData["title"] = params.Title
	}

	if params.Description != "" {
		eventData["description"] = params.Description
	}

	if params.StartAt != "" {
		eventData["start_at"] = params.StartAt
	}

	if params.EndAt != "" {
		eventData["end_at"] = params.EndAt
	}

	if params.LocationName != "" {
		eventData["location_name"] = params.LocationName
	}

	if params.LocationAddress != "" {
		eventData["location_address"] = params.LocationAddress
	}

	if params.TimeZoneEdited != "" {
		eventData["time_zone_edited"] = params.TimeZoneEdited
	}

	if params.AllDay {
		eventData["all_day"] = true
	}

	if params.RRule != "" {
		eventData["rrule"] = params.RRule
	}

	if params.BlackoutDate {
		eventData["blackout_date"] = true
	}

	var event CalendarEvent
	if err := s.client.PostJSON(ctx, path, body, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

// UpdateCalendarEventParams holds parameters for updating a calendar event
type UpdateCalendarEventParams struct {
	ContextCode     *string
	Title           *string
	Description     *string
	StartAt         *string
	EndAt           *string
	LocationName    *string
	LocationAddress *string
	TimeZoneEdited  *string
	AllDay          *bool
	RRule           *string
	BlackoutDate    *bool
	Which           string // one, all, following (for series)
}

// Update updates an existing calendar event
func (s *CalendarService) Update(ctx context.Context, eventID int64, params *UpdateCalendarEventParams) (*CalendarEvent, error) {
	path := fmt.Sprintf("/api/v1/calendar_events/%d", eventID)

	if params.Which != "" {
		path += "?which=" + params.Which
	}

	body := map[string]interface{}{
		"calendar_event": make(map[string]interface{}),
	}

	eventData := body["calendar_event"].(map[string]interface{})

	if params.ContextCode != nil {
		eventData["context_code"] = *params.ContextCode
	}

	if params.Title != nil {
		eventData["title"] = *params.Title
	}

	if params.Description != nil {
		eventData["description"] = *params.Description
	}

	if params.StartAt != nil {
		eventData["start_at"] = *params.StartAt
	}

	if params.EndAt != nil {
		eventData["end_at"] = *params.EndAt
	}

	if params.LocationName != nil {
		eventData["location_name"] = *params.LocationName
	}

	if params.LocationAddress != nil {
		eventData["location_address"] = *params.LocationAddress
	}

	if params.TimeZoneEdited != nil {
		eventData["time_zone_edited"] = *params.TimeZoneEdited
	}

	if params.AllDay != nil {
		eventData["all_day"] = *params.AllDay
	}

	if params.RRule != nil {
		eventData["rrule"] = *params.RRule
	}

	if params.BlackoutDate != nil {
		eventData["blackout_date"] = *params.BlackoutDate
	}

	var event CalendarEvent
	if err := s.client.PutJSON(ctx, path, body, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

// DeleteOptions holds options for deleting a calendar event
type DeleteOptions struct {
	CancelReason string
	Which        string // one, all, following (for series)
}

// Delete deletes a calendar event
func (s *CalendarService) Delete(ctx context.Context, eventID int64, opts *DeleteOptions) error {
	path := fmt.Sprintf("/api/v1/calendar_events/%d", eventID)

	if opts != nil {
		query := url.Values{}

		if opts.CancelReason != "" {
			query.Add("cancel_reason", opts.CancelReason)
		}

		if opts.Which != "" {
			query.Add("which", opts.Which)
		}

		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	_, err := s.client.Delete(ctx, path)
	return err
}

// Reserve reserves a time slot
func (s *CalendarService) Reserve(ctx context.Context, eventID int64, participantID *int64, comments string, cancelExisting bool) (*CalendarEvent, error) {
	var path string
	if participantID != nil {
		path = fmt.Sprintf("/api/v1/calendar_events/%d/reservations/%d", eventID, *participantID)
	} else {
		path = fmt.Sprintf("/api/v1/calendar_events/%d/reservations", eventID)
	}

	body := make(map[string]interface{})

	if comments != "" {
		body["comments"] = comments
	}

	if cancelExisting {
		body["cancel_existing"] = true
	}

	var event CalendarEvent
	if err := s.client.PostJSON(ctx, path, body, &event); err != nil {
		return nil, err
	}

	return &event, nil
}
