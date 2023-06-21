package dns

import (
	"reflect"
	"testing"
)

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
			got := BuildQuery(tt.id, tt.domainName, tt.recordType)
			if !reflect.DeepEqual(got, []byte(tt.want)) {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}
