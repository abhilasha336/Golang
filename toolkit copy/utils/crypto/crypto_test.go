package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

var minNonceLength = 12

func TestEncryptDecrypt(t *testing.T) {
	// Define a slice of test cases
	testCases := []struct {
		plaintext string
	}{
		{"This is a secret message."},
		{"Another message."},
		{"1234567890"},
		{"a"}, // Empty plaintext
	}

	// Replace with your 32-byte encryption key
	key := []byte("your_32_byte_encryption_key_here")

	// Iterate through the test cases
	for _, tc := range testCases {
		// Encrypt the plaintext
		encryptedText, err := Encrypt(tc.plaintext, key)
		if err != nil {
			require.Error(t, err, "Encrypt error for plaintext")
		}

		// Decrypt the ciphertext
		decryptedText, err := Decrypt(encryptedText, key)
		if err != nil {
			require.Error(t, err, "Decrypt error for ciphertext")
		}

		// Check that the decrypted text matches the original plaintext
		if decryptedText != tc.plaintext {
			require.Error(t, err, "Decrypted text  does not match original plaintext")
		}

		noNonceCiphertext := base64.StdEncoding.EncodeToString([]byte("invalid_ciphertext"))
		_, err = Decrypt(noNonceCiphertext, key)
		if err == nil {
			require.Error(t, err, "invalid_ciphertext")
		}

		// Test nonce generation
		nonce, err := generateNonce()
		if err != nil {
			require.Error(t, err, "Test nonce generation")
		}
		if len(nonce) != minNonceLength {
			require.Error(t, err)
		}
		// Test case with invalid base64
		invalidBase64 := "invalid_base64"
		_, err = Decrypt(invalidBase64, key)
		if err == nil {
			require.Error(t, err, "Decrypt should return an error with invalid base64")
		}

		// Test case with short ciphertext (should return an error)
		shortCiphertext := "s" // Ciphertext shorter than nonce size
		_, err = Decrypt(shortCiphertext, key)
		if err == nil {
			require.Error(t, err, "Decrypt should return an error with invalid base64")
		}

	}

	testCasesFalse := []struct {
		plaintext string
	}{
		{"This is a secret message."},
		{"Another message."},
		{"1234567890"},
		{"a"}, // Empty plaintext
	}

	// Replace with your 32-byte encryption key
	keyNew := []byte("your_32_byte_encr")

	// Iterate through the test cases
	for _, tc := range testCasesFalse {
		// Encrypt the plaintext
		encryptedText, err := Encrypt(tc.plaintext, keyNew)
		require.Error(t, err, "Encrypt the plaintext")

		// Decrypt the ciphertext
		decryptedText, err := Decrypt(encryptedText, keyNew)
		require.Error(t, err, "Decrypt the ciphertext")

		// Check that the decrypted text matches the original plaintext
		if decryptedText != tc.plaintext {
			require.Error(t, err, " Check that the decrypted text matches the original plaintext")
		}
		// Test case with ciphertext that has no nonce
		noNonceCiphertext := base64.StdEncoding.EncodeToString([]byte("invalid_ciphertext"))
		_, err = Decrypt(noNonceCiphertext, key)
		if err == nil {
			require.Error(t, err, "Test case with ciphertext that has no nonce")
		}

		// Test nonce generation
		nonce, err := generateNonce()
		if err != nil {
			require.Error(t, err, " Test nonce generation")
		}
		if len(nonce) != minNonceLength {
			require.Error(t, err, "Test length of nonce generation")
		}
	}

}
func generateNonce() ([]byte, error) {
	nonce := make([]byte, minNonceLength) // Adjust the size as needed
	_, err := io.ReadFull(rand.Reader, nonce)
	return nonce, err
}

// TestHash is a unit test for the Hash function .
// It tests the hash generation of the given message using the specified algorithm.
func TestHash(t *testing.T) {
	// Test cases with input messages, algorithms, and expected hash results
	testCases := []struct {
		message          string
		algorithm        string
		expectedHash     string
		expectedHashSize int
	}{
		{"Hello, World!", "md5", "65a8e27d8879283831b664bd8b7f0ad4", 16},
		{"Test", "sha256", "532eaabd9574880dbf76b9b8cc00832c20a6ec113d682299550d7a6e0f345e25", 32},
	}

	// Iterate through test cases and run the test for each case
	for _, tc := range testCases {
		hash, err := Hash(tc.message, tc.algorithm)

		// Check if the result matches the expected hash value
		if err != nil {
			t.Errorf("Error hashing message: %v", err)
		} else if hash != tc.expectedHash {
			t.Errorf("For message '%s' and algorithm '%s', expected hash '%s', but got '%s'", tc.message, tc.algorithm, tc.expectedHash, hash)
		}

		// Check if the hash size matches the expected size
		if len(hash)/2 != tc.expectedHashSize {
			t.Errorf("For message '%s' and algorithm '%s', expected hash size %d bytes, but got %d bytes", tc.message, tc.algorithm, tc.expectedHashSize, len(hash)/2)
		}
	}
}
