package options

import "testing"

func TestWebhookListenOptions_Validate(t *testing.T) {
	opts := &WebhookListenOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("WebhookListenOptions.Validate() error = %v, want nil", err)
	}
}

func TestWebhookEventsOptions_Validate(t *testing.T) {
	opts := &WebhookEventsOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("WebhookEventsOptions.Validate() error = %v, want nil", err)
	}
}
