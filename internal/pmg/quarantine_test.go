package pmg

import (
	"testing"
)

func TestFilterQuarantine(t *testing.T) {
	mails := []QuarantineMail{
		{ID: "1", From: "spam@evil.com", Receiver: "user@example.com", Subject: "Buy now"},
		{ID: "2", From: "info@example.org", Receiver: "admin@examplecorp.com", Subject: "Invoice"},
		{ID: "3", From: "info@example.org", Receiver: "user@example.com", Subject: "Hello"},
	}

	tests := []struct {
		name     string
		sender   string
		receiver string
		want     int
	}{
		{"no filter", "*", "*", 3},
		{"sender filter", "example.org", "*", 2},
		{"receiver filter", "*", "examplecorp.com", 1},
		{"both filters", "example.org", "examplecorp.com", 1},
		{"no match", "noexist.xyz", "*", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterQuarantine(mails, tt.sender, tt.receiver)
			if len(result) != tt.want {
				t.Errorf("FilterQuarantine(%q, %q) = %d mails, want %d", tt.sender, tt.receiver, len(result), tt.want)
			}
		})
	}
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		s       string
		pattern string
		want    bool
	}{
		{"spam@evil.com", "*", true},
		{"spam@evil.com", "", true},
		{"spam@evil.com", "evil.com", true},
		{"spam@evil.com", "EVIL.COM", true},
		{"spam@evil.com", "spam", true},
		{"spam@evil.com", "good.com", false},
		{"spam@evil.com", "*@evil.com", true},
		{"info@example.org", "*.org", true},
		{"hello.world", "hello.*", true},
	}

	for _, tt := range tests {
		got := matchPattern(tt.s, tt.pattern)
		if got != tt.want {
			t.Errorf("matchPattern(%q, %q) = %v, want %v", tt.s, tt.pattern, got, tt.want)
		}
	}
}