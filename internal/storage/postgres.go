package storage

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	appcrypto "url-shortener/internal/crypto"
)

const (
	codeAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	codeLength   = 8
)

type PostgresURLStore struct {
	pool      *pgxpool.Pool
	encryptor *appcrypto.Encryptor
}

func NewPostgresURLStore(pool *pgxpool.Pool, encryptor *appcrypto.Encryptor) *PostgresURLStore {
	return &PostgresURLStore{
		pool:      pool,
		encryptor: encryptor,
	}
}

func (s *PostgresURLStore) Save(ctx context.Context, longURL string) (URL, error) {
	encryptedURL, err := s.encryptor.Encrypt(longURL)
	if err != nil {
		return URL{}, fmt.Errorf("encrypt url: %w", err)
	}

	for attempts := 0; attempts < 5; attempts++ {
		code, err := randomCode(codeLength)
		if err != nil {
			return URL{}, fmt.Errorf("generate code: %w", err)
		}

		var createdAt time.Time
		err = s.pool.QueryRow(ctx, `
			INSERT INTO urls(code, encrypted_url)
			VALUES($1, $2)
			RETURNING created_at
		`, code, encryptedURL).Scan(&createdAt)
		if err == nil {
			return URL{
				Code:      code,
				LongURL:   longURL,
				CreatedAt: createdAt,
			}, nil
		}

		if isUniqueViolation(err) {
			continue
		}

		return URL{}, fmt.Errorf("insert url: %w", err)
	}

	return URL{}, errors.New("failed to generate unique short code")
}

func (s *PostgresURLStore) FindByCode(ctx context.Context, code string) (URL, error) {
	var encryptedURL []byte
	var createdAt time.Time

	err := s.pool.QueryRow(ctx, `
		SELECT encrypted_url, created_at
		FROM urls
		WHERE code = $1
	`, code).Scan(&encryptedURL, &createdAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return URL{}, ErrURLNotFound
	}
	if err != nil {
		return URL{}, fmt.Errorf("select url: %w", err)
	}

	longURL, err := s.encryptor.Decrypt(encryptedURL)
	if err != nil {
		return URL{}, fmt.Errorf("decrypt url: %w", err)
	}

	return URL{
		Code:      code,
		LongURL:   longURL,
		CreatedAt: createdAt,
	}, nil
}

func randomCode(length int) (string, error) {
	result := make([]byte, length)
	max := big.NewInt(int64(len(codeAlphabet)))

	for i := range result {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}

		result[i] = codeAlphabet[n.Int64()]
	}

	return string(result), nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
