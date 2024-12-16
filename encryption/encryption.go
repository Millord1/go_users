package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

var gcm cipher.AEAD

func EncryptData(data string) (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Reader.Read(key); err != nil {
		fmt.Println("error generating random encryption key ", err)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("error creating aes block cipher", err)
		return "", err
	}

	cipherGcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("error setting gcm mode", err)
		return "", err
	}
	gcm = cipherGcm

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("error generating the nonce ", err)
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)

	return hex.EncodeToString(ciphertext), nil
}

func DecryptData(enc string) (string, error) {
	decodedCipherText, err := hex.DecodeString(enc)
	if err != nil {
		fmt.Println("error decoding hex", err)
		return "", err
	}

	decryptedData, err := gcm.Open(nil, decodedCipherText[:gcm.NonceSize()], decodedCipherText[gcm.NonceSize():], nil)
	if err != nil {
		fmt.Println("error decrypting data", err)
		return "", err
	}

	return string(decryptedData), nil
}

/* func EncryptData(data []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	// Using Galois Counter Mode (GCM) for authenticated encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nil, nonce, data, nil)
	return ciphertext, nil
}

func DecryptData(cipherText []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, fmt.Errorf("ciphertext is too short")
	}
	nonce, ciphertext := cipherText[:nonceSize], cipherText[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
} */
