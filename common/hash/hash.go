package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(v ...string) string {
	h := sha256.New()
	for _, s := range v {
		h.Write([]byte(s))
	}
	return hex.EncodeToString(h.Sum(nil))
}
