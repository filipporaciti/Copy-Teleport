package cipher

import(
	"bytes"
	"errors"
)


// Add padding to input
// Input: src (input bytes), size (padding block size)
// Output: (src+padding)
func Pad(src []byte, size int) []byte {
    padding := size - len(src)%size
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(src, padtext...)
}

// Remove padding from src (with pad)
// Input: src (with pad)
// Output: src (without padding), error (nil if no error)
func Unpad(src []byte) ([]byte, error) {
    length := len(src)
    unpadding := int(src[length-1])

    if unpadding > length {
        return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
    }

    return src[:(length - unpadding)], nil
}