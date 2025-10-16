package todo

import (
	"testing"
)

func TestTranslateError(t *testing.T) {
	tests := []struct {
		name      string
		update    UpdateTodo
		wantError bool
		checkMsg  func(map[string]string) bool
	}{
		{
			name: "title too short",
			update: UpdateTodo{
				Title: strPtr(""),
			},
			wantError: true,
			checkMsg: func(errs map[string]string) bool {
				msg, ok := errs["Title"]
				return ok && msg != "" // Should have a translated message
			},
		},
		{
			name: "title too long",
			update: UpdateTodo{
				Title: strPtr(string(make([]byte, 300))),
			},
			wantError: true,
			checkMsg: func(errs map[string]string) bool {
				msg, ok := errs["Title"]
				return ok && msg != "" // Should have a translated message
			},
		},
		{
			name: "valid update",
			update: UpdateTodo{
				Title: strPtr("Valid Title"),
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.update.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}

			if err != nil {
				// Get translated errors
				translatedErrs := TranslateError(err)
				if translatedErrs == nil {
					t.Error("TranslateError() returned nil for validation error")
				}

				// Print the translated errors for verification
				t.Logf("Translated errors: %v", translatedErrs)

				if tt.checkMsg != nil && !tt.checkMsg(translatedErrs) {
					t.Error("Translated error message validation failed")
				}
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
