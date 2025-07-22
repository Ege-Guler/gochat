package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

type SecureComm struct {
	PrivateKey *ecdh.PrivateKey
	PublicKey  []byte
}

func GenerateECDH() (*SecureComm, error) {
	curve := ecdh.P256()
	priv, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ECDH key: %w", err)
	}
	pubBytes := priv.PublicKey().Bytes()

	return &SecureComm{
		PrivateKey: priv,
		PublicKey:  pubBytes,
	}, nil
}

func ComputeKEK(ownPrivateKey *ecdh.PrivateKey, peerPublicKey []byte) ([]byte, error) {
	peerPubKey, err := ecdh.P256().NewPublicKey(peerPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to validate peer public key: %w", err)
	}

	sharedSecret, err := ownPrivateKey.ECDH(peerPubKey)
	if err != nil {
		return nil, fmt.Errorf("ECDH computation failed: %w", err)
	}

	hash := sha256.Sum256(sharedSecret)
	return hash[:], nil
}

func GenerateAES() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %w", err)
	}
	return key, nil
}

func EncryptAES(key []byte, plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return nonce, ciphertext, nil
}

func DecryptAES(key, nonce, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}
	return plaintext, nil
}
