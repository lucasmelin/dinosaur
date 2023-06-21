package dns

import (
	"fmt"
	"testing"
)

func TestDNSName_toBytes(t *testing.T) {
	tests := []struct {
		name     string
		dnsName  string
		expected string
	}{
		{
			name:     "basic",
			dnsName:  "google.com",
			expected: "\"\\x06google\\x03com\\x00\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprintf("%q", EncodeName(tt.dnsName))
			if got != tt.expected {
				t.Errorf("expected %q but got %q", tt.expected, got)
			}
		})
	}
}
