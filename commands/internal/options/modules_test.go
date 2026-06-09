package options

import "testing"

func TestModulesListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ModulesListOptions{CourseID: 1},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesListOptions{CourseID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ModulesGetOptions{CourseID: 1, ModuleID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesGetOptions{CourseID: 0, ModuleID: 2},
			wantErr: true,
		},
		{
			name:    "zero module ID",
			opts:    &ModulesGetOptions{CourseID: 1, ModuleID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesCreateOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ModulesCreateOptions{CourseID: 1, Name: "Week 1"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesCreateOptions{CourseID: 0, Name: "Week 1"},
			wantErr: true,
		},
		{
			name:    "missing name",
			opts:    &ModulesCreateOptions{CourseID: 1, Name: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid with one field set",
			opts:    &ModulesUpdateOptions{CourseID: 1, ModuleID: 2, NameSet: true, Name: "New Name"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesUpdateOptions{CourseID: 0, ModuleID: 2, NameSet: true},
			wantErr: true,
		},
		{
			name:    "zero module ID",
			opts:    &ModulesUpdateOptions{CourseID: 1, ModuleID: 0, NameSet: true},
			wantErr: true,
		},
		{
			name:    "no fields set",
			opts:    &ModulesUpdateOptions{CourseID: 1, ModuleID: 2},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ModulesDeleteOptions{CourseID: 1, ModuleID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesDeleteOptions{CourseID: 0, ModuleID: 2},
			wantErr: true,
		},
		{
			name:    "zero module ID",
			opts:    &ModulesDeleteOptions{CourseID: 1, ModuleID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesRelockOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesRelockOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ModulesRelockOptions{CourseID: 1, ModuleID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesRelockOptions{CourseID: 0, ModuleID: 2},
			wantErr: true,
		},
		{
			name:    "zero module ID",
			opts:    &ModulesRelockOptions{CourseID: 1, ModuleID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesRelockOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesPublishOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesPublishOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ModulesPublishOptions{CourseID: 1, ModuleID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesPublishOptions{CourseID: 0, ModuleID: 2},
			wantErr: true,
		},
		{
			name:    "zero module ID",
			opts:    &ModulesPublishOptions{CourseID: 1, ModuleID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesPublishOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesUnpublishOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesUnpublishOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ModulesUnpublishOptions{CourseID: 1, ModuleID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesUnpublishOptions{CourseID: 0, ModuleID: 2},
			wantErr: true,
		},
		{
			name:    "zero module ID",
			opts:    &ModulesUnpublishOptions{CourseID: 1, ModuleID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesUnpublishOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesItemsListOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesItemsListOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ModulesItemsListOptions{CourseID: 1, ModuleID: 2},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesItemsListOptions{CourseID: 0, ModuleID: 2},
			wantErr: true,
		},
		{
			name:    "zero module ID",
			opts:    &ModulesItemsListOptions{CourseID: 1, ModuleID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesItemsListOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesItemsGetOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesItemsGetOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ModulesItemsGetOptions{CourseID: 1, ModuleID: 2, ItemID: 3},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesItemsGetOptions{CourseID: 0, ModuleID: 2, ItemID: 3},
			wantErr: true,
		},
		{
			name:    "zero module ID",
			opts:    &ModulesItemsGetOptions{CourseID: 1, ModuleID: 0, ItemID: 3},
			wantErr: true,
		},
		{
			name:    "zero item ID",
			opts:    &ModulesItemsGetOptions{CourseID: 1, ModuleID: 2, ItemID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesItemsGetOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesItemsCreateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesItemsCreateOptions
		wantErr bool
	}{
		{
			name:    "valid Page type",
			opts:    &ModulesItemsCreateOptions{CourseID: 1, ModuleID: 2, Type: "Page", Title: "My Page"},
			wantErr: false,
		},
		{
			name:    "valid Assignment type",
			opts:    &ModulesItemsCreateOptions{CourseID: 1, ModuleID: 2, Type: "Assignment", Title: "HW1"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesItemsCreateOptions{CourseID: 0, ModuleID: 2, Type: "Page", Title: "My Page"},
			wantErr: true,
		},
		{
			name:    "zero module ID",
			opts:    &ModulesItemsCreateOptions{CourseID: 1, ModuleID: 0, Type: "Page", Title: "My Page"},
			wantErr: true,
		},
		{
			name:    "missing type",
			opts:    &ModulesItemsCreateOptions{CourseID: 1, ModuleID: 2, Type: "", Title: "My Page"},
			wantErr: true,
		},
		{
			name:    "missing title",
			opts:    &ModulesItemsCreateOptions{CourseID: 1, ModuleID: 2, Type: "Page", Title: ""},
			wantErr: true,
		},
		{
			name:    "invalid type",
			opts:    &ModulesItemsCreateOptions{CourseID: 1, ModuleID: 2, Type: "Invalid", Title: "My Page"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesItemsCreateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesItemsUpdateOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesItemsUpdateOptions
		wantErr bool
	}{
		{
			name:    "valid with one field set",
			opts:    &ModulesItemsUpdateOptions{CourseID: 1, ModuleID: 2, ItemID: 3, TitleSet: true, Title: "New Title"},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesItemsUpdateOptions{CourseID: 0, ModuleID: 2, ItemID: 3, TitleSet: true},
			wantErr: true,
		},
		{
			name:    "zero module ID",
			opts:    &ModulesItemsUpdateOptions{CourseID: 1, ModuleID: 0, ItemID: 3, TitleSet: true},
			wantErr: true,
		},
		{
			name:    "zero item ID",
			opts:    &ModulesItemsUpdateOptions{CourseID: 1, ModuleID: 2, ItemID: 0, TitleSet: true},
			wantErr: true,
		},
		{
			name:    "no fields set",
			opts:    &ModulesItemsUpdateOptions{CourseID: 1, ModuleID: 2, ItemID: 3},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesItemsUpdateOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesItemsDeleteOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesItemsDeleteOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ModulesItemsDeleteOptions{CourseID: 1, ModuleID: 2, ItemID: 3},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesItemsDeleteOptions{CourseID: 0, ModuleID: 2, ItemID: 3},
			wantErr: true,
		},
		{
			name:    "zero module ID",
			opts:    &ModulesItemsDeleteOptions{CourseID: 1, ModuleID: 0, ItemID: 3},
			wantErr: true,
		},
		{
			name:    "zero item ID",
			opts:    &ModulesItemsDeleteOptions{CourseID: 1, ModuleID: 2, ItemID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesItemsDeleteOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModulesItemsDoneOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ModulesItemsDoneOptions
		wantErr bool
	}{
		{
			name:    "valid",
			opts:    &ModulesItemsDoneOptions{CourseID: 1, ModuleID: 2, ItemID: 3},
			wantErr: false,
		},
		{
			name:    "zero course ID",
			opts:    &ModulesItemsDoneOptions{CourseID: 0, ModuleID: 2, ItemID: 3},
			wantErr: true,
		},
		{
			name:    "zero module ID",
			opts:    &ModulesItemsDoneOptions{CourseID: 1, ModuleID: 0, ItemID: 3},
			wantErr: true,
		},
		{
			name:    "zero item ID",
			opts:    &ModulesItemsDoneOptions{CourseID: 1, ModuleID: 2, ItemID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ModulesItemsDoneOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
