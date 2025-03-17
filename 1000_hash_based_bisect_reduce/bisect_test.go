package main

import (
	"testing"
)

func TestCutMarker(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantShort string
		wantID    uint64
		wantOk    bool
	}{
		{
			name:      "valid hex marker",
			line:      "foo [bisect-match 0x1234] bar",
			wantShort: "foo bar",
			wantID:    0x1234,
			wantOk:    true,
		},
		{
			name:      "valid binary marker",
			line:      "foo [bisect-match 0101] bar",
			wantShort: "foo bar",
			wantID:    0x5,
			wantOk:    true,
		},
		{
			name:      "marker with spaces around",
			line:      "foo [bisect-match 0x1234] bar",
			wantShort: "foo bar",
			wantID:    0x1234,
			wantOk:    true,
		},
		{
			name:      "marker at start of line",
			line:      "[bisect-match 0x1234] bar",
			wantShort: "bar",
			wantID:    0x1234,
			wantOk:    true,
		},
		{
			name:      "marker at end of line",
			line:      "foo [bisect-match 0x1234]",
			wantShort: "foo",
			wantID:    0x1234,
			wantOk:    true,
		},
		{
			name:      "invalid marker missing [",
			line:      "foo bisect-match 0x1234] bar",
			wantShort: "foo bisect-match 0x1234] bar",
			wantID:    0,
			wantOk:    false,
		},
		{
			name:      "invalid marker missing ]",
			line:      "foo [bisect-match 0x1234 bar",
			wantShort: "foo [bisect-match 0x1234 bar",
			wantID:    0,
			wantOk:    false,
		},
		{
			name:      "invalid hex marker too long",
			line:      "foo [bisect-match 0x123456789abcdef1234] bar",
			wantShort: "foo [bisect-match 0x123456789abcdef1234] bar",
			wantID:    0,
			wantOk:    false,
		},
		{
			name:      "invalid binary marker too long",
			line:      "foo [bisect-match 010101010101010101010101010101010101010101010101010101010101010101] bar",
			wantShort: "foo [bisect-match 010101010101010101010101010101010101010101010101010101010101010101] bar",
			wantID:    0,
			wantOk:    false,
		},
		{
			name:      "no marker",
			line:      "foo bar",
			wantShort: "foo bar",
			wantID:    0,
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotShort, gotID, gotOk := CutMarker(tt.line)
			if gotShort != tt.wantShort {
				t.Errorf("CutMarker() gotShort = %v, want %v", gotShort, tt.wantShort)
			}
			if gotID != tt.wantID {
				t.Errorf("CutMarker() gotID = %v, want %v", gotID, tt.wantID)
			}
			if gotOk != tt.wantOk {
				t.Errorf("CutMarker() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
