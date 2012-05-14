// A small wrapper around liblo (OSC) for (de)serializing OSC packets
package golo

// #include <stdlib.h>
// #cgo pkg-config: liblo
// #include <lo/lo.h>
// #include <lo/lo.h>
// #include "golo.h"
import "C"

import (
	"errors"
	"unsafe"
	"fmt"
)

var (
	ErrUnknown = errors.New("Unknown type")
)

// A message represents an OSC message.
// A parameter may be one of the following types:
//
// int64, int32, float64, float32, string
//
// OSC supports more types than this,
// but the code has not been ported yet.
type Message struct {
	Path   string
	Params []interface{}
}

// Takes a serialized OSC packet and deserializes it into a
// message. Parameters with an unsupported type will
// be deserialized as ErrUnknown.
func Deserialize(data []byte) (*Message, error) {
	cdata := unsafe.Pointer(&data[0])
	cdatalen := C.size_t(len(data))

	msg := C.lo_message_deserialise(cdata, cdatalen, (*C.int)(unsafe.Pointer(nil)))
	if msg == nil {
		// TODO: Obtain and parse actual error?
		return nil, errors.New("Deserialization failed")
	}

	result := &Message{
		Path:   C.GoString(C.lo_get_path(cdata, C.ssize_t(cdatalen))),
		Params: make([]interface{}, 0, 10),
	}
	argc := int(C.lo_message_get_argc(msg))
	for i := 0; i < argc; i++ {
		result.Params = append(result.Params, extractArgument(msg, i))
	}
	return result, nil
}

// Serializes a message to OSC wire format.
// If a unsupported type is encountered, serialization
// will be stopped.
func Serialize(m *Message) ([]byte, error) {
	msg := C.lo_message_new()
	for i, param := range m.Params {
		switch x := param.(type) {
		case int32:
			C.lo_message_add_int32(msg, C.int32_t(x))
		case int64:
			C.lo_message_add_int64(msg, C.int64_t(x))
		case float32:
			C.lo_message_add_float(msg, C.float(x))
		case float64:
			C.lo_message_add_double(msg, C.double(x))
		case string:
			cstr := C.CString(x)
			defer C.free(unsafe.Pointer(cstr))
			C.lo_message_add_string(msg, cstr)
		default:
			return nil, fmt.Errorf("Parameter %d has invalid type", i)
		}
	}

	cpath := C.CString(m.Path)
	defer C.free(unsafe.Pointer(cpath))
	var size int

	tmpbuffer := C.lo_message_serialise(msg, cpath, unsafe.Pointer(nil), (*C.size_t)(unsafe.Pointer((&size))))
	defer C.free(unsafe.Pointer(tmpbuffer))
	longbuffer := C.GoBytes(tmpbuffer, C.int(size))

	shortbuffer := make([]byte, size)
	copy(shortbuffer, longbuffer)
	return shortbuffer, nil
}

func extractArgument(msg C.lo_message, idx int) interface{} {
	argtypes := C.GoString(C.lo_message_get_types(msg))
	argv := C.lo_message_get_argv(msg)
	switch argtypes[idx] {
	case C.LO_INT32:
		return int32(C.msg_extract_int32(argv, C.int(idx)))
	case C.LO_INT64:
		return int64(C.msg_extract_int64(argv, C.int(idx)))
	case C.LO_FLOAT:
		return float32(C.msg_extract_float32(argv, C.int(idx)))
	case C.LO_DOUBLE:
		return float64(C.msg_extract_float64(argv, C.int(idx)))
	case C.LO_STRING:
		return C.GoString(C.msg_extract_string(argv, C.int(idx)))
	}
	return ErrUnknown
}
