package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkGenerateRandomString(b *testing.B) {
	for b.Loop() {
		GenerateRandomString(ALPHA, 6)
	}
}

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		charSet string
		length  int
	}{{
		name:    "charset an length",
		charSet: ALPHA,
		length:  10,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateRandomString(tt.charSet, tt.length)
			for _, char := range got {
				assert.Contains(t, tt.charSet, string(char))
			}
		})
	}
}
