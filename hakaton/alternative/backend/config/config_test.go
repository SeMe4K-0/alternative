package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DB_HOST", "test_db_host")
	os.Setenv("DB_PORT", "1234")
	os.Setenv("DB_USER", "test_db_user")
	os.Setenv("DB_PASSWORD", "test_db_password")
	os.Setenv("DB_NAME", "test_db_name")
	os.Setenv("DB_SSLMODE", "require")

	os.Setenv("REDIS_HOST", "test_redis_host")
	os.Setenv("REDIS_PORT", "5678")
	os.Setenv("REDIS_PASSWORD", "test_redis_password")

	os.Setenv("MINIO_ENDPOINT", "test_minio_endpoint")
	os.Setenv("MINIO_ACCESS_KEY", "test_minio_access_key")
	os.Setenv("MINIO_SECRET_KEY", "test_minio_secret_key")
	os.Setenv("MINIO_USE_SSL", "true")
	os.Setenv("MINIO_BUCKET", "test_minio_bucket")

	os.Setenv("SERVER_PORT", ":9000")
	os.Setenv("PYTHON_SERVICE_URL", "http://test_python_service:5000")

	cfg := Load()

	if cfg.Database.Host != "test_db_host" {
		t.Errorf("Expected DB_HOST 'test_db_host', got %s", cfg.Database.Host)
	}
	if cfg.Database.Port != 1234 {
		t.Errorf("Expected DB_PORT 1234, got %d", cfg.Database.Port)
	}
	if cfg.Database.User != "test_db_user" {
		t.Errorf("Expected DB_USER 'test_db_user', got %s", cfg.Database.User)
	}
	if cfg.Database.Password != "test_db_password" {
		t.Errorf("Expected DB_PASSWORD 'test_db_password', got %s", cfg.Database.Password)
	}
	if cfg.Database.DBName != "test_db_name" {
		t.Errorf("Expected DB_NAME 'test_db_name', got %s", cfg.Database.DBName)
	}
	if cfg.Database.SSLMode != "require" {
		t.Errorf("Expected DB_SSLMODE 'require', got %s", cfg.Database.SSLMode)
	}

	if cfg.Redis.Host != "test_redis_host" {
		t.Errorf("Expected REDIS_HOST 'test_redis_host', got %s", cfg.Redis.Host)
	}
	if cfg.Redis.Port != 5678 {
		t.Errorf("Expected REDIS_PORT 5678, got %d", cfg.Redis.Port)
	}
	if cfg.Redis.Password != "test_redis_password" {
		t.Errorf("Expected REDIS_PASSWORD 'test_redis_password', got %s", cfg.Redis.Password)
	}

	if cfg.MinIO.Endpoint != "test_minio_endpoint" {
		t.Errorf("Expected MINIO_ENDPOINT 'test_minio_endpoint', got %s", cfg.MinIO.Endpoint)
	}
	if cfg.MinIO.AccessKeyID != "test_minio_access_key" {
		t.Errorf("Expected MINIO_ACCESS_KEY 'test_minio_access_key', got %s", cfg.MinIO.AccessKeyID)
	}
	if cfg.MinIO.SecretAccessKey != "test_minio_secret_key" {
		t.Errorf("Expected MINIO_SECRET_KEY 'test_minio_secret_key', got %s", cfg.MinIO.SecretAccessKey)
	}
	if !cfg.MinIO.UseSSL {
		t.Error("Expected MINIO_USE_SSL true, got false")
	}
	if cfg.MinIO.BucketName != "test_minio_bucket" {
		t.Errorf("Expected MINIO_BUCKET 'test_minio_bucket', got %s", cfg.MinIO.BucketName)
	}

	if cfg.Server.Port != ":9000" {
		t.Errorf("Expected SERVER_PORT ':9000', got %s", cfg.Server.Port)
	}

	if cfg.PythonService.URL != "http://test_python_service:5000" {
		t.Errorf("Expected PYTHON_SERVICE_URL 'http://test_python_service:5000', got %s", cfg.PythonService.URL)
	}

	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")
	os.Unsetenv("REDIS_HOST")
	os.Unsetenv("REDIS_PORT")
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("MINIO_ENDPOINT")
	os.Unsetenv("MINIO_ACCESS_KEY")
	os.Unsetenv("MINIO_SECRET_KEY")
	os.Unsetenv("MINIO_USE_SSL")
	os.Unsetenv("MINIO_BUCKET")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("PYTHON_SERVICE_URL")
}

func TestLoadConfig_Defaults(t *testing.T) {
	cfg := Load()

	if cfg.Database.Host == "" {
		t.Error("Database host should have a default value")
	}

	if cfg.Database.Port == 0 {
		t.Error("Database port should have a default value")
	}

	if cfg.Redis.Host == "" {
		t.Error("Redis host should have a default value")
	}

	if cfg.Redis.Port == 0 {
		t.Error("Redis port should have a default value")
	}

	if cfg.MinIO.Endpoint == "" {
		t.Error("MinIO endpoint should have a default value")
	}

	if cfg.Server.Port == "" {
		t.Error("Server port should have a default value")
	}

	if cfg.PythonService.URL == "" {
		t.Error("Python service URL should have a default value")
	}
}

func TestLoadConfig_EnvironmentVariables(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	value := os.Getenv("TEST_VAR")
	if value != "test_value" {
		t.Errorf("Expected TEST_VAR 'test_value', got %s", value)
	}
}

func TestLoadConfig_EmptyEnvironment(t *testing.T) {
	os.Clearenv()

	cfg := Load()

	if cfg.Database.Host == "" {
		t.Error("Database host should have a default value even with empty environment")
	}

	if cfg.Redis.Host == "" {
		t.Error("Redis host should have a default value even with empty environment")
	}

	if cfg.MinIO.Endpoint == "" {
		t.Error("MinIO endpoint should have a default value even with empty environment")
	}
}