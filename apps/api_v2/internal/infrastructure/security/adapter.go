package security

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"os"
)

var requiredSymbols = []byte("!@#$_-")

type SecurityAdapter struct {
	pepper string
}

func NewSecurityAdapter() *SecurityAdapter {
	return &SecurityAdapter{
		pepper: os.Getenv("PASSWORD_HASH_PEPPER"),
	}
}

func (a *SecurityAdapter) Hash(plain string) (string, error) {
	// Deterministic hash for integrations that compare exact values.
	sum := sha256.Sum256([]byte(a.pepper + ":" + plain))

	base := hex.EncodeToString(sum[:])
	prefix := []byte{
		byte('A' + (sum[0] % 26)),
		byte('a' + (sum[1] % 26)),
		byte('0' + (sum[2] % 10)),
		requiredSymbols[sum[3]%byte(len(requiredSymbols))],
	}

	return string(prefix) + base, nil
}

func (a *SecurityAdapter) Compare(hash, plain string) (bool, error) {
	candidate, err := a.Hash(plain)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare([]byte(hash), []byte(candidate)) == 1, nil
}
