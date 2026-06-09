package options

import "testing"

func TestPeerReviewsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PeerReviewsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PeerReviewsListOptions{CourseID: 1, AssignmentID: 2},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &PeerReviewsListOptions{CourseID: 0, AssignmentID: 2},
			wantErr: true,
		},
		{
			name:    "missing assignment ID",
			opts:    &PeerReviewsListOptions{CourseID: 1, AssignmentID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerReviewsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerReviewsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PeerReviewsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PeerReviewsCreateOptions{CourseID: 1, AssignmentID: 2, SubmissionID: 3, UserID: 4},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &PeerReviewsCreateOptions{CourseID: 0, AssignmentID: 2, SubmissionID: 3, UserID: 4},
			wantErr: true,
		},
		{
			name:    "missing assignment ID",
			opts:    &PeerReviewsCreateOptions{CourseID: 1, AssignmentID: 0, SubmissionID: 3, UserID: 4},
			wantErr: true,
		},
		{
			name:    "missing submission ID",
			opts:    &PeerReviewsCreateOptions{CourseID: 1, AssignmentID: 2, SubmissionID: 0, UserID: 4},
			wantErr: true,
		},
		{
			name:    "missing user ID",
			opts:    &PeerReviewsCreateOptions{CourseID: 1, AssignmentID: 2, SubmissionID: 3, UserID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerReviewsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerReviewsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *PeerReviewsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &PeerReviewsDeleteOptions{CourseID: 1, AssignmentID: 2, SubmissionID: 3, UserID: 4},
			wantErr: false,
		},
		{
			name:    "missing course ID",
			opts:    &PeerReviewsDeleteOptions{CourseID: 0, AssignmentID: 2, SubmissionID: 3, UserID: 4},
			wantErr: true,
		},
		{
			name:    "missing assignment ID",
			opts:    &PeerReviewsDeleteOptions{CourseID: 1, AssignmentID: 0, SubmissionID: 3, UserID: 4},
			wantErr: true,
		},
		{
			name:    "missing submission ID",
			opts:    &PeerReviewsDeleteOptions{CourseID: 1, AssignmentID: 2, SubmissionID: 0, UserID: 4},
			wantErr: true,
		},
		{
			name:    "missing user ID",
			opts:    &PeerReviewsDeleteOptions{CourseID: 1, AssignmentID: 2, SubmissionID: 3, UserID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerReviewsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
