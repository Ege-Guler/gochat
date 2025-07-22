package mqtt

import (
	"github.com/Ege-Guler/gochat/internal/crypto"
	"github.com/google/uuid"
)

type Session struct {
	SessionID       string
	OwnKeys         *crypto.SecureComm
	PeerPubKey      []byte
	KEK             []byte
	AESKey          []byte
	EncryptedAESKey []byte
	Nonce           []byte
}

func (s *Session) GenerateSessionID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	s.SessionID = id.String()
	return nil
}

func (s *Session) GenerateOwnKeys() error {
	keys, err := crypto.GenerateECDH()
	if err != nil {
		return err
	}
	s.OwnKeys = keys
	return nil
}

func (s *Session) GenerateKEK() error {
	KEK, err := crypto.ComputeKEK(s.OwnKeys.PrivateKey, s.PeerPubKey)
	if err != nil {
		return err
	}
	s.KEK = KEK
	return nil
}

func (s *Session) SetAES() error {
	aesKey, err := crypto.GenerateAES()
	if err != nil {
		return err
	}

	nonce, encryptedAESKey, err := crypto.EncryptAES(s.KEK, aesKey)
	if err != nil {
		return err
	}

	s.AESKey = aesKey
	s.Nonce = nonce
	s.EncryptedAESKey = encryptedAESKey

	return nil
}

func (s *Session) DecryptPeerAES(peerNonce, peerEncryptedAESKey []byte) ([]byte, error) {
	return crypto.DecryptAES(s.KEK, peerNonce, peerEncryptedAESKey)
}
