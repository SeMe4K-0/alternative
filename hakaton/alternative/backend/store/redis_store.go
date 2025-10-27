package store

import (
	"backend/config"
	"backend/models"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(cfg *config.RedisConfig) (*RedisStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		// Password: cfg.Password,
	})

	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisStore{
		client: rdb,
		ctx:    ctx,
	}, nil
}

func (r *RedisStore) CreateSession(userID uint64) (map[string]interface{}, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	session := map[string]interface{}{
		"user_id":    userID,
		"token":      token,
		"expires_at": expiresAt.Format(time.RFC3339),
	}

	sessionData, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	ttl := 24 * time.Hour
	if err := r.client.Set(r.ctx, token, sessionData, ttl).Err(); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *RedisStore) GetUserBySession(token string) (*models.User, bool) {
	sessionData, err := r.client.Get(r.ctx, token).Result()
	if err != nil {
		return nil, false
	}

	var session map[string]interface{}
	if err := json.Unmarshal([]byte(sessionData), &session); err != nil {
		return nil, false
	}

	expiresAt, ok := session["expires_at"].(string)
	if !ok {
		return nil, false
	}
	
	expiresTime, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return nil, false
	}
	
	if time.Now().After(expiresTime) {
		r.DeleteSession(token)
		return nil, false
	}

	userID, ok := session["user_id"].(float64)
	if !ok {
		return nil, false
	}

	user := &models.User{
		ID: uint64(userID),
	}

	return user, true
}

func (r *RedisStore) DeleteSession(token string) error {
	return r.client.Del(r.ctx, token).Err()
}

func (r *RedisStore) CleanExpiredSessions() error {
	return nil
}

func (r *RedisStore) Close() error {
	return r.client.Close()
}

func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
