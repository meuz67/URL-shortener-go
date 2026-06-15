package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

type Encryptor struct {
	gcm cipher.AEAD
}

func NewEncryptor(secret string) (*Encryptor, error) {
	key := sha256.Sum256([]byte(secret))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}

	return &Encryptor{gcm: gcm}, nil
}

func (e *Encryptor) Encrypt(value string) ([]byte, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := e.gcm.Seal(nil, nonce, []byte(value), nil)
	return append(nonce, ciphertext...), nil
}

func (e *Encryptor) Decrypt(value []byte) (string, error) {
	nonceSize := e.gcm.NonceSize()
	if len(value) < nonceSize {
		return "", fmt.Errorf("encrypted value is too short")
	}

	nonce := value[:nonceSize]
	ciphertext := value[nonceSize:]

	plaintext, err := e.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt value: %w", err)
	}

	return string(plaintext), nil
}
