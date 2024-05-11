package main

import (
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	fmt.Println("=== int32 ===")
	var data1 int32 = 268_435_456
	i32 := &wrapperspb.Int32Value{
		Value: data1,
	}
	d1, w1 := serializedSize(data1, i32)

	fmt.Printf("in memory: %d\npb: %d\n", d1, w1)

	fmt.Println("\n=== uint64 ===")
	var data2 uint64 = 72_057_594_037_927_936
	ui64 := &wrapperspb.UInt64Value{
		Value: data2,
	}
	d2, w2 := serializedSize(data2, ui64)

	fmt.Printf("in memory: %d\npb: %d\n", d2, w2)
}

func serializedSize[D constraints.Integer, W protoreflect.ProtoMessage](data D, wrapper W) (uintptr, int) {
	out, err := proto.Marshal(wrapper)

	if err != nil {
		log.Fatal(err)
	}

	return unsafe.Sizeof(data), len(out) - 1
}
