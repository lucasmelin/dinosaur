package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDNSHeader_toBytes(t *testing.T) {
	tests := []struct {
		name     string
		header   DNSHeader
		expected string
	}{
		{
			name: "basic",
			header: DNSHeader{
				id:             0x1314,
				flags:          0,
				numQuestions:   1,
				numAnswers:     0,
				numAuthorities: 0,
				numAdditionals: 0,
			},
			expected: "\"\\x13\\x14\\x00\\x00\\x00\\x01\\x00\\x00\\x00\\x00\\x00\\x00\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprintf("%q", tt.header.toBytes())
			if got != tt.expected {
				t.Errorf("expected %q but got %q", tt.expected, got)
			}
		})
	}
}

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
			got := fmt.Sprintf("%q", encodeDNSName(tt.dnsName))
			if got != tt.expected {
				t.Errorf("expected %q but got %q", tt.expected, got)
			}
		})
	}
}

func Test_buildQuery(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		domainName string
		recordType int
		want       string
	}{
		{
			name:       "basic",
			id:         0x8298,
			domainName: "example.com",
			recordType: TypeA,
			want:       "\x82\x98\x01\x00\x00\x01\x00\x00\x00\x00\x00\x00\aexample\x03com\x00\x00\x01\x00\x01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildQuery(tt.id, tt.domainName, tt.recordType)
			if !reflect.DeepEqual(got, []byte(tt.want)) {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func Test_parseHeader(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   DNSHeader
	}{
		{
			name:   "basic",
			header: "`V\x81\x80\x00\x01\x00\x01\x00\x00\x00\x00",
			want: DNSHeader{
				id:             24662,
				flags:          33152,
				numQuestions:   1,
				numAnswers:     1,
				numAuthorities: 0,
				numAdditionals: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseHeader([]byte(tt.header)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
