package oauth

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/pkg/errors"
)

// GenerateOAuthState TODO: maybe do something harder
func GenerateOAuthState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "failed to read random bytes")
	}
	state := base64.URLEncoding.EncodeToString(b)

	return state, nil
}

func VerifyOAuthState(state string, original string) bool {
	return state == original
}
