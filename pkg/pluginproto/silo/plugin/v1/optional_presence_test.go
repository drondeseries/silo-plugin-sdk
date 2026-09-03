package pluginv1

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func assertOptionalInt32Presence(t *testing.T, withZero proto.Message, empty proto.Message, get func(proto.Message) *int32) {
	t.Helper()

	data, err := proto.Marshal(withZero)
	if err != nil {
		t.Fatalf("marshal %T with explicit zero: %v", withZero, err)
	}
	decoded := withZero.ProtoReflect().New().Interface()
	if err := proto.Unmarshal(data, decoded); err != nil {
		t.Fatalf("unmarshal %T with explicit zero: %v", withZero, err)
	}
	value := get(decoded)
	if value == nil || *value != 0 {
		t.Fatalf("%T explicit zero = %v, want present value 0", decoded, value)
	}

	data, err = proto.Marshal(empty)
	if err != nil {
		t.Fatalf("marshal empty %T: %v", empty, err)
	}
	decoded = empty.ProtoReflect().New().Interface()
	if err := proto.Unmarshal(data, decoded); err != nil {
		t.Fatalf("unmarshal empty %T: %v", empty, err)
	}
	if value := get(decoded); value != nil {
		t.Fatalf("%T absent value = %v, want nil", decoded, value)
	}
}
