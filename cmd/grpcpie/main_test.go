package main

import (
	"testing"

	"github.com/konstantinwirz/grpc-pie/pkg/service"
	"github.com/stretchr/testify/assert"
)

func validEndpoint(s string) service.Endpoint {
	endpoint, err := service.NewEndpoint(s)
	if err != nil {
		panic(err)
	}

	return endpoint
}

func TestCreateServiceCall(t *testing.T) {
	type input struct {
		protoFile string
		args      []string
	}
	var tests = []struct {
		name string
		in   input
		// expected results
		call service.Call
		err  error
	}{
		{
			name: "empty argument list",
			in:   input{"", []string{}},
			call: nil,
			err:  ErrHostMissing,
		},
		{
			name: "host without port number",
			in:   input{"", []string{"localhost"}},
			call: nil,
			err:  service.ErrInvalidEndpoint,
		},
		{
			name: "valid endpoint",
			in:   input{"file", []string{"localhost:30000/service/method"}},
			call: service.NewCall(validEndpoint("localhost:30000/service/method"), service.ProtoFile("file")),
			err:  nil,
		},
		{
			name: "valid endpoint/arguments",
			in:   input{"file", []string{"localhost:30000/service/method", "a=b", "c=d"}},
			call: service.NewCall(validEndpoint("localhost:30000/service/method"), service.Field("a", "b"), service.Field("c", "d"), service.ProtoFile("file")),
			err:  nil,
		},
		{
			name: "valid endpoint and invalid arguments",
			in:   input{"", []string{"localhost:30000/service/method", "a=b", "c"}},
			call: nil,
			err:  ErrInvalidFieldAssignment},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := createServiceCall(tt.in.protoFile, tt.in.args)
			assert.Equal(t, tt.call, actual)
			if tt.err == nil { // no error expected
				assert.NoError(t, err)
			} else { // error expected
				assert.EqualError(t, err, tt.err.Error())
			}
		})
	}
}
