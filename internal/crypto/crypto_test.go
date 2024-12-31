package crypto_test

import (
	"fmt"
	"testing"

	mycrypto "github.com/Isaac-Fate/myst/internal/crypto"
)

const passphrase string = "hello, world"
const password string = "password123456!"

func TestEncryptDecrypt(t *testing.T) {
	encryptedPassword, err := mycrypto.Encrypt(passphrase, password)

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("encrypted password: %s\n", encryptedPassword)

	recoveredPassword, err := mycrypto.Decrypt(passphrase, encryptedPassword)

	if err != nil {
		t.Error(err)
	}

	if recoveredPassword != password {
		t.Errorf("expected %s, got %s", password, recoveredPassword)
	}

	fmt.Printf("recovered password: %s\n", recoveredPassword)

}

func TestPassphrase(t *testing.T) {
	digestedPassphrase := mycrypto.DigestPassphrase(passphrase)

	fmt.Println("digested passphrase: ", digestedPassphrase)

	isPassphraseCorrect := mycrypto.VerifyPassphrase(passphrase, digestedPassphrase)
	if !isPassphraseCorrect {
		t.Errorf("expected true, got false")
	}
}
