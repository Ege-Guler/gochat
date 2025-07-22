package mqtt

import (
	"encoding/json"
)

const (
	IM = "init_message"
	KM = "key_message"
	CM = "connection_message"
)

type generic struct {
	Type string `json:"type"`
}

type InitExchangeMessage struct {
	Type      string `json:"type"`
	SessionID string `json:"session_id"`
	PublicKey []byte `json:"public_key"`
}

type KeyExchangeMessage struct {
	Type            string `json:"type"`
	SessionID       string `json:"session_id"`
	EncryptedAESKey []byte `json:"encrypted_aes_key"`
	Nonce           []byte `json:"nonce"`
}

type ExchangeMessage struct {
	Type            string `json:"type"`
	SessionID       string `json:"session_id"`
	PublicKey       []byte `json:"public_key"`
	EncryptedAESKey []byte `json:"encrypted_aes_key"`
	Nonce           []byte `json:"nonce"`
}

func (s *Session) ToInitExchange() InitExchangeMessage {
	return InitExchangeMessage{
		Type:      IM,
		SessionID: s.SessionID,
		PublicKey: s.OwnKeys.PublicKey,
	}
}

func (s *Session) ToKeyExchangeMessage() KeyExchangeMessage {
	return KeyExchangeMessage{
		Type:            KM,
		SessionID:       s.SessionID,
		EncryptedAESKey: s.EncryptedAESKey,
		Nonce:           s.Nonce,
	}
}

func (i *InitExchangeMessage) Encode() ([]byte, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (i *InitExchangeMessage) Decode(data []byte) error {
	if err := json.Unmarshal(data, i); err != nil {
		return err
	}
	return nil
}

func (k *KeyExchangeMessage) Encode() ([]byte, error) {
	data, err := json.Marshal(k)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (k *KeyExchangeMessage) Decode(data []byte) error {
	if err := json.Unmarshal(data, k); err != nil {
		return err
	}
	return nil
}

func InitExchangeFunc(s *Session) ([]byte, error) {
	im := s.ToInitExchange()
	payload, err := im.Encode()
	if err != nil {
		return nil, err
	}
	return payload, nil

}

// func (s *Session) handler(em *ExchangeMessage) error {

// 	if em.Type == EM {
// 		s.PeerPubKey = em.PublicKey

// 		// generates key exchange key
// 		if err := s.GenerateKEK(); err != nil {
// 			return err
// 		}

// 		// sets aeskey, nonce and encrypted AES key
// 		if err := s.SetAES(); err != nil {
// 			return err
// 		}

// 	}

// 	return nil
// }

// func (e *ExchangeMessage) EncodeExchangeMessage() ([]byte, error) {
// 	data, err := json.Marshal(e)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }

// func (e *ExchangeMessage) DecodeExchangeMessage(data []byte) error {
// 	if err := json.Unmarshal(data, e); err != nil {
// 		return err
// 	}
// 	return nil
// }
