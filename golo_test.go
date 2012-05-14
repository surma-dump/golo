package golo

import (
	"testing"
)

func TestIdentity(t *testing.T) {
	in_msg := &Message{
		Path:   "/osc/path",
		Params: make([]interface{}, 5),
	}
	in_msg.Params[0] = float32(0.0)
	in_msg.Params[1] = float64(1.0)
	in_msg.Params[2] = int32(2.0)
	in_msg.Params[3] = int64(3.0)
	in_msg.Params[4] = "4"

	data, e := Serialize(in_msg)
	if e != nil {
		t.Fatalf("Could not serialize message: %s", e)
	}

	out_msg, e := Deserialize(data)
	if e != nil {
		t.Fatalf("Could not deserialize message: %s", e)
	}

	if out_msg.Path != in_msg.Path {
		t.Fatalf("Paths do not match (\"%s\" != \"%s\")", out_msg.Path, in_msg.Path)
	}

	for i := range in_msg.Params {
		if in_msg.Params[i] != out_msg.Params[i] {
			t.Fatalf("Paramerter %d does not match (%v != %v)", i, in_msg.Params[i], out_msg.Params[i])
		}
	}
}
