package refreshtokens

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type hasher struct {
	secret []byte
}

func NewHasher(secret string) *hasher {
	return &hasher{
		secret: []byte(secret),
	}
}

func (h *hasher) hash(token string) string {
	mac := hmac.New(sha256.New, h.secret)
	mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}
