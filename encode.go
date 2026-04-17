package blurhash

import (
	"fmt"
	"math"
	"strings"
)

// Encode generates a Blurhash string from pixel data.
//
// pixels is a flat array of RGB values in row-major order (length = width * height * 3).
// xComponents and yComponents control the level of detail (1-9 each).
func Encode(pixels []byte, width, height, xComponents, yComponents int) (string, error) {
	if xComponents < 1 || xComponents > 9 || yComponents < 1 || yComponents > 9 {
		return "", fmt.Errorf("components must be between 1 and 9, got x=%d y=%d", xComponents, yComponents)
	}
	if len(pixels) != width*height*3 {
		return "", fmt.Errorf("pixel data length mismatch: expected %d, got %d", width*height*3, len(pixels))
	}

	factors := make([][3]float64, xComponents*yComponents)

	for j := range yComponents {
		for i := range xComponents {
			var r, g, b float64
			for y := range height {
				for x := range width {
					basis := math.Cos(math.Pi*float64(i)*float64(x)/float64(width)) *
						math.Cos(math.Pi*float64(j)*float64(y)/float64(height))
					idx := (y*width + x) * 3
					r += basis * sRGBToLinear(int(pixels[idx]))
					g += basis * sRGBToLinear(int(pixels[idx+1]))
					b += basis * sRGBToLinear(int(pixels[idx+2]))
				}
			}
			scale := 1.0 / float64(width*height)
			if i != 0 || j != 0 {
				scale = 2.0 / float64(width*height)
			}
			factors[j*xComponents+i] = [3]float64{r * scale, g * scale, b * scale}
		}
	}

	var buf strings.Builder

	// Size flag
	sizeFlag := (xComponents - 1) + (yComponents-1)*9
	buf.WriteString(encodeBase83(sizeFlag, 1))

	// Quantised maximum value
	var maximumValue float64
	if len(factors) > 1 {
		var actualMaximumValue float64
		for _, f := range factors[1:] {
			for _, v := range f {
				if math.Abs(v) > actualMaximumValue {
					actualMaximumValue = math.Abs(v)
				}
			}
		}
		quantisedMaximumValue := int(math.Max(0, math.Min(82, math.Floor(actualMaximumValue*166-0.5))))
		maximumValue = (float64(quantisedMaximumValue) + 1) / 166.0
		buf.WriteString(encodeBase83(quantisedMaximumValue, 1))
	} else {
		maximumValue = 1
		buf.WriteString(encodeBase83(0, 1))
	}

	// DC value
	dc := factors[0]
	buf.WriteString(encodeBase83(encodeDC(dc[0], dc[1], dc[2]), 4))

	// AC values
	for _, f := range factors[1:] {
		buf.WriteString(encodeBase83(encodeAC(f[0], f[1], f[2], maximumValue), 2))
	}

	return buf.String(), nil
}
