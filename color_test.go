package blurhash

import (
	"math"
	"testing"
)

func TestSRGBLinearRoundTrip(t *testing.T) {
	for v := 0; v <= 255; v++ {
		linear := sRGBToLinear(v)
		if linear < 0 || linear > 1 {
			t.Errorf("sRGBToLinear(%d) = %f, out of range [0, 1]", v, linear)
		}
		back := linearToSRGB(linear)
		if abs(back-v) > 1 {
			t.Errorf("round trip failed for %d: got %d", v, back)
		}
	}
}

func TestSRGBToLinear_Bounds(t *testing.T) {
	if got := sRGBToLinear(0); got != 0 {
		t.Errorf("sRGBToLinear(0) = %f, want 0", got)
	}
	if got := sRGBToLinear(255); math.Abs(got-1.0) > 0.001 {
		t.Errorf("sRGBToLinear(255) = %f, want ~1.0", got)
	}
}

func TestLinearToSRGB_Clamp(t *testing.T) {
	if got := linearToSRGB(-0.5); got != 0 {
		t.Errorf("linearToSRGB(-0.5) = %d, want 0", got)
	}
	if got := linearToSRGB(1.5); got != 255 {
		t.Errorf("linearToSRGB(1.5) = %d, want 255", got)
	}
}

func TestSignPow(t *testing.T) {
	if got := signPow(4, 0.5); math.Abs(got-2.0) > 0.001 {
		t.Errorf("signPow(4, 0.5) = %f, want 2.0", got)
	}
	if got := signPow(-4, 0.5); math.Abs(got+2.0) > 0.001 {
		t.Errorf("signPow(-4, 0.5) = %f, want -2.0", got)
	}
}

func TestEncodeDC(t *testing.T) {
	// All black
	got := encodeDC(0, 0, 0)
	if got != 0 {
		t.Errorf("encodeDC(0,0,0) = %d, want 0", got)
	}
	// All white
	got = encodeDC(1.0, 1.0, 1.0)
	expected := (255 << 16) + (255 << 8) + 255
	if got != expected {
		t.Errorf("encodeDC(1,1,1) = %d, want %d", got, expected)
	}
}

func TestEncodeAC_Range(t *testing.T) {
	// AC value should be in range [0, 18*19*19+18*19+18 = 6858]
	got := encodeAC(0, 0, 0, 1.0)
	if got < 0 || got > 18*19*19+18*19+18 {
		t.Errorf("encodeAC(0,0,0,1.0) = %d, out of range", got)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
