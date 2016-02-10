package keyman

import (
	"crypto/aes"
	"crypto/sha256"
)

type Key struct {
	bytes       []byte
	fingerprint [32]byte // The SHA-256 fingerprint of this key.
}

func NewKey(data []byte) (*Key, error) {
	switch len(data) {
	case 16, 24, 32:
		break
	default:
		return nil, aes.KeySizeError(len(data))
	}

	return &Key{
		bytes:       data,
		fingerprint: sha256.Sum256(data),
	}, nil
}

func (key Key) GetBytes() []byte {
	return key.bytes
}

func (key Key) GetFingerprint() [32]byte {
	return key.fingerprint
}
