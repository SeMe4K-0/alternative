package models

import (
	"testing"
)

func TestUser_HashPassword(t *testing.T) {
	user := &User{}
	password := "testpassword"

	err := user.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if user.Password == password {
		t.Error("Password should be hashed")
	}

	if user.Password == "" {
		t.Error("Hashed password should not be empty")
	}

	if len(user.Password) < 10 {
		t.Error("Hashed password should be longer than 10 characters")
	}
}

func TestUser_HashPassword_Empty(t *testing.T) {
	user := &User{}
	password := ""

	err := user.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword with empty password failed: %v", err)
	}

	if user.Password == "" {
		t.Error("Hashed password should not be empty even for empty input")
	}
}

func TestUser_CheckPassword(t *testing.T) {
	user := &User{}
	password := "testpassword"
	user.HashPassword(password)

	if !user.CheckPassword(password) {
		t.Error("CheckPassword failed for correct password")
	}

	if user.CheckPassword("wrongpassword") {
		t.Error("CheckPassword succeeded for wrong password")
	}

	if user.CheckPassword("") {
		t.Error("CheckPassword succeeded for empty password")
	}
}

func TestUser_CheckPassword_Empty(t *testing.T) {
	user := &User{}
	user.HashPassword("")

	if user.CheckPassword("testpassword") {
		t.Error("CheckPassword succeeded for wrong password with empty hashed password")
	}

	if !user.CheckPassword("") {
		t.Error("CheckPassword failed for correct empty password")
	}
}

func TestUser_CheckPassword_Invalid(t *testing.T) {
	user := &User{}
	user.Password = "invalid-hash"

	if user.CheckPassword("anypassword") {
		t.Error("CheckPassword succeeded for invalid hash")
	}
}

func TestUser_PasswordStrength(t *testing.T) {
	passwords := []string{
		"password123",
		"secretpass",
		"mypassword",
	}

	for _, password := range passwords {
		user := &User{}
		err := user.HashPassword(password)
		if err != nil {
			t.Fatalf("HashPassword failed for %s: %v", password, err)
		}

		if !user.CheckPassword(password) {
			t.Errorf("CheckPassword failed for password %s", password)
		}
	}
}

func TestUser_PasswordUniqueness(t *testing.T) {
	password := "testpassword"
	user1 := &User{}
	user2 := &User{}

	user1.HashPassword(password)
	user2.HashPassword(password)

	if user1.Password == user2.Password {
		t.Error("Hashed passwords should be unique even for same input")
	}
}

func TestUser_PasswordLength(t *testing.T) {
	shortPassword := "123"
	longPassword := "thisisaverylongpasswordthatshouldworkfine"

	user1 := &User{}
	user2 := &User{}

	err := user1.HashPassword(shortPassword)
	if err != nil {
		t.Fatalf("HashPassword failed for short password: %v", err)
	}

	err = user2.HashPassword(longPassword)
	if err != nil {
		t.Fatalf("HashPassword failed for long password: %v", err)
	}

	if !user1.CheckPassword(shortPassword) {
		t.Error("CheckPassword failed for short password")
	}

	if !user2.CheckPassword(longPassword) {
		t.Error("CheckPassword failed for long password")
	}
}

func TestUser_PasswordSpecialChars(t *testing.T) {
	specialPassword := "p@ssw0rd!@#$%^&*()"

	user := &User{}
	err := user.HashPassword(specialPassword)
	if err != nil {
		t.Fatalf("HashPassword failed for special password: %v", err)
	}

	if !user.CheckPassword(specialPassword) {
		t.Error("CheckPassword failed for special password")
	}
}

func TestUser_PasswordUnicode(t *testing.T) {
	unicodePassword := "пароль123"

	user := &User{}
	err := user.HashPassword(unicodePassword)
	if err != nil {
		t.Fatalf("HashPassword failed for unicode password: %v", err)
	}

	if !user.CheckPassword(unicodePassword) {
		t.Error("CheckPassword failed for unicode password")
	}
}