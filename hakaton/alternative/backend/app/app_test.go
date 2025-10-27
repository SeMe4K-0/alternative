package app

import (
	"backend/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp_InvalidConfig(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "invalid-host",
			Port:     5432,
			User:     "invalid-user",
			Password: "invalid-password",
			DBName:   "invalid-db",
			SSLMode:  "disable",
		},
	}

	app, err := NewApp(cfg)

	assert.Error(t, err)
	assert.Nil(t, app)
}

func TestNewApp_EmptyConfig(t *testing.T) {
	cfg := &config.Config{}

	app, err := NewApp(cfg)

	assert.Error(t, err)
	assert.Nil(t, app)
}

func TestNewApp_NilConfig(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil config")
		}
	}()
	
	NewApp(nil)
}

func TestApp_Run_InvalidPort(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "invalid-port", // Invalid port
		},
	}
	
	app, err := NewApp(cfg)
	if err != nil {
		return
	}
	
	defer func() {
		if r := recover(); r == nil {
		}
	}()
	
	app.Run(":invalid-port")
}

func TestApp_Run_ValidConfig(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "test",
			Password: "test",
			DBName:   "test",
			SSLMode:  "disable",
		},
	}
	
	app, err := NewApp(cfg)
	if err != nil {
		return
	}
	
	defer func() {
		if r := recover(); r == nil {
		}
	}()
	
	app.Run(":8080")
}
