package repository

import (
	"backend/models"
	"fmt"
)


func (r *Repository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("repository.CreateUser: %w", err)
	}
	return nil
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("repository.GetUserByEmail: %w", err)
	}
	return &user, nil
}

func (r *Repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("repository.GetUserByUsername: %w", err)
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id uint64) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("repository.GetUserByID: %w", err)
	}
	return &user, nil
}

func (r *Repository) UpdateUser(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("repository.UpdateUser: %w", err)
	}
	return nil
}

func (r *Repository) CreatePasswordResetToken(token *models.PasswordResetToken) error {
	if err := r.db.Create(token).Error; err != nil {
		return fmt.Errorf("repository.CreatePasswordResetToken: %w", err)
	}
	return nil
}

func (r *Repository) GetPasswordResetToken(token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	if err := r.db.Where("token = ? AND used = ?", token, false).First(&resetToken).Error; err != nil {
		return nil, fmt.Errorf("repository.GetPasswordResetToken: %w", err)
	}
	return &resetToken, nil
}

func (r *Repository) MarkPasswordResetTokenAsUsed(token string) error {
	if err := r.db.Model(&models.PasswordResetToken{}).Where("token = ?", token).Update("used", true).Error; err != nil {
		return fmt.Errorf("repository.MarkPasswordResetTokenAsUsed: %w", err)
	}
	return nil
}

func (r *Repository) DeleteExpiredPasswordResetTokens() error {
	if err := r.db.Where("expires_at < ?", "NOW()").Delete(&models.PasswordResetToken{}).Error; err != nil {
		return fmt.Errorf("repository.DeleteExpiredPasswordResetTokens: %w", err)
	}
	return nil
}