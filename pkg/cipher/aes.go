package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func Encrypt(data []byte, passphrase string) ([]byte, error) {

	var result []byte

	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return result, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return result, err
	}

	result = gcm.Seal(nonce, nonce, data, nil)

	return result, nil
}

func Decrypt(data []byte, passphrase string) ([]byte, error) {

	var result []byte

	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return result, err
	}

	nonceSize := gcm.NonceSize()

	if nonceSize > len(data)-1 {
		return result, fmt.Errorf("Invalid encrypted text")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	result, err = gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return result, err
	}

	return result, nil
}

func createHash(key string) string {

	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
