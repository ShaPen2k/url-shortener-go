package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 20",
			size: 20,
		},
		{
			name: "size = 30",
			size: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate two strings to check for randomness
			str1 := NewRandomString(tt.size)
			str2 := NewRandomString(tt.size)

			// Assert string length
			assert.Len(t, str1, tt.size)

			// Validate non-empty result for positive sizes
			if tt.size > 0 {
				assert.NotEmpty(t, str1)
			}

			// Assert uniqueness to ensure random output
			// Skip for very small sizes due to high collision probability
			if tt.size > 5 {
				assert.NotEqual(t, str1, str2)
			}
		})
	}
}
