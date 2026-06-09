package options

import "testing"

func TestDiscussionsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *DiscussionsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &DiscussionsListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &DiscussionsListOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscussionsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscussionsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *DiscussionsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &DiscussionsGetOptions{CourseID: 1, TopicID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &DiscussionsGetOptions{CourseID: 0, TopicID: 2},
			wantErr: true,
		},
		{
			name:    "zero topic ID",
			opts:    &DiscussionsGetOptions{CourseID: 1, TopicID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscussionsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscussionsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *DiscussionsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &DiscussionsCreateOptions{CourseID: 1, Title: "Test Discussion"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &DiscussionsCreateOptions{CourseID: 0, Title: "Test Discussion"},
			wantErr: true,
		},
		{
			name:    "missing title",
			opts:    &DiscussionsCreateOptions{CourseID: 1, Title: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscussionsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscussionsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *DiscussionsUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &DiscussionsUpdateOptions{CourseID: 1, TopicID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &DiscussionsUpdateOptions{CourseID: 0, TopicID: 2},
			wantErr: true,
		},
		{
			name:    "zero topic ID",
			opts:    &DiscussionsUpdateOptions{CourseID: 1, TopicID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscussionsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscussionsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *DiscussionsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &DiscussionsDeleteOptions{CourseID: 1, TopicID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &DiscussionsDeleteOptions{CourseID: 0, TopicID: 2},
			wantErr: true,
		},
		{
			name:    "zero topic ID",
			opts:    &DiscussionsDeleteOptions{CourseID: 1, TopicID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscussionsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscussionsEntriesOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *DiscussionsEntriesOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &DiscussionsEntriesOptions{CourseID: 1, TopicID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &DiscussionsEntriesOptions{CourseID: 0, TopicID: 2},
			wantErr: true,
		},
		{
			name:    "zero topic ID",
			opts:    &DiscussionsEntriesOptions{CourseID: 1, TopicID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscussionsEntriesOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscussionsPostOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *DiscussionsPostOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &DiscussionsPostOptions{CourseID: 1, TopicID: 2, Message: "Hello"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &DiscussionsPostOptions{CourseID: 0, TopicID: 2, Message: "Hello"},
			wantErr: true,
		},
		{
			name:    "zero topic ID",
			opts:    &DiscussionsPostOptions{CourseID: 1, TopicID: 0, Message: "Hello"},
			wantErr: true,
		},
		{
			name:    "missing message",
			opts:    &DiscussionsPostOptions{CourseID: 1, TopicID: 2, Message: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscussionsPostOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscussionsReplyOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *DiscussionsReplyOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &DiscussionsReplyOptions{CourseID: 1, TopicID: 2, EntryID: 3, Message: "Reply"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &DiscussionsReplyOptions{CourseID: 0, TopicID: 2, EntryID: 3, Message: "Reply"},
			wantErr: true,
		},
		{
			name:    "zero topic ID",
			opts:    &DiscussionsReplyOptions{CourseID: 1, TopicID: 0, EntryID: 3, Message: "Reply"},
			wantErr: true,
		},
		{
			name:    "zero entry ID",
			opts:    &DiscussionsReplyOptions{CourseID: 1, TopicID: 2, EntryID: 0, Message: "Reply"},
			wantErr: true,
		},
		{
			name:    "missing message",
			opts:    &DiscussionsReplyOptions{CourseID: 1, TopicID: 2, EntryID: 3, Message: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscussionsReplyOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscussionsSubscribeOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *DiscussionsSubscribeOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &DiscussionsSubscribeOptions{CourseID: 1, TopicID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &DiscussionsSubscribeOptions{CourseID: 0, TopicID: 2},
			wantErr: true,
		},
		{
			name:    "zero topic ID",
			opts:    &DiscussionsSubscribeOptions{CourseID: 1, TopicID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscussionsSubscribeOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscussionsUnsubscribeOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *DiscussionsUnsubscribeOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &DiscussionsUnsubscribeOptions{CourseID: 1, TopicID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &DiscussionsUnsubscribeOptions{CourseID: 0, TopicID: 2},
			wantErr: true,
		},
		{
			name:    "zero topic ID",
			opts:    &DiscussionsUnsubscribeOptions{CourseID: 1, TopicID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscussionsUnsubscribeOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
