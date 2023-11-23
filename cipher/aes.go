package cipher


import(
		"crypto/aes"
		"crypto/cipher"
crand 	"crypto/rand"
		"math/rand"
		"time"
		"fmt"
		"io"
)


var( 
	privateAESKey []byte = GenerateAESKey() // local private key
)


// Decrypt ciphertext with local private key
//
// Input: ciphertext
//
// Output: plaintext, error (nil if no error)
func LocalAESDecrypt(ciphertext []byte) ([]byte, error) {

	ciphertext, err := Base64ToByte(string(ciphertext))

	block, err := aes.NewCipher(privateAESKey)
	if err != nil {
		fmt.Println("\033[31m[Error] AES decrypt creation cipher:", err.Error(), "\033[0m")
		return nil, err
	}
	
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]


	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ciphertext, ciphertext)


	plaintext, err := Unpad(ciphertext)
	if err != nil {
		fmt.Println("\033[31m[Error] AES decrypt unpad:", err.Error(), "\033[0m")
		return nil, err
	}


	return plaintext, nil


}


// Decrypt ciphertext with input key
//
// Input: key, ciphertext
//
// Output: plaintext, error (nil if no error)
func AESDecrypt(key []byte, ciphertext []byte) ([]byte, error) {

	ciphertext, err := Base64ToByte(string(ciphertext))

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("\033[31m[Error] AES decrypt creation cipher:", err.Error(), "\033[0m")
		return nil, err
	}
	
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ciphertext, ciphertext)

	// Unpad the plaintext
	plaintext, err := Unpad(ciphertext)
	if err != nil {
		fmt.Println("\033[31m[Error] AES decrypt unpad:", err.Error(), "\033[0m")
		return nil, err
	}

	return plaintext, nil

}

// Encrypt plaintext with local private key
//
// Input: plaintext
//
// Output: ciphertext, error (nil if no error)
func LocalAESEncrypt(plaintext []byte) ([]byte, error) {

	plaintext = Pad(plaintext, aes.BlockSize)

	block, err := aes.NewCipher(privateAESKey)
	if err != nil {
		fmt.Println("\033[31m[Error] AES encrypt creation cipher:", err.Error(), "\033[0m")
		return nil, err
	}

	ciphertext := make([]byte, len(plaintext))
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		fmt.Println("\033[31m[Error] AES encrypt generate iv:", err.Error(), "\033[0m")
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext, plaintext)

	ciphertext = append(iv, ciphertext...)

	ciphertext = []byte(ByteToBase64(ciphertext))

	return ciphertext, nil

}

// Encrypt plaintext with input key
//
// Input: key, plaintext
//
// Output: ciphertext, error (nil if no error)
func AESEncrypt(key []byte, plaintext []byte) ([]byte, error) {

	plaintext = Pad(plaintext, aes.BlockSize)

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("\033[31m[Error] AES encrypt creation cipher:", err.Error(), "\033[0m")
		return nil, err
	}

	ciphertext := make([]byte, len(plaintext))
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		fmt.Println("\033[31m[Error] AES encrypt generate iv:", err.Error(), "\033[0m")
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext, plaintext)

	ciphertext = append(iv, ciphertext...)

	return ciphertext, nil

}


// Return new random private key
//
// Input:
//
// Output: private key
func GenerateAESKey() []byte {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 16)
    for i := range b {
        b[i] = byte(rand.Intn(255))
    }
    return b
}


