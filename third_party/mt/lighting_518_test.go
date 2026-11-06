package mt

import (
	"bytes"
	"testing"
)

func TestLighting518RoundTrip(t *testing.T) {
	for _, set := range []bool{false, true} {
		in := &ToCltLighting{LightCurveSet: set, AmbientOcclusionGamma: 2.5}
		in.LightCurve[0], in.LightCurve[15] = 1, 16
		var wire bytes.Buffer
		if err := serialize(&wire, in); err != nil {
			t.Fatal(err)
		}
		if want := 69; set {
			want += 16
			if wire.Len() != want {
				t.Fatalf("packet length = %d, want %d", wire.Len(), want)
			}
		} else if wire.Len() != want {
			t.Fatalf("packet length = %d, want %d", wire.Len(), want)
		}
		var out ToCltLighting
		if err := deserialize(&wire, &out); err != nil {
			t.Fatal(err)
		}
		if out.LightCurveSet != set || out.AmbientOcclusionGamma != in.AmbientOcclusionGamma || (set && out.LightCurve != in.LightCurve) {
			t.Fatalf("lighting mismatch: %+v", out)
		}
	}
}

func TestLightingWithoutStaticFields(t *testing.T) {
	var wire bytes.Buffer
	if err := serialize(&wire, &ToCltLighting{}); err != nil {
		t.Fatal(err)
	}
	legacy := wire.Bytes()[:wire.Len()-5]
	var out ToCltLighting
	if err := deserialize(bytes.NewReader(legacy), &out); err != nil {
		t.Fatal(err)
	}
	if out.LightCurveSet || out.AmbientOcclusionGamma != 0 {
		t.Fatalf("unexpected static fields: %+v", out)
	}
}
