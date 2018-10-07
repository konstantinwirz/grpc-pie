package service

import "fmt"

// Call represents a service call
type Call interface {
	// Endpoint returns the endpoint of this call
	Endpoint() Endpoint
	// Fields returns all the field names of this call
	Fields() []string
	// FieldValue returns the value of the given argument
	FieldValue(field string) (string, bool)
	// ProtoFiles returns proto files describing this call
	ProtoFiles() []string
}

// defaultCall is a default implementaion of
// the Call interface
type defaultCall struct {
	endpoint   Endpoint
	fields     map[string]string
	protoFiles []string
}

func (c *defaultCall) Endpoint() Endpoint {
	return c.endpoint
}

func (c *defaultCall) Fields() []string {
	sz := len(c.fields)
	fields := make([]string, sz)
	var i int
	for k := range c.fields {
		fields[i] = k
		i++
	}

	return fields
}

func (c *defaultCall) FieldValue(field string) (string, bool) {
	value, ok := c.fields[field]
	return value, ok
}

func (c *defaultCall) String() string {
	return fmt.Sprintf("ServiceCall {endpoint: %v, fields: %v}", c.endpoint, c.fields)
}

func (c *defaultCall) ProtoFiles() []string {
	return c.protoFiles
}

// CallOpt is a option type for the Call constructor
type CallOpt func(c *defaultCall) *defaultCall

// Field represents a call parameter, a key-value pair
func Field(key, value string) CallOpt {
	return func(c *defaultCall) *defaultCall {
		if c.fields == nil {
			c.fields = make(map[string]string)
		}

		c.fields[key] = value
		return c
	}
}

// ProtoFile allows to add a proto file
func ProtoFile(file string) CallOpt {
	return func(c *defaultCall) *defaultCall {
		if c.protoFiles == nil {
			c.protoFiles = make([]string, 1)
			c.protoFiles[0] = file
		} else {
			c.protoFiles = append(c.protoFiles, file)
		}
		return c
	}
}

// NewCall creates and returrns a configured Call
func NewCall(endpoint Endpoint, opts ...CallOpt) Call {
	var call = &defaultCall{endpoint: endpoint}

	for _, opt := range opts {
		call = opt(call)
	}

	return call
}
