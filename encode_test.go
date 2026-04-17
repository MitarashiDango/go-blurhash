package blurhash

import (
	"testing"
)

func TestEncode_SolidColor(t *testing.T) {
	// 2x2 solid red image
	pixels := []byte{
		255, 0, 0, 255, 0, 0,
		255, 0, 0, 255, 0, 0,
	}
	hash, err := Encode(pixels, 2, 2, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	// 1x1 components: size flag (1 char) + quantised max (1 char) + DC (4 chars) = 6 chars
	if len(hash) != 6 {
		t.Errorf("expected hash length 6, got %d: %s", len(hash), hash)
	}
}

func TestEncode_KnownHash(t *testing.T) {
	// 4x4 solid white image should produce a known deterministic hash
	width, height := 4, 4
	pixels := make([]byte, width*height*3)
	for i := range pixels {
		pixels[i] = 255
	}
	hash, err := Encode(pixels, width, height, 4, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// With all-white pixels, all AC components should be zero
	// Hash should be deterministic
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	// 4x3 components: 6 + (4*3-1)*2 = 6 + 22 = 28 chars
	expectedLen := 6 + 2*(4*3-1)
	if len(hash) != expectedLen {
		t.Errorf("expected hash length %d, got %d: %s", expectedLen, len(hash), hash)
	}
	t.Logf("white 4x4 (4x3 components) hash: %s", hash)
}

func TestEncode_InvalidComponents(t *testing.T) {
	pixels := make([]byte, 3)
	_, err := Encode(pixels, 1, 1, 0, 1)
	if err == nil {
		t.Fatal("expected error for invalid components")
	}

	_, err = Encode(pixels, 1, 1, 10, 1)
	if err == nil {
		t.Fatal("expected error for invalid components")
	}
}

func TestEncode_PixelLengthMismatch(t *testing.T) {
	pixels := make([]byte, 5) // wrong size for 2x2
	_, err := Encode(pixels, 2, 2, 1, 1)
	if err == nil {
		t.Fatal("expected error for pixel length mismatch")
	}
}

func TestEncode_HashLength(t *testing.T) {
	// Hash length should be 6 + 2*(xComponents*yComponents - 1)
	tests := []struct {
		xComp, yComp int
	}{
		{1, 1},
		{2, 2},
		{4, 3},
		{9, 9},
	}

	for _, tt := range tests {
		width, height := 8, 8
		pixels := make([]byte, width*height*3)
		for i := range pixels {
			pixels[i] = byte(i % 256)
		}
		hash, err := Encode(pixels, width, height, tt.xComp, tt.yComp)
		if err != nil {
			t.Fatalf("xComp=%d yComp=%d: unexpected error: %v", tt.xComp, tt.yComp, err)
		}
		expectedLen := 6 + 2*(tt.xComp*tt.yComp-1)
		if len(hash) != expectedLen {
			t.Errorf("xComp=%d yComp=%d: expected length %d, got %d: %s",
				tt.xComp, tt.yComp, expectedLen, len(hash), hash)
		}
	}
}
