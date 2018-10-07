package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
)

// CallExecutor executes calls
type CallExecutor interface {
	Exec(Call, ...grpc.DialOption) (proto.Message, error)
}

type defaultCallExecutor struct {
}

// NewExecutor creates and returns a CallExecutor instance
func NewExecutor() CallExecutor {
	return &defaultCallExecutor{}
}

func (e *defaultCallExecutor) Exec(call Call, dialOpts ...grpc.DialOption) (proto.Message, error) {
	descriptors, err := fileDescriptors(call.ProtoFiles())
	if err != nil {
		return nil, fmt.Errorf("failed to parse proto files: %v", err.Error())
	}

	_, svcDescr, err := findService(call, descriptors)
	if err != nil {
		return nil, err
	}

	// get the method descriptor
	methodDescr := svcDescr.FindMethodByName(call.Endpoint().Method())
	if methodDescr == nil {
		return nil, ErrMethodNotFound
	}

	// get message descriptor for the request
	messageDescr := methodDescr.GetInputType()
	if messageDescr == nil {
		return nil, ErrMessageNotFound
	}

	// now create a message
	msg := dynamic.NewMessage(messageDescr)
	if msg == nil {
		return nil, ErrMessageCreationFailed
	}

	msg, err = prepareMessage(call, msg)
	if err != nil {
		return nil, err
	}

	// establish a connection
	conn, err := grpc.Dial(call.Endpoint().HostAndPort(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	stub := grpcdynamic.NewStub(conn)
	response, err := stub.InvokeRpc(context.Background(), methodDescr, msg)
	if err != nil {
		fmt.Printf("failed to make a rpc call: %v", err)
		os.Exit(1)
	}

	return response, nil
}

func fileDescriptors(files []string) ([]*desc.FileDescriptor, error) {
	parser := protoparse.Parser{}
	return parser.ParseFiles(files...)
}

var (
	// ErrServiceNotFound signals that a service couldn't be found
	ErrServiceNotFound = errors.New("service not found")
	// ErrMethodNotFound signals that a method couldn't be found
	ErrMethodNotFound = errors.New("method not found")
	// ErrMessageNotFound signals that a message couldn't be found
	ErrMessageNotFound = errors.New("message not found")
	// ErrMessageCreationFailed signals that a message couldn't be created
	ErrMessageCreationFailed = errors.New("failed to create a message")
	// ErrFieldNotFound signals that a field couldn't be found
	ErrFieldNotFound = errors.New("field not found")
)

func findService(call Call, descriptors []*desc.FileDescriptor) (*desc.FileDescriptor, *desc.ServiceDescriptor, error) {
	for _, descr := range descriptors {
		serviceName := qualifyServiceName(descr.GetPackage(), call.Endpoint().Service())
		svcDescr := descr.FindService(serviceName)
		if svcDescr != nil {
			return descr, svcDescr, nil
		}
	}

	// nothing found
	return nil, nil, ErrServiceNotFound
}

func qualifyServiceName(pkg, service string) string {
	if len(pkg) > 0 {
		return pkg + "." + service
	}

	return service
}

func prepareMessage(call Call, msg *dynamic.Message) (*dynamic.Message, error) {
	// iterate over fields
	for _, field := range call.Fields() {
		value, ok := call.FieldValue(field)
		if !ok {
			panic("no value for field " + field)
		}

		fieldDescr := msg.FindFieldDescriptorByName(field)
		if fieldDescr == nil {
			return nil, ErrFieldNotFound
		}

		setFieldValue(msg, fieldDescr, value)
	}

	return msg, nil
}

func setFieldValue(msg *dynamic.Message, fd *desc.FieldDescriptor, value string) {
	switch fd.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			panic(err)
		}
		msg.SetField(fd, v)
	default:
		msg.SetField(fd, value)
	}
}
