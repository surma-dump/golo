// A small wrapper around liblo (OSC) for (de)serializing OSC packets
package golo

// #cgo pkg-config: liblo
// #include <lo/lo.h>
// #include <lo/lo.h>
// #include "golo.h"
import "C"

import (
	"errors"
	"unsafe"
)

var (
	ErrUnknown = errors.New("Unknown type")
)

// A message represents an OSC message.
// A parameter may be one of the following types:
//
// int64, int32, float64, float32
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

func extractArgument(msg C.lo_message, idx int) interface{} {
	argtypes := C.GoString(C.lo_message_get_types(msg))
	argv := C.lo_message_get_argv(msg)
	switch argtypes[idx] {
	case 'i':
		return int32(C.msg_extract_int32(argv, C.int(idx)))
	case 'h':
		return int64(C.msg_extract_int64(argv, C.int(idx)))
	case 'f':
		return float32(C.msg_extract_float32(argv, C.int(idx)))
	case 'd':
		return float64(C.msg_extract_float64(argv, C.int(idx)))
	}
	return ErrUnknown
}
