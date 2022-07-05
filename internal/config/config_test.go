package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		expected Interface
	}{
		{
			name:     "creates mew config",
			expected: Config{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig()
			assert.NoError(t, err)
			assert.IsType(t, tt.expected, got)
		})
	}
}
