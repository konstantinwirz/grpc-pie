package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// validEndpoint returns a valid endpoint
func validEndpoint() Endpoint {
	endpoint, err := NewEndpoint("localhost:30000/Service/Method")
	if err != nil {
		panic(err)
	}
	return endpoint
}

func TestNewCall(t *testing.T) {
	type input struct {
		endpoint Endpoint
		params   []CallOpt
	}
	var tests = []struct {
		name     string
		in       input
		expected Call
	}{
		{
			name:     "nil input parameters",
			in:       input{nil, nil},
			expected: &defaultCall{},
		},
		{
			name:     "without arguments",
			in:       input{validEndpoint(), nil},
			expected: &defaultCall{endpoint: validEndpoint()},
		},
		{
			name:     "one argument",
			in:       input{validEndpoint(), []CallOpt{Field("a", "b")}},
			expected: &defaultCall{endpoint: validEndpoint(), fields: map[string]string{"a": "b"}},
		},
		{
			name:     "multiple argument",
			in:       input{validEndpoint(), []CallOpt{Field("a", "b"), Field("c", "d"), ProtoFile("file1"), ProtoFile("file2")}},
			expected: &defaultCall{endpoint: validEndpoint(), fields: map[string]string{"a": "b", "c": "d"}, protoFiles: []string{"file1", "file2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewCall(tt.in.endpoint, tt.in.params...)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
