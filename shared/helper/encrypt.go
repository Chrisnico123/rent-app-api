package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"rent-application/configs"
)

func GetMasterKey() []byte {
	keyHex := configs.Cfg.SecretKey.Key
	key, _ := hex.DecodeString(keyHex)
	return key
}

func Encrypt(plaintext string) (string, string, error) {
	key := GetMasterKey()
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)

	iv := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", "", err
	}

	ciphertext := gcm.Seal(nil, iv, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), hex.EncodeToString(iv), nil
}

func Decrypt(ciphertextHex string, ivHex string) (string, error) {
	key := GetMasterKey()
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)

	ciphertext, _ := hex.DecodeString(ciphertextHex)
	iv, _ := hex.DecodeString(ivHex)

	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
