package dns

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestDNSHeader_toBytes(t *testing.T) {
	tests := []struct {
		name     string
		header   Header
		expected string
	}{
		{
			name: "basic",
			header: Header{
				ID:             0x1314,
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

func Test_parseHeader(t *testing.T) {
	tests := []struct {
		name   string
		header *bytes.Reader
		want   Header
	}{
		{
			name:   "basic",
			header: bytes.NewReader([]byte("`V\x81\x80\x00\x01\x00\x01\x00\x00\x00\x00")),
			want: Header{
				ID:             24662,
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
			if got := ParseHeader(tt.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestHeader_String(t *testing.T) {
	type fields struct {
		ID             uint16
		Flags          uint16
		NumQuestions   uint16
		NumAnswers     uint16
		NumAuthorities uint16
		NumAdditionals uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "empty header",
			fields: fields{},
			want: `Header{
    ID: 0,
    Flags: 0,
    NumQuestions: 0,
    NumAnswers: 0,
    NumAuthorities: 0,
    NumAdditionals: 0
  }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Header{
				ID:             tt.fields.ID,
				Flags:          tt.fields.Flags,
				NumQuestions:   tt.fields.NumQuestions,
				NumAnswers:     tt.fields.NumAnswers,
				NumAuthorities: tt.fields.NumAuthorities,
				NumAdditionals: tt.fields.NumAdditionals,
			}
			if got := h.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
