package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckLuhn(t *testing.T) {
	tests := []struct {
		name string
		cnn  string
		want bool
	}{
		{
			name: "Correct card",
			cnn:  "4539148803436467",
			want: true,
		},
		{
			name: "Incorrect card",
			cnn:  "4485275702468698",
			want: false,
		},
		{
			name: "Empty card",
			cnn:  "",
			want: false,
		},
		{
			name: "Word card",
			cnn:  "sdfsdfsdfdsf",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckLuhn(tt.cnn)
			assert.Equal(t, tt.want, got)
		})
	}
}
