package resources

import (
	"reflect"
	"testing"
)

func TestExtractIDsFromURI(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		segments []string
		want     map[string]string
		wantErr  bool
	}{
		{
			name:     "single segment found",
			uri:      "bugsnag://projects/123",
			segments: []string{"projects"},
			want:     map[string]string{"projects": "123"},
			wantErr:  false,
		},
		{
			name:     "multiple segments found",
			uri:      "bugsnag://projects/123/events/456",
			segments: []string{"projects", "events"},
			want:     map[string]string{"projects": "123", "events": "456"},
			wantErr:  false,
		},
		{
			name:     "segment not found",
			uri:      "bugsnag://organizations/789",
			segments: []string{"projects"},
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "empty uri",
			uri:      "",
			segments: []string{"projects"},
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "segment at end with no id",
			uri:      "bugsnag://projects",
			segments: []string{"projects"},
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "extra segments in uri",
			uri:      "bugsnag://organizations/789/projects/123/events/456",
			segments: []string{"organizations", "projects", "events"},
			want:     map[string]string{"organizations": "789", "projects": "123", "events": "456"},
			wantErr:  false,
		},
		{
			name:     "segment appears twice, only first is used",
			uri:      "bugsnag://projects/123/projects/456",
			segments: []string{"projects"},
			want:     map[string]string{"projects": "456"},
			wantErr:  false,
		},
		{
			name:     "segment present but no following id",
			uri:      "bugsnag://projects/",
			segments: []string{"projects"},
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "segment with trailing slash",
			uri:      "bugsnag://projects/123/",
			segments: []string{"projects"},
			want:     map[string]string{"projects": "123"},
			wantErr:  false,
		},
		{
			name:     "multiple segments, one missing id",
			uri:      "bugsnag://projects/123/events",
			segments: []string{"projects", "events"},
			want:     map[string]string{"projects": "123"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractIDsFromURI(tt.uri, tt.segments...)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractIDsFromURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractIDsFromURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
