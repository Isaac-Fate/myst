package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"
)

// Reference: https://gist.github.com/tscholl2/dc7dc15dc132ea70a98e8542fefffa28

// The encrypted password is a hex string which consists of 3 parts:
//
// 1. ciphertext: ciphertext of the password
// 2. salt: randomly generated salt
// 3. iv: initialization vector
//
// The 3 parts are separated by '-':
// <ciphertext>-<salt>-<iv>

// func Encrypt(passphrase string, password string) string {}

const saltLength int = 32
const ivLength int = 12
const secretKeyLength int = 32
const separator string = "-"

func Encrypt(passphrase string, password string) (string, error) {
	// Derive the secret key from the passphrase and generate a random salt
	key, salt := deriveKey(passphrase, nil)

	// Generate an initialization vector
	iv := make([]byte, ivLength)
	rand.Read(iv)

	// Create an AES block cipher
	blockCipher, err := aes.NewCipher(key)

	if err != nil {
		return "", err
	}

	// Wrap the block cipher in GCM mode
	gcmCipher, err := cipher.NewGCM(blockCipher)

	if err != nil {
		return "", err
	}

	// Encrypt the password
	ciphertext := gcmCipher.Seal(nil, iv, []byte(password), nil)

	// Create the encrypted password
	encryptedPassword := strings.Join(
		[]string{hex.EncodeToString(ciphertext), hex.EncodeToString(salt), hex.EncodeToString(iv)},
		separator,
	)

	return encryptedPassword, nil
}

func Decrypt(passphrase string, encryptedPassword string) (string, error) {
	// Separate the salt, iv, and ciphertext
	parts := strings.Split(encryptedPassword, separator)

	if len(parts) != 3 {
		return "", errors.New("invalid encrypted password")
	}

	// Decode into bytes

	ciphertext, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", err
	}

	salt, err := hex.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	iv, err := hex.DecodeString(parts[2])
	if err != nil {
		return "", err
	}

	// Derive the secret key from the passphrase and salt
	key, _ := deriveKey(passphrase, salt)

	// Create an AES block cipher
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Wrap the block cipher in GCM mode
	gcmCipher, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}

	// Decrypt the password
	password, err := gcmCipher.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func DigestPassphrase(passphrase string) string {
	// Derive a key from the passphrase with a randomly generated salt
	key, salt := deriveKey(passphrase, nil)

	// The digested passphrase is composed of the key and the salt
	digestedPassphrase := strings.Join(
		[]string{hex.EncodeToString(key), hex.EncodeToString(salt)},
		separator,
	)

	return digestedPassphrase
}

func VerifyPassphrase(passphrase string, digestedPassphrase string) bool {
	// Separate the salt and key
	parts := strings.Split(digestedPassphrase, separator)

	if len(parts) != 2 {
		return false
	}

	// Get the ground truth groundTruthKey
	groundTruthKey := parts[0]

	// Get the salt
	salt, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}

	// Derive a key from the input passphrase and the salt from the ground truth
	key, _ := deriveKey(passphrase, salt)

	// Compare with the ground truth key
	return hex.EncodeToString(key) == groundTruthKey
}

func deriveKey(passphrase string, salt []byte) ([]byte, []byte) {
	// Gernate a random salt if it is not provided
	if salt == nil {
		// Allocate space for the salt
		salt = make([]byte, saltLength)

		// Generate the salt randomly
		rand.Read(salt)
	}

	// Derive the secret key from the passphrase and salt
	key := pbkdf2.Key([]byte(passphrase), salt, 100000, secretKeyLength, sha256.New)

	return key, salt
}
