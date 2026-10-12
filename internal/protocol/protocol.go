package protocol

import "encoding/json"

const Version = 1

type Message struct {
	V       int             `json:"v"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

func New(msgType string, payload any) (Message, error) {
	m := Message{V: Version, Type: msgType}
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return Message{}, err
		}
		m.Payload = b
	}
	return m, nil
}
