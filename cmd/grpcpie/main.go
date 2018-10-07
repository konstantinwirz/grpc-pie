package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/konstantinwirz/grpc-pie/pkg/service"
	"google.golang.org/grpc"
)

func main() {
	protoFile := flag.String("proto", "", "path to the proto file")
	flag.Usage = usage
	flag.Parse()
	call, err := createServiceCall(*protoFile, flag.Args())
	if err != nil {
		fmt.Printf("failed to parse arguments: %v\n", err)
		os.Exit(1)
	}

	executor := service.NewExecutor()
	response, err := executor.Exec(call, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to make a rpc call: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("RESPONSE = %v\n", response)
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

func usage() {
	fmt.Printf("Usage of %s:\n\n", os.Args[0])
	fmt.Printf("%s [OPTIONS] host:port/service/method field1=value1 field2=value2 ... fieldN=valueN\n\n", os.Args[0])
	flag.PrintDefaults()
}

var (
	// ErrHostMissing signals that the sericce host is missing
	ErrHostMissing = errors.New("host is missing")
	// ErrInvalidFieldAssignment signal a argument assignment is invalid
	ErrInvalidFieldAssignment = errors.New("invalid argument assigment")
)

// creates a service call using given proto file and args that should contain
// endpoint and field assignments
func createServiceCall(protoFile string, args []string) (service.Call, error) {
	if len(args) == 0 {
		return nil, ErrHostMissing
	}

	// expect the first part of the string to be endpoint
	endpoint, err := service.NewEndpoint(args[0])
	if err != nil {
		return nil, err
	}

	fields := []service.CallOpt{}
	for _, arg := range args[1:] {
		var parts = strings.Split(arg, "=")
		if len(parts) != 2 {
			return nil, ErrInvalidFieldAssignment
		}

		fieldName, fieldValue := parts[0], parts[1]
		fields = append(fields, service.Field(fieldName, fieldValue))
	}

	// add proto file
	fields = append(fields, service.ProtoFile(protoFile))

	return service.NewCall(endpoint, fields...), nil
}
