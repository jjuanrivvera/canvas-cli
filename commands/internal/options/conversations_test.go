package options

import "testing"

func TestConversationsListOptions_Validate(t *testing.T) {
	opts := &ConversationsListOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("ConversationsListOptions.Validate() error = %v, want nil", err)
	}
}

func TestConversationsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConversationsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConversationsGetOptions{ConversationID: 1},
			wantErr: false,
		},
		{
			name:    "zero conversation ID",
			opts:    &ConversationsGetOptions{ConversationID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConversationsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConversationsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConversationsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConversationsCreateOptions{Recipients: "123", Body: "Hello"},
			wantErr: false,
		},
		{
			name:    "missing recipients",
			opts:    &ConversationsCreateOptions{Recipients: "", Body: "Hello"},
			wantErr: true,
		},
		{
			name:    "missing body",
			opts:    &ConversationsCreateOptions{Recipients: "123", Body: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConversationsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConversationsReplyOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConversationsReplyOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConversationsReplyOptions{ConversationID: 1, Body: "Reply text"},
			wantErr: false,
		},
		{
			name:    "zero conversation ID",
			opts:    &ConversationsReplyOptions{ConversationID: 0, Body: "Reply text"},
			wantErr: true,
		},
		{
			name:    "missing body",
			opts:    &ConversationsReplyOptions{ConversationID: 1, Body: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConversationsReplyOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConversationsAddRecipientsOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConversationsAddRecipientsOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConversationsAddRecipientsOptions{ConversationID: 1, Recipients: "456"},
			wantErr: false,
		},
		{
			name:    "zero conversation ID",
			opts:    &ConversationsAddRecipientsOptions{ConversationID: 0, Recipients: "456"},
			wantErr: true,
		},
		{
			name:    "missing recipients",
			opts:    &ConversationsAddRecipientsOptions{ConversationID: 1, Recipients: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConversationsAddRecipientsOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConversationsArchiveOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConversationsArchiveOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConversationsArchiveOptions{ConversationID: 1},
			wantErr: false,
		},
		{
			name:    "zero conversation ID",
			opts:    &ConversationsArchiveOptions{ConversationID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConversationsArchiveOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConversationsUnarchiveOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConversationsUnarchiveOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConversationsUnarchiveOptions{ConversationID: 1},
			wantErr: false,
		},
		{
			name:    "zero conversation ID",
			opts:    &ConversationsUnarchiveOptions{ConversationID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConversationsUnarchiveOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConversationsStarOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConversationsStarOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConversationsStarOptions{ConversationID: 1},
			wantErr: false,
		},
		{
			name:    "zero conversation ID",
			opts:    &ConversationsStarOptions{ConversationID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConversationsStarOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConversationsUnstarOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConversationsUnstarOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConversationsUnstarOptions{ConversationID: 1},
			wantErr: false,
		},
		{
			name:    "zero conversation ID",
			opts:    &ConversationsUnstarOptions{ConversationID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConversationsUnstarOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConversationsMarkReadOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConversationsMarkReadOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConversationsMarkReadOptions{ConversationID: 1},
			wantErr: false,
		},
		{
			name:    "zero conversation ID",
			opts:    &ConversationsMarkReadOptions{ConversationID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConversationsMarkReadOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConversationsMarkAllReadOptions_Validate(t *testing.T) {
	opts := &ConversationsMarkAllReadOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("ConversationsMarkAllReadOptions.Validate() error = %v, want nil", err)
	}
}

func TestConversationsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ConversationsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ConversationsDeleteOptions{ConversationID: 1},
			wantErr: false,
		},
		{
			name:    "zero conversation ID",
			opts:    &ConversationsDeleteOptions{ConversationID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConversationsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConversationsUnreadCountOptions_Validate(t *testing.T) {
	opts := &ConversationsUnreadCountOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("ConversationsUnreadCountOptions.Validate() error = %v, want nil", err)
	}
}
