package mt

import (
	"bytes"
	"testing"
)

func TestAOPositionRotationFlag(t *testing.T) {
	for _, want := range []bool{false, true} {
		in := &IDAOMsg{ID: 12, Msg: &AOCmdPos{Pos: AOPos{InterpolateRotation: want}}}
		var wire bytes.Buffer
		if err := serialize(&wire, in); err != nil {
			t.Fatal(err)
		}
		var out IDAOMsg
		if err := deserialize(&wire, &out); err != nil {
			t.Fatal(err)
		}
		if got := out.Msg.(*AOCmdPos).Pos.InterpolateRotation; got != want {
			t.Fatalf("rotation flag = %v, want %v", got, want)
		}
	}
}

func TestAOPositionWithoutRotationFlag(t *testing.T) {
	var wire bytes.Buffer
	if err := serialize(&wire, &AOPos{InterpolateRotation: true}); err != nil {
		t.Fatal(err)
	}
	legacy := wire.Bytes()[:wire.Len()-1]
	var out AOPos
	if err := deserialize(bytes.NewReader(legacy), &out); err != nil {
		t.Fatal(err)
	}
	if !out.InterpolateRotation {
		t.Fatal("legacy position should interpolate rotation")
	}
}
