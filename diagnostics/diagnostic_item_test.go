package diagnostics

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiagnosticLevel_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		level    DiagnosticLevel
		expected string
	}{
		{
			name:     "Info",
			level:    Info,
			expected: "\"Info\"",
		},
		{
			name:     "Warning",
			level:    Warning,
			expected: "\"Warning\"",
		},
		{
			name:     "Error",
			level:    Error,
			expected: "\"Error\"",
		},
		{
			name:     "Trace",
			level:    Trace,
			expected: "\"Trace\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.level)
			assert.NoError(t, err, "Unexpected error while marshaling")
			assert.Equal(t, []byte(tt.expected), b, "Unexpected marshaled JSON")
		})
	}
}

func TestDiagnosticLevel_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected DiagnosticLevel
	}{
		{
			name:     "Info",
			input:    "\"Info\"",
			expected: Info,
		},
		{
			name:     "Warning",
			input:    "\"Warning\"",
			expected: Warning,
		},
		{
			name:     "Error",
			input:    "\"Error\"",
			expected: Error,
		},
		{
			name:     "Trace",
			input:    "\"Trace\"",
			expected: Trace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var level DiagnosticLevel
			err := json.Unmarshal([]byte(tt.input), &level)
			assert.NoError(t, err, "Unexpected error while unmarshaling")
			assert.Equal(t, tt.expected, level, "Unexpected unmarshaled value")
		})
	}
}
