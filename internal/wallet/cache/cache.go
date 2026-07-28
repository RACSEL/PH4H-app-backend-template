package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type entry struct {
	raw       string
	expiresAt time.Time
}

var (
	mu    sync.RWMutex
	store = make(map[string]entry)
	ttl   = 24 * time.Hour
)

func keyFromClaims(claims map[string]interface{}) (string, error) {
	data, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}

func Set(claims map[string]interface{}, raw string) {
	k, err := keyFromClaims(claims)
	if err != nil {
		return
	}
	mu.Lock()
	store[k] = entry{raw: raw, expiresAt: time.Now().Add(ttl)}
	mu.Unlock()
}

func Get(claims map[string]interface{}) (string, bool) {
	deleteExpired()

	k, err := keyFromClaims(claims)
	if err != nil {
		return "", false
	}
	mu.RLock()
	e, ok := store[k]
	mu.RUnlock()
	if !ok || time.Now().After(e.expiresAt) {
		return "", false
	}
	return e.raw, true
}

func deleteExpired() {
	mu.Lock()
	defer mu.Unlock()
	now := time.Now()
	for k, e := range store {
		if now.After(e.expiresAt) {
			delete(store, k)
		}
	}
}
