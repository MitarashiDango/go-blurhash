package blurhash

import "testing"

func TestPow83(t *testing.T) {
	tests := []struct {
		exp      int
		expected int
	}{
		{0, 1},
		{1, 83},
		{2, 83 * 83},
		{3, 83 * 83 * 83},
	}
	for _, tt := range tests {
		got := pow83(tt.exp)
		if got != tt.expected {
			t.Errorf("pow83(%d) = %d, want %d", tt.exp, got, tt.expected)
		}
	}
}

func TestEncodeBase83(t *testing.T) {
	tests := []struct {
		value    int
		length   int
		expected string
	}{
		{0, 1, "0"},
		{1, 1, "1"},
		{82, 1, "~"},
		{0, 2, "00"},
		{83, 2, "10"},
	}
	for _, tt := range tests {
		got := encodeBase83(tt.value, tt.length)
		if got != tt.expected {
			t.Errorf("encodeBase83(%d, %d) = %q, want %q", tt.value, tt.length, got, tt.expected)
		}
	}
}

func TestEncodeBase83_Length(t *testing.T) {
	for length := 1; length <= 4; length++ {
		got := encodeBase83(0, length)
		if len(got) != length {
			t.Errorf("encodeBase83(0, %d) returned length %d", length, len(got))
		}
	}
}
