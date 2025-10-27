//go:build integration
// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestFullIntegration(t *testing.T) {
	baseURL := "http://localhost:8080"
	
	// Step 1: Register a test user
	fmt.Println("Step 1: Registering test user...")
	
	registerReq := map[string]string{
		"email":    "test@example.com",
		"password": "testpassword123",
		"username": "testuser",
	}
	
	jsonData, _ := json.Marshal(registerReq)
	resp, err := http.Post(baseURL+"/api/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Register Status: %d\n", resp.StatusCode)
	fmt.Printf("Register Response: %s\n", string(body))
	
	// Step 2: Request password reset
	fmt.Println("\nStep 2: Requesting password reset...")
	
	resetReq := map[string]string{
		"email": "test@example.com",
	}
	
	jsonData, _ = json.Marshal(resetReq)
	resp, err = http.Post(baseURL+"/api/auth/forgot-password", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Reset Request Status: %d\n", resp.StatusCode)
	fmt.Printf("Reset Request Response: %s\n", string(body))
	
	// Step 3: Try to login with old password (should fail after reset)
	fmt.Println("\nStep 3: Testing login with old password...")
	
	loginReq := map[string]string{
		"email":    "test@example.com",
		"password": "testpassword123",
	}
	
	jsonData, _ = json.Marshal(loginReq)
	resp, err = http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Login Status: %d\n", resp.StatusCode)
	fmt.Printf("Login Response: %s\n", string(body))
	
	fmt.Println("\nFull test completed!")
}
