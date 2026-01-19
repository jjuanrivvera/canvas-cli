package options

// WebhookListenOptions contains options for starting webhook listener
type WebhookListenOptions struct {
	Addr          string
	Secret        string
	JWKsURL       string
	CanvasDataSvc bool
	Events        []string
	Log           bool
}

// Validate validates the options
func (o *WebhookListenOptions) Validate() error {
	// No required fields - all options are optional with defaults
	return nil
}

// WebhookEventsOptions contains options for listing webhook events
type WebhookEventsOptions struct {
	// No fields needed for listing events
}

// Validate validates the options
func (o *WebhookEventsOptions) Validate() error {
	return nil
}
