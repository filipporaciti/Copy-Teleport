package cipher

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"encoding/pem"
	"fmt"
	"os"
)

var(
	bits int = 1024
	privateRSAKey *rsa.PrivateKey = generateRSAKey() // local private key
	RemotePublicRSAKey *rsa.PublicKey // remote private key
)
 

// Encrypt plaintext with local private key
//
// Input: plaintext
//
// Output: ciphertext, error (nil if no error)
func LocalRSAEncrypt(plaintext []byte) ([]byte, error) {
	
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &privateRSAKey.PublicKey, plaintext, []byte(""))
	if err != nil {
		fmt.Println("[Error] local RSA encrypt")
		return nil, err
	}

	return ciphertext, nil
}

// Encrypt plaintext with input key
//
// Input: key, plaintext
//
// Output: ciphertext, error (nil if no error)
func RSAEncrypt(key *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, plaintext, []byte(""))
	if err != nil {
		fmt.Println("[Error] local RSA encrypt")
		return nil, err
	}

	return ciphertext, nil
}

// Decrypt ciphertext with local private key
//
// Input: ciphertext
//
// Output: plaintext, error (nil if no error)
func LocalRSADecrypt(ciphertext []byte) ([]byte, error) {

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateRSAKey, ciphertext, []byte(""))
	if err != nil {
		fmt.Println("[Error] RSA decrypt")
		return nil, err
	}
	return plaintext, nil

}

// Decrypt ciphertext with input key
//
// Input: key, ciphertext
//
// Output: plaintext, error (nil if no error)
func RSADecrypt(key *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, ciphertext, []byte(""))
	if err != nil {
		fmt.Println("[Error] RSA decrypt: " + err.Error())
		return nil, err
	}
	return plaintext, nil

}

// Return new private key
//
// Input:
//
// Output: private key
func generateRSAKey() *rsa.PrivateKey {
	k, err := rsa.GenerateKey(rand.Reader,bits)
	if err != nil {
		fmt.Println("[Error] creating RSA keys")
		os.Exit(1)
	}
	return k
}

// Return public key of local private key with PEM certificate
//
// Input:
//
// Output: PEM public key
func EncodeRSAPublicKeyPEM(pk *rsa.PublicKey) string {
	pubkey_pem := string(pem.EncodeToMemory(&pem.Block{Type:  "RSA PUBLIC KEY",Bytes: x509.MarshalPKCS1PublicKey(pk)}))
    return pubkey_pem
}


func DecodeRSAPublicKeyPEM(p string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(p))
    if block == nil {
            return nil, errors.New("failed to parse PEM block containing the key")
    }

    priv, err := x509.ParsePKCS1PublicKey(block.Bytes)
    if err != nil {
            return nil, err
    }

    return priv, nil
}