package security

import "golang.org/x/crypto/bcrypt"

type SecurityAdapter struct {
}

func NewSecurityAdapter() *SecurityAdapter {
	return &SecurityAdapter{}
}

func (a *SecurityAdapter) Hash(plain string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func (a *SecurityAdapter) Compare(hash, plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err == nil {
		return true, nil
	}
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	return false, err
}
