package blurhash

const base83Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz#$%*+,-.:;=?@[]^_{|}~"

func encodeBase83(value, length int) string {
	result := make([]byte, length)
	for i := 1; i <= length; i++ {
		digit := (value / pow83(length-i)) % 83
		result[i-1] = base83Chars[digit]
	}
	return string(result)
}

func pow83(exp int) int {
	result := 1
	for range exp {
		result *= 83
	}
	return result
}
