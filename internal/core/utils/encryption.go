package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"log"
)

var EncryptionKey []byte

func InitEncryption(key string) {
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		log.Fatal("Invalid encryption key: ", err)
	}

	EncryptionKey = keyBytes
}

//https://stackoverflow.com/questions/18817336/encrypting-a-string-with-aes-and-base64

func EncryptString(text []byte) ([]byte, error) {
	block, err := aes.NewCipher(EncryptionKey)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func DecryptString(text []byte) ([]byte, error) {
	block, err := aes.NewCipher(EncryptionKey)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Encode raw bytes into UTF-8
func EncryptToString(plaintext []byte) (string, error) {
	encrypted, err := EncryptString(plaintext)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func DecryptFromString(encoded string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return DecryptString(decoded)
}
