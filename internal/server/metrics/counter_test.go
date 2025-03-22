//go:build server
// +build server

package metrics_test

import (
	"github.com/mobypolo/ya-41go/internal/server/metrics"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterProcessor_ValidateName(t *testing.T) {
	p := metrics.NewCounterProcessor()

	assert.NoError(t, p.ValidateName("requests"))
	assert.Error(t, p.ValidateName("unknown"))
}

func TestCounterProcessor_ParseValue(t *testing.T) {
	p := metrics.NewCounterProcessor()

	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{
			name:     "valid integer",
			input:    "123",
			expected: 123,
			wantErr:  false,
		},
		{
			name:    "invalid string",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:     "negative integer",
			input:    "-42",
			expected: -42,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := p.ParseValue(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, val)
			}
		})
	}
}
