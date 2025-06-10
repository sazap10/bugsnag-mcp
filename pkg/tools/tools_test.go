package tools

import (
	"testing"
)

func TestGetEventIDFromIDOrLink(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantID    string
		wantError bool
	}{
		{
			name:      "plain event id",
			input:     "abc123",
			wantID:    "abc123",
			wantError: false,
		},
		{
			name:      "event link with event_id param",
			input:     "https://app.bugsnag.com/org/proj/errors/errid/events/event?event_id=xyz789",
			wantID:    "xyz789",
			wantError: false,
		},
		{
			name:      "event link without event_id param",
			input:     "https://app.bugsnag.com/org/proj/errors/errid/events/event",
			wantID:    "",
			wantError: true,
		},
		{
			name:      "malformed url",
			input:     "http://%41:8080/", // invalid URL
			wantID:    "",
			wantError: true,
		},
		{
			name:      "slash but not http",
			input:     "foo/bar",
			wantID:    "",
			wantError: true,
		},
		{
			name:      "http but no query",
			input:     "http://example.com/foo/bar",
			wantID:    "",
			wantError: true,
		},
		{
			name:      "event_id param but empty",
			input:     "https://app.bugsnag.com/org/proj/errors/errid/events/event?event_id=",
			wantID:    "",
			wantError: true,
		},
		{
			name:      "event_id param with multiple values",
			input:     "https://app.bugsnag.com/org/proj/errors/errid/events/event?event_id=one&event_id=two",
			wantID:    "one",
			wantError: false,
		},
		// On-premise and custom domain tests
		{
			name:      "on-premise domain with event_id",
			input:     "https://bugsnag.mycompany.com/org/proj/errors/errid/events/event?event_id=onprem123",
			wantID:    "onprem123",
			wantError: false,
		},
		{
			name:      "on-premise host:port with event_id",
			input:     "http://localhost:8080/org/proj/errors/errid/events/event?event_id=local456",
			wantID:    "local456",
			wantError: false,
		},
		{
			name:      "on-premise host:port without event_id",
			input:     "http://localhost:8080/org/proj/errors/errid/events/event",
			wantID:    "",
			wantError: true,
		},
		{
			name:      "custom domain with event_id",
			input:     "https://errors.example.org/org/proj/errors/errid/events/event?event_id=custom789",
			wantID:    "custom789",
			wantError: false,
		},
		{
			name:      "custom domain without event_id",
			input:     "https://errors.example.org/org/proj/errors/errid/events/event",
			wantID:    "",
			wantError: true,
		},
		{
			name:      "custom domain with empty event_id",
			input:     "https://errors.example.org/org/proj/errors/errid/events/event?event_id=",
			wantID:    "",
			wantError: true,
		},
		{
			name:      "custom domain with multiple event_id values",
			input:     "https://errors.example.org/org/proj/errors/errid/events/event?event_id=first&event_id=second",
			wantID:    "first",
			wantError: false,
		},
		// Edge cases
		{
			name:      "event_id param with special characters",
			input:     "https://app.bugsnag.com/org/proj/errors/errid/events/event?event_id=ev%20ent%2Fid",
			wantID:    "ev ent/id",
			wantError: false,
		},
		{
			name:      "event_id param with unicode",
			input:     "https://app.bugsnag.com/org/proj/errors/errid/events/event?event_id=%E2%9C%93",
			wantID:    "✓",
			wantError: false,
		},
		{
			name:      "empty string",
			input:     "",
			wantID:    "",
			wantError: false,
		},
		{
			name:      "just a slash",
			input:     "/",
			wantID:    "",
			wantError: true,
		},
		{
			name:      "http only",
			input:     "http://",
			wantID:    "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := getEventIDFromIDOrLink(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("getEventIDFromIDOrLink() error = %v, wantError %v", err, tt.wantError)
			}
			if gotID != tt.wantID {
				t.Errorf("getEventIDFromIDOrLink() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}
