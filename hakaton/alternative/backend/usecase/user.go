package usecase

import (
	"backend/models"
	"backend/named_errors"
	"backend/repository"
	"backend/services"
	"backend/store"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserUsecase struct {
	repo         *repository.Repository
	redisStore   *store.RedisStore
	minioStore   *store.MinIOStore
	emailService *services.EmailService
}

func NewUserUsecase(repo *repository.Repository, redisStore *store.RedisStore, minioStore *store.MinIOStore, emailService *services.EmailService) *UserUsecase {
	return &UserUsecase{
		repo:         repo,
		redisStore:   redisStore,
		minioStore:   minioStore,
		emailService: emailService,
	}
}

func (u *UserUsecase) RegisterUser(email string, username *string, password string) (*models.User, error) {
	_, err := u.repo.GetUserByEmail(email)
	if err == nil {
		return nil, named_errors.ErrConflict
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("usecase.RegisterUser: failed to check email: %w", err)
	}

	var finalUsername string
	if username != nil && *username != "" {
		finalUsername = *username
	} else {
		finalUsername = email
	}
	
	user := &models.User{
		Email:    email,
		Username: finalUsername,
	}

	if err := user.HashPassword(password); err != nil {
		return nil, fmt.Errorf("usecase.RegisterUser: failed to hash password: %w", err)
	}

	if err := u.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("usecase.RegisterUser: failed to create user: %w", err)
	}

	return user, nil
}

func (u *UserUsecase) LoginUser(email, password string) (map[string]interface{}, error) {
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("invalid credentials")
	}

	session, err := u.redisStore.CreateSession(user.ID)
	if err != nil {
		return nil, fmt.Errorf("usecase.LoginUser: failed to create session: %w", err)
	}

	return session, nil
}

func (u *UserUsecase) GetUserProfile(userID uint64) (*models.User, error) {
	user, err := u.repo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, named_errors.ErrNotFound
		}
		return nil, fmt.Errorf("usecase.GetUserProfile: %w", err)
	}
	return user, nil
}

func (u *UserUsecase) UpdateUser(user *models.User) error {
	if err := u.repo.UpdateUser(user); err != nil {
		return fmt.Errorf("usecase.UpdateUser: %w", err)
	}
	return nil
}

func (u *UserUsecase) UpdateUserAvatar(userID uint64, avatarURL string) error {
	user, err := u.repo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("usecase.UpdateUserAvatar: user not found: %w", err)
	}
	user.AvatarURL = &avatarURL
	return u.repo.UpdateUser(user)
}

func (u *UserUsecase) LogoutUser(token string) error {
	return u.redisStore.DeleteSession(token)
}

func (u *UserUsecase) GetUserBySession(token string) (*models.User, bool) {
	return u.redisStore.GetUserBySession(token)
}

func (u *UserUsecase) RequestPasswordReset(email string) error {
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal if user exists or not for security
			return nil
		}
		return fmt.Errorf("usecase.RequestPasswordReset: failed to get user: %w", err)
	}

	// Generate a secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("usecase.RequestPasswordReset: failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Create password reset token
	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour), // Token expires in 1 hour
	}

	// Save token to database
	if err := u.repo.CreatePasswordResetToken(resetToken); err != nil {
		return fmt.Errorf("usecase.RequestPasswordReset: failed to create reset token: %w", err)
	}

	// Send email
	if err := u.emailService.SendPasswordResetEmail(user.Email, token); err != nil {
		return fmt.Errorf("usecase.RequestPasswordReset: failed to send email: %w", err)
	}

	return nil
}

func (u *UserUsecase) ResetPassword(token, newPassword string) error {
	// Get the reset token
	resetToken, err := u.repo.GetPasswordResetToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return named_errors.ErrNotFound
		}
		return fmt.Errorf("usecase.ResetPassword: failed to get reset token: %w", err)
	}

	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		return named_errors.ErrNotFound
	}

	// Get user
	user, err := u.repo.GetUserByID(resetToken.UserID)
	if err != nil {
		return fmt.Errorf("usecase.ResetPassword: failed to get user: %w", err)
	}

	// Update password
	if err := user.HashPassword(newPassword); err != nil {
		return fmt.Errorf("usecase.ResetPassword: failed to hash password: %w", err)
	}

	if err := u.repo.UpdateUser(user); err != nil {
		return fmt.Errorf("usecase.ResetPassword: failed to update user: %w", err)
	}

	// Mark token as used
	if err := u.repo.MarkPasswordResetTokenAsUsed(token); err != nil {
		return fmt.Errorf("usecase.ResetPassword: failed to mark token as used: %w", err)
	}

	return nil
}