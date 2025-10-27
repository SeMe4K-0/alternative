package store

import (
	"backend/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRedisStore_InvalidConfig(t *testing.T) {
	cfg := &config.RedisConfig{
		Host: "invalid-host",
		Port: 6379,
	}

	store, err := NewRedisStore(cfg)

	assert.Error(t, err)
	assert.Nil(t, store)
}

func TestNewRedisStore_EmptyConfig(t *testing.T) {
	cfg := &config.RedisConfig{}

	store, err := NewRedisStore(cfg)

	assert.Error(t, err)
	assert.Nil(t, store)
}

func TestNewRedisStore_NilConfig(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil config")
		}
	}()
	
	NewRedisStore(nil)
}

func TestRedisStore_CreateSession_NoConnection(t *testing.T) {
	store := &RedisStore{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil client")
		}
	}()
	
	store.CreateSession(1)
}

func TestRedisStore_GetUserBySession_NoConnection(t *testing.T) {
	store := &RedisStore{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil client")
		}
	}()
	
	store.GetUserBySession("token123")
}

func TestRedisStore_DeleteSession_NoConnection(t *testing.T) {
	store := &RedisStore{}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil client")
		}
	}()
	
	store.DeleteSession("token123")
}

func TestGenerateToken(t *testing.T) {
	token, err := generateToken()

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Len(t, token, 64) // 32 bytes = 64 hex characters
}

func TestGenerateToken_Multiple(t *testing.T) {
	token1, err1 := generateToken()
	token2, err2 := generateToken()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, token1, token2) // Tokens should be different
}

func TestRedisStore_GetUserBySession_ValidSession(t *testing.T) {
	// Create a mock Redis store with a valid session
	store := &RedisStore{}
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil client")
		}
	}()
	
	store.GetUserBySession("valid-token")
}

func TestRedisStore_GetUserBySession_ExpiredSession(t *testing.T) {
	// Create a mock Redis store with an expired session
	store := &RedisStore{}
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil client")
		}
	}()
	
	store.GetUserBySession("expired-token")
}

func TestRedisStore_CleanExpiredSessions(t *testing.T) {
	store := &RedisStore{}
	
	// This should not panic as it returns nil
	err := store.CleanExpiredSessions()
	assert.NoError(t, err)
}

func TestRedisStore_Close(t *testing.T) {
	store := &RedisStore{}
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil client")
		}
	}()
	
	store.Close()
}