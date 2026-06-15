package storage

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"
)

var ErrURLNotFound = errors.New("url not found")

type URL struct {
	Code      string
	LongURL   string
	CreatedAt time.Time
}

type URLStore interface {
	Save(ctx context.Context, longURL string) (URL, error)
	FindByCode(ctx context.Context, code string) (URL, error)
}

type MemoryURLStore struct {
	mu      sync.RWMutex
	counter uint64
	byCode  map[string]URL
}

func NewMemoryURLStore() *MemoryURLStore {
	return &MemoryURLStore{
		byCode: make(map[string]URL),
	}
}

func (s *MemoryURLStore) Save(_ context.Context, longURL string) (URL, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	code := strconv.FormatUint(s.counter, 36)

	item := URL{
		Code:      code,
		LongURL:   longURL,
		CreatedAt: time.Now().UTC(),
	}

	s.byCode[code] = item

	return item, nil
}

func (s *MemoryURLStore) FindByCode(_ context.Context, code string) (URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.byCode[code]
	if !ok {
		return URL{}, ErrURLNotFound
	}

	return item, nil
}
