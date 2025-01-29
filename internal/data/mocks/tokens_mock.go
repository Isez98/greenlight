package mocks

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"greenlight.isez.dev/internal/data"
)

type TokenModel_Mock struct {
	DB *sql.DB
}

func generateToken(userID int64, ttl time.Duration, scope string) (*data.Token, error) {
	token := &data.Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func (m TokenModel_Mock) New(userID int64, ttl time.Duration, scope string) (*data.Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, nil
	}

	err = m.Insert(token)
	return token, nil
}

func (m TokenModel_Mock) Insert(token *data.Token) error {
	return nil
}

func (m TokenModel_Mock) DeleteAllForUser(scope string, userID int64) error {
	return nil
}

func (m TokenModel_Mock) ValidateToken(tokenString string) (bool, error) {
	return false, nil

}
