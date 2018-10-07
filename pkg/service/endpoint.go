package service

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

// Endpoint represents a service endpoint
type Endpoint interface {
	Host() string
	Port() uint16
	HostAndPort() string
	Service() string
	Method() string
}

// defaultEndpoint is the default implementation
// of the Endpoint
type defaultEndpoint struct {
	host, service, method string
	port                  uint16
}

func (e *defaultEndpoint) Host() string {
	return e.host
}

func (e *defaultEndpoint) Port() uint16 {
	return e.port
}

func (e *defaultEndpoint) HostAndPort() string {
	return e.host + ":" + strconv.Itoa(int(e.port))
}

func (e *defaultEndpoint) Service() string {
	return e.service
}

func (e *defaultEndpoint) Method() string {
	return e.method
}

func (e *defaultEndpoint) String() string {
	return fmt.Sprintf("Endpoint {host = %s, port = %d, service = %s, method = %s}",
		e.host, e.port, e.service, e.method)
}

var (
	// ErrEmptyEndpoint siganl that a empty string as endpoint is passed
	ErrEmptyEndpoint = errors.New("empty endpoint")
	// ErrInvalidEndpoint signal that a endpoint isn't well formed
	ErrInvalidEndpoint = errors.New("invalid endpoint")
	// ErrInvalidPort signals that the port is out of range
	ErrInvalidPort = errors.New("invalid port")
)

// NewEndpoint creates and returns an Endpoint instance
// parsed from the given string
func NewEndpoint(endpoint string) (Endpoint, error) {
	if endpoint == "" {
		return nil, ErrEmptyEndpoint
	}

	var r = regexp.MustCompile("^(\\w+):(\\d+)/(\\w+)/(\\w+)")
	var groups = r.FindStringSubmatch(endpoint)
	if groups == nil {
		return nil, ErrInvalidEndpoint
	}

	port, err := strconv.ParseUint(groups[2], 10, 16)
	if err != nil || port > math.MaxUint16 {
		return nil, ErrInvalidPort
	}

	return &defaultEndpoint{
		host:    groups[1],
		port:    uint16(port),
		service: groups[3],
		method:  groups[4]}, nil
}
