package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"time"
)

const ScopeAuthentication = "authentication"

type Token struct {
	PlainText string    `json:"token"`
	UserID    int64     `json:"-"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: int64(userID),
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base64.URLEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]
	return token, nil
}

func (m *DBModel) InsertToken(t *Token, u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `
		DELETE FROM tokens WHERE user_id = $1;
	`
	_, err := m.DB.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	stmt = `
		INSERT INTO tokens (user_id, name, email, token_hash, expiry, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = m.DB.ExecContext(ctx, stmt, u.ID, u.LastName, u.Email, t.Hash, t.Expiry, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) GetUserForToken(token string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(token))

	var user User
	stmt := `
		SELECT u.id, u.first_name, u.last_name, u.email
		FROM USERS u
		INNER JOIN tokens t ON t.user_id = u.id
		WHERE t.token_hash = $1 AND t.expiry > $2;
	`

	err := m.DB.QueryRowContext(ctx, stmt, tokenHash[:], time.Now()).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
