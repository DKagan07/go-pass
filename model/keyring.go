package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/pbkdf2"
)

const (
	NONCE_SIZE     = 12
	KEY_SIZE       = 32
	NUM_ITERATIONS = 100000

	DefaultKeyringService = "gopass"
	DefaultKeyringAccount = "encryption_key"
)

// MasterAESKeyManager is the struct that contains the logic to handle the
// keyring encryption and decryption
type MasterAESKeyManager struct {
	Masterpassword string
	// KeyringService and KeyringAccount allow tests to use isolated keyring entries
	KeyringService string
	KeyringAccount string
}

// NewMasterAESKeyManager returns a new MasterAESKeyManager with the passed-in
// password
func NewMasterAESKeyManager(mp string) *MasterAESKeyManager {
	return &MasterAESKeyManager{
		Masterpassword: mp,
		KeyringService: DefaultKeyringService,
		KeyringAccount: DefaultKeyringAccount,
	}
}

// NewTestMasterAESKeyManager creates a keyring manager for testing with isolated keyring entries
func NewTestMasterAESKeyManager(mp string) *MasterAESKeyManager {
	return &MasterAESKeyManager{
		Masterpassword: mp,
		KeyringService: "gopass-test",
		KeyringAccount: "test_encryption_key",
	}
}

// InitializeKeychain creates a new keyring and sets it
func (k *MasterAESKeyManager) InitializeKeychain() error {
	randomKey := make([]byte, 32)
	if _, err := rand.Read(randomKey); err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(randomKey)

	if err := keyring.Set(k.KeyringService, k.KeyringAccount, encoded); err != nil {
		return err
	}

	return nil
}

// DeleteKeychain removes the keyring entry (useful for tests and cleanup)
func (k *MasterAESKeyManager) DeleteKeychain() error {
	return keyring.Delete(k.KeyringService, k.KeyringAccount)
}

// GetEncryptionKey returns the encrypted key that is stored in the keyring
func (k *MasterAESKeyManager) GetEncryptionKey() ([]byte, error) {
	salt := GetSalt()

	encoded, err := keyring.Get(k.KeyringService, k.KeyringAccount)
	if err != nil {
		return nil, err
	}

	baseKey, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	key := pbkdf2.Key(
		append(baseKey, []byte(k.Masterpassword)...),
		salt,
		NUM_ITERATIONS,
		KEY_SIZE,
		sha256.New,
	)

	return key, nil
}

// Encrypt encrypts the []byte using the keyring, and returns the hex-encoded
// representation of the encrypted text
func (k *MasterAESKeyManager) Encrypt(plaintext []byte) (string, error) {
	key, err := k.GetEncryptionKey()
	if err != nil {
		return "", err
	}

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("creating cipher block: %v", err)
	}

	aesgcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", fmt.Errorf("creating aes gcm: %v", err)
	}

	nonce, err := GenerateNonce()
	if err != nil {
		return "", err
	}

	cipherText := aesgcm.Seal(nil, nonce, plaintext, nil)

	// Appending the nonce to the data bytes
	cipherText = append(nonce, cipherText...)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts the passed-in ciphertext using the keyring, and returns the
// []byte represnetatino of the text. This []byte representation can be used in
// conjunctino with `string()` to have it be in string form
func (k *MasterAESKeyManager) Decrypt(ciphertext string) ([]byte, error) {
	key, err := k.GetEncryptionKey()
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}

	nonce, cipher := decoded[:NONCE_SIZE], decoded[NONCE_SIZE:]
	plaintext, err := aesgcm.Open(nil, nonce, cipher, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// GenerateNonce generates a Number Once, used for AES-256 encryption.
func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, NONCE_SIZE)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("creating nonce: %v", err)
	}

	if len(nonce) != NONCE_SIZE {
		return nil, fmt.Errorf("nonce not correct length")
	}
	return nonce, nil
}

// GetSalt is a helper function that gets the key from the environment
// variable called 'SECRET_PASSWORD_KEY'. 'SECRET_PASSWORD_KEY' should be a
// strong, 32-byte password that's unique. You can generate a strong password
// from something like: 'https://passwords-generator.org/32-character'
func GetSalt() []byte {
	key := []byte(os.Getenv("SECRET_PASSWORD_KEY"))
	if len(key) != KEY_SIZE {
		log.Fatal("Salt not appropriate length or not present. Please set.")
	}
	return key
}
