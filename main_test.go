package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/lucasmelin/dinosaur/dns"
)

func TestDNSHeader_toBytes(t *testing.T) {
	tests := []struct {
		name     string
		header   dns.Header
		expected string
	}{
		{
			name: "basic",
			header: dns.Header{
				Id:             0x1314,
				Flags:          0,
				NumQuestions:   1,
				NumAnswers:     0,
				NumAuthorities: 0,
				NumAdditionals: 0,
			},
			expected: "\"\\x13\\x14\\x00\\x00\\x00\\x01\\x00\\x00\\x00\\x00\\x00\\x00\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprintf("%q", tt.header.ToBytes())
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
			got := fmt.Sprintf("%q", dns.EncodeDNSName(tt.dnsName))
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
		recordType string
		want       string
	}{
		{
			name:       "basic",
			id:         0x8298,
			domainName: "example.com",
			recordType: "A",
			want:       "\x82\x98\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\aexample\x03com\x00\x00\x01\x00\x01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dns.BuildQuery(tt.id, tt.domainName, tt.recordType)
			if !reflect.DeepEqual(got, []byte(tt.want)) {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func Test_parseHeader(t *testing.T) {
	tests := []struct {
		name   string
		header *bytes.Reader
		want   dns.Header
	}{
		{
			name:   "basic",
			header: bytes.NewReader([]byte("`V\x81\x80\x00\x01\x00\x01\x00\x00\x00\x00")),
			want: dns.Header{
				Id:             24662,
				Flags:          33152,
				NumQuestions:   1,
				NumAnswers:     1,
				NumAuthorities: 0,
				NumAdditionals: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dns.ParseHeader(tt.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
