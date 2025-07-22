package mqtt

import (
	"encoding/json"
)

type ExchangeMessage struct {
	SessionID string `json:"session_id"`
	PublicKey []byte `json:"public_key"`
	Nonce     []byte `json:"nonce"`
}

func (s *Session) ToExchangeMessage() ExchangeMessage {
	return ExchangeMessage{
		SessionID: s.SessionID,
		PublicKey: s.OwnKeys.PublicKey,
		Nonce:     s.Nonce,
	}
}

func (e *ExchangeMessage) EncodeExchangeMessage() ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *ExchangeMessage) DecodeExchangeMessage(data []byte) error {
	if err := json.Unmarshal(data, e); err != nil {
		return err
	}
	return nil
}
