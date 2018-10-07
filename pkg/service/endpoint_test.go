package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEndpoint(t *testing.T) {
	var tests = []struct {
		name     string
		in       string
		endpoint Endpoint
		err      error
	}{
		{
			name:     "empty string",
			in:       "",
			endpoint: nil,
			err:      ErrEmptyEndpoint,
		},
		{
			name:     "invalid endpoint",
			in:       "abc",
			endpoint: nil,
			err:      ErrInvalidEndpoint,
		},
		{
			name:     "valid endpoint",
			in:       "localhost:3000/DocumentService/FindById",
			endpoint: &defaultEndpoint{host: "localhost", port: 3000, service: "DocumentService", method: "FindById"},
			err:      nil,
		},
		{
			name:     "port out of range",
			in:       "localhost:70000/DocumentService/FindById",
			endpoint: nil,
			err:      ErrInvalidPort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewEndpoint(tt.in)
			assert.Equal(t, tt.endpoint, actual)
			if tt.err == nil { // no error expected
				assert.Nil(t, err)
			} else { // error expected
				assert.Error(t, err, tt.err.Error())
			}
		})
	}
}

func TestHostAndPort(t *testing.T) {
	e, err := NewEndpoint("localhost:30000/service/method")
	assert.Nil(t, err)
	assert.Equal(t, "localhost:30000", e.HostAndPort())
}
