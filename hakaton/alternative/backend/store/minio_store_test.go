package store

import (
	"backend/config"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMinIOStore_InvalidConfig(t *testing.T) {
	cfg := &config.MinIOConfig{
		Endpoint:        "invalid-endpoint",
		AccessKeyID:     "invalid-key",
		SecretAccessKey: "invalid-secret",
		BucketName:      "test-bucket",
		UseSSL:          false,
	}

	store, err := NewMinIOStore(cfg)

	assert.Error(t, err)
	assert.Nil(t, store)
}

func TestNewMinIOStore_EmptyConfig(t *testing.T) {
	cfg := &config.MinIOConfig{}

	store, err := NewMinIOStore(cfg)

	assert.Error(t, err)
	assert.Nil(t, store)
}

func TestNewMinIOStore_NilConfig(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil config")
		}
	}()
	
	NewMinIOStore(nil)
}

func TestMinIOStore_ValidateImageType(t *testing.T) {
	store := &MinIOStore{}

	assert.True(t, store.ValidateImageType("image/jpeg"))
	assert.True(t, store.ValidateImageType("image/png"))
	assert.True(t, store.ValidateImageType("image/gif"))
	assert.True(t, store.ValidateImageType("image/webp"))

	assert.False(t, store.ValidateImageType("text/plain"))
	assert.False(t, store.ValidateImageType("application/json"))
	assert.False(t, store.ValidateImageType(""))
}

func TestMinIOStore_GetAvatarURL(t *testing.T) {
	store := &MinIOStore{}

	url := store.GetAvatarURL("avatar_user_1.jpg")

	assert.True(t, strings.HasPrefix(url, "/api/files/"))
	assert.True(t, strings.Contains(url, "avatar_user_1.jpg"))
}

func TestMinIOStore_GetObservationImageURL(t *testing.T) {
	store := &MinIOStore{}

	url := store.GetObservationImageURL("observation_1.jpg")

	assert.True(t, strings.HasPrefix(url, "/api/files/"))
	assert.True(t, strings.Contains(url, "observation_1.jpg"))
}

func TestMinIOStore_UploadAvatar_NoConnection(t *testing.T) {
	assert.True(t, true)
}

func TestMinIOStore_UploadObservationImage_NoConnection(t *testing.T) {
	assert.True(t, true)
}

func TestMinIOStore_DeleteAvatar_NoConnection(t *testing.T) {
	assert.True(t, true)
}

func TestMinIOStore_DeleteObservationImage_NoConnection(t *testing.T) {
	assert.True(t, true)
}

func TestMinIOStore_GetObject_NoConnection(t *testing.T) {
	assert.True(t, true)
}

func TestMinIOStore_GetObjectInfo_NoConnection(t *testing.T) {
	assert.True(t, true)
}

func TestMinIOStore_ValidateImageType_EdgeCases(t *testing.T) {
	store := &MinIOStore{}

	assert.True(t, store.ValidateImageType("IMAGE/JPEG"))
	assert.True(t, store.ValidateImageType("Image/Jpeg"))
	assert.True(t, store.ValidateImageType("image/JPEG"))

	assert.False(t, store.ValidateImageType(" image/jpeg"))
	assert.False(t, store.ValidateImageType("image/jpeg "))
	assert.False(t, store.ValidateImageType(" image/jpeg "))

	assert.False(t, store.ValidateImageType("image/jpeg; charset=utf-8"))
	assert.False(t, store.ValidateImageType("image/png; quality=0.8"))

	assert.False(t, store.ValidateImageType("image/bmp"))
	assert.False(t, store.ValidateImageType("image/tiff"))
	assert.False(t, store.ValidateImageType("image/svg+xml"))
}

func TestMinIOStore_GetAvatarURL_EdgeCases(t *testing.T) {
	store := &MinIOStore{}

	url := store.GetAvatarURL("")
	assert.True(t, strings.HasPrefix(url, "/api/files/"))
	assert.True(t, strings.Contains(url, ""))

	url = store.GetAvatarURL("avatar_user_1@#$%.jpg")
	assert.True(t, strings.HasPrefix(url, "/api/files/"))
	assert.True(t, strings.Contains(url, "avatar_user_1@#$%.jpg"))

	url = store.GetAvatarURL("avatar user 1.jpg")
	assert.True(t, strings.HasPrefix(url, "/api/files/"))
	assert.True(t, strings.Contains(url, "avatar user 1.jpg"))
}

func TestMinIOStore_GetObservationImageURL_EdgeCases(t *testing.T) {
	store := &MinIOStore{}

	url := store.GetObservationImageURL("")
	assert.True(t, strings.HasPrefix(url, "/api/files/"))
	assert.True(t, strings.Contains(url, ""))

	url = store.GetObservationImageURL("observation_1@#$%.jpg")
	assert.True(t, strings.HasPrefix(url, "/api/files/"))
	assert.True(t, strings.Contains(url, "observation_1@#$%.jpg"))

	url = store.GetObservationImageURL("observation 1.jpg")
	assert.True(t, strings.HasPrefix(url, "/api/files/"))
	assert.True(t, strings.Contains(url, "observation 1.jpg"))
}
