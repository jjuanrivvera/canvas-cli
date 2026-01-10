package api

import (
	"testing"
	"time"
)

func TestNormalizeCourse(t *testing.T) {
	tests := []struct {
		name  string
		input *Course
		want  *Course
	}{
		{
			name:  "nil course",
			input: nil,
			want:  nil,
		},
		{
			name: "course with nil permissions",
			input: &Course{
				ID:   1,
				Name: "Test Course",
			},
			want: &Course{
				ID:                                1,
				Name:                              "Test Course",
				Permissions:                       make(map[string]bool),
				BlueprintRestrictions:             make(map[string]bool),
				BlueprintRestrictionsByObjectType: make(map[string]map[string]bool),
			},
		},
		{
			name: "course with existing permissions",
			input: &Course{
				ID:          1,
				Name:        "Test Course",
				Permissions: map[string]bool{"read": true},
			},
			want: &Course{
				ID:                                1,
				Name:                              "Test Course",
				Permissions:                       map[string]bool{"read": true},
				BlueprintRestrictions:             make(map[string]bool),
				BlueprintRestrictionsByObjectType: make(map[string]map[string]bool),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeCourse(tt.input)
			if got == nil && tt.want == nil {
				return
			}
			if got == nil || tt.want == nil {
				t.Errorf("NormalizeCourse() = %v, want %v", got, tt.want)
				return
			}
			if got.Permissions == nil {
				t.Error("NormalizeCourse() Permissions should not be nil")
			}
			if got.BlueprintRestrictions == nil {
				t.Error("NormalizeCourse() BlueprintRestrictions should not be nil")
			}
			if got.BlueprintRestrictionsByObjectType == nil {
				t.Error("NormalizeCourse() BlueprintRestrictionsByObjectType should not be nil")
			}
		})
	}
}

func TestNormalizeUser(t *testing.T) {
	tests := []struct {
		name  string
		input *User
		want  *User
	}{
		{
			name:  "nil user",
			input: nil,
			want:  nil,
		},
		{
			name: "user with nil enrollments",
			input: &User{
				ID:   1,
				Name: "Test User",
			},
			want: &User{
				ID:          1,
				Name:        "Test User",
				Enrollments: []Enrollment{},
			},
		},
		{
			name: "user with existing enrollments",
			input: &User{
				ID:          1,
				Name:        "Test User",
				Enrollments: []Enrollment{{ID: 1}},
			},
			want: &User{
				ID:          1,
				Name:        "Test User",
				Enrollments: []Enrollment{{ID: 1}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeUser(tt.input)
			if got == nil && tt.want == nil {
				return
			}
			if got == nil || tt.want == nil {
				t.Errorf("NormalizeUser() = %v, want %v", got, tt.want)
				return
			}
			if got.Enrollments == nil {
				t.Error("NormalizeUser() Enrollments should not be nil")
			}
		})
	}
}

func TestNormalizeAssignment(t *testing.T) {
	tests := []struct {
		name  string
		input *Assignment
		want  *Assignment
	}{
		{
			name:  "nil assignment",
			input: nil,
			want:  nil,
		},
		{
			name: "assignment with nil slices",
			input: &Assignment{
				ID:   1,
				Name: "Test Assignment",
			},
			want: &Assignment{
				ID:                         1,
				Name:                       "Test Assignment",
				AllowedExtensions:          []string{},
				SubmissionTypes:            []string{},
				NeedsGradingCountBySection: []SectionGradingCount{},
				FrozenAttributes:           []string{},
				Rubric:                     []RubricCriterion{},
				AssignmentVisibility:       []int64{},
				Overrides:                  []AssignmentOverride{},
				TurnitinSettings:           make(map[string]interface{}),
				ExternalToolTagAttributes:  make(map[string]interface{}),
				IntegrationData:            make(map[string]interface{}),
				RubricSettings:             make(map[string]interface{}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeAssignment(tt.input)
			if got == nil && tt.want == nil {
				return
			}
			if got == nil || tt.want == nil {
				t.Errorf("NormalizeAssignment() = %v, want %v", got, tt.want)
				return
			}
			if got.AllowedExtensions == nil {
				t.Error("NormalizeAssignment() AllowedExtensions should not be nil")
			}
			if got.SubmissionTypes == nil {
				t.Error("NormalizeAssignment() SubmissionTypes should not be nil")
			}
		})
	}
}

func TestNormalizeSubmission(t *testing.T) {
	tests := []struct {
		name  string
		input *Submission
		want  *Submission
	}{
		{
			name:  "nil submission",
			input: nil,
			want:  nil,
		},
		{
			name: "submission with nil slices",
			input: &Submission{
				ID:     1,
				UserID: 1,
			},
			want: &Submission{
				ID:                 1,
				UserID:             1,
				Attachments:        []Attachment{},
				SubmissionComments: []SubmissionComment{},
				Rubric:             []RubricAssessment{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeSubmission(tt.input)
			if got == nil && tt.want == nil {
				return
			}
			if got == nil || tt.want == nil {
				t.Errorf("NormalizeSubmission() = %v, want %v", got, tt.want)
				return
			}
			if got.Attachments == nil {
				t.Error("NormalizeSubmission() Attachments should not be nil")
			}
			if got.SubmissionComments == nil {
				t.Error("NormalizeSubmission() SubmissionComments should not be nil")
			}
		})
	}
}

func TestNormalizeEnrollment(t *testing.T) {
	tests := []struct {
		name  string
		input *Enrollment
		want  *Enrollment
	}{
		{
			name:  "nil enrollment",
			input: nil,
			want:  nil,
		},
		{
			name: "enrollment with nil grades",
			input: &Enrollment{
				ID:     1,
				UserID: 1,
			},
			want: &Enrollment{
				ID:     1,
				UserID: 1,
				Grades: &Grades{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeEnrollment(tt.input)
			if got == nil && tt.want == nil {
				return
			}
			if got == nil || tt.want == nil {
				t.Errorf("NormalizeEnrollment() = %v, want %v", got, tt.want)
				return
			}
			if got.Grades == nil {
				t.Error("NormalizeEnrollment() Grades should not be nil")
			}
		})
	}
}

func TestNormalizeCourses(t *testing.T) {
	tests := []struct {
		name  string
		input []Course
		want  int
	}{
		{
			name:  "nil courses",
			input: nil,
			want:  0,
		},
		{
			name:  "empty courses",
			input: []Course{},
			want:  0,
		},
		{
			name: "multiple courses",
			input: []Course{
				{ID: 1, Name: "Course 1"},
				{ID: 2, Name: "Course 2"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeCourses(tt.input)
			if len(got) != tt.want {
				t.Errorf("NormalizeCourses() returned %d courses, want %d", len(got), tt.want)
			}
			for _, course := range got {
				if course.Permissions == nil {
					t.Error("NormalizeCourses() course Permissions should not be nil")
				}
			}
		})
	}
}

func TestNormalizeUsers(t *testing.T) {
	tests := []struct {
		name  string
		input []User
		want  int
	}{
		{
			name:  "nil users",
			input: nil,
			want:  0,
		},
		{
			name:  "empty users",
			input: []User{},
			want:  0,
		},
		{
			name: "multiple users",
			input: []User{
				{ID: 1, Name: "User 1"},
				{ID: 2, Name: "User 2"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeUsers(tt.input)
			if len(got) != tt.want {
				t.Errorf("NormalizeUsers() returned %d users, want %d", len(got), tt.want)
			}
			for _, user := range got {
				if user.Enrollments == nil {
					t.Error("NormalizeUsers() user Enrollments should not be nil")
				}
			}
		})
	}
}

func TestNormalizeTerm(t *testing.T) {
	tests := []struct {
		name string
		term *Term
		want *Term
	}{
		{
			name: "nil term",
			term: nil,
			want: nil,
		},
		{
			name: "term with nil overrides",
			term: &Term{
				ID:        1,
				Name:      "Fall 2024",
				Overrides: nil,
			},
			want: &Term{
				ID:        1,
				Name:      "Fall 2024",
				Overrides: make(map[string]TermOverride),
			},
		},
		{
			name: "term with existing overrides",
			term: &Term{
				ID:   1,
				Name: "Fall 2024",
				Overrides: map[string]TermOverride{
					"key": {StartAt: time.Now(), EndAt: time.Now().Add(time.Hour)},
				},
			},
			want: &Term{
				ID:   1,
				Name: "Fall 2024",
				Overrides: map[string]TermOverride{
					"key": {StartAt: time.Now(), EndAt: time.Now().Add(time.Hour)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeTerm(tt.term)
			if tt.want == nil {
				if got != nil {
					t.Errorf("NormalizeTerm() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("NormalizeTerm() returned nil, want non-nil")
			}

			if got.ID != tt.want.ID {
				t.Errorf("NormalizeTerm() ID = %d, want %d", got.ID, tt.want.ID)
			}

			if got.Name != tt.want.Name {
				t.Errorf("NormalizeTerm() Name = %s, want %s", got.Name, tt.want.Name)
			}

			if got.Overrides == nil {
				t.Error("NormalizeTerm() Overrides should not be nil")
			}

			if len(got.Overrides) != len(tt.want.Overrides) {
				t.Errorf("NormalizeTerm() Overrides length = %d, want %d", len(got.Overrides), len(tt.want.Overrides))
			}
		})
	}
}
