package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"
	"io"
)

// Encrypt function accept 32 byte key to encrypt plaintext
// ciphertext as output
func Encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt function decypt the encrypted cipher text
func Decrypt(ciphertext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(decodedCiphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext is too short")
	}

	nonce := decodedCiphertext[:gcm.NonceSize()]
	decodedCiphertext = decodedCiphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, decodedCiphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Hash calculates a hash of the given message using the specified algorithm like MD5,SHA-256,SHA-2 etc(default is MD5).
func Hash(message string, algorithms ...string) (string, error) {
	var hasher hash.Hash

	var algorithm string

	if len(algorithms) > 0 {
		algorithm = algorithms[0]
	}

	switch algorithm {
	case "sha256":
		hasher = sha256.New()
	default:
		hasher = md5.New()
	}

	_, err := hasher.Write([]byte(message))
	if err != nil {
		return "", err
	}

	hashBytes := hasher.Sum(nil)
	hashHex := hex.EncodeToString(hashBytes)

	return hashHex, nil
}
