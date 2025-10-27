package delivery

import (
	"backend/apiutils"
	"backend/middleware"
	"backend/named_errors"
	"backend/store"
	"backend/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
	minioStore  *store.MinIOStore
}

func NewUserHandler(userUsecase *usecase.UserUsecase, minioStore *store.MinIOStore) *UserHandler {
	return &UserHandler{userUsecase: userUsecase, minioStore: minioStore}
}

type RegisterRequest struct {
	Email    string  `json:"email"`
	Username *string `json:"username,omitempty"`
	Password string  `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
}

type PasswordResetRequest struct {
	Email string `json:"email"`
}

type PasswordResetConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// @Summary Register a new user
// @Description Creates a new user account. Can accept 'application/json' or 'multipart/form-data' if an avatar is being uploaded.
// @Tags users
// @Accept  json
// @Accept  mpfd
// @Produce json
// @Param   registerRequest body RegisterRequest true "User registration data (if using JSON)"
// @Param   email    formData string true  "User's email (if using multipart)"
// @Param   password formData string true  "User's password (if using multipart)"
// @Param   username formData string false "User's username (if using multipart)"
// @Param   avatar   formData file   false "User's avatar image (if using multipart)"
// @Success 201 {object} models.User
// @Failure 400 {object} apiutils.ErrorResponse "Invalid request body or required fields are missing"
// @Failure 409 {object} apiutils.ErrorResponse "User with this email or username already exists"
// @Failure 500 {object} apiutils.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// ... (логика парсинга multipart/form или json остается)
	contentType := r.Header.Get("Content-Type")

	var email, password, username string
	var usernamePtr *string

	if strings.HasPrefix(contentType, "multipart/form-data") {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse multipart form")
			apiutils.WriteError(w, http.StatusBadRequest, "failed to parse form")
			return
		}

		email = r.FormValue("email")
		password = r.FormValue("password")
		username = r.FormValue("username")
	} else {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("failed to decode register request")
			apiutils.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		email = req.Email
		password = req.Password
		if req.Username != nil {
			username = *req.Username
		}
	}

	if email == "" || password == "" {
		apiutils.WriteError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	if username != "" {
		usernamePtr = &username
	}

	user, err := h.userUsecase.RegisterUser(email, usernamePtr, password)
	if err != nil {
		if errors.Is(err, named_errors.ErrConflict) {
			log.Warn().Err(err).Msg("registration conflict")
			apiutils.WriteError(w, http.StatusConflict, "user with this email or username already exists")
			return
		}

		log.Error().Err(err).Msg("failed to register user")
		apiutils.WriteError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	if strings.HasPrefix(contentType, "multipart/form-data") {
		file, header, err := r.FormFile("avatar")
		if err == nil {
			defer file.Close()

			contentType := header.Header.Get("Content-Type")
			if h.minioStore.ValidateImageType(contentType) {
				if header.Size <= 5<<20 {
					oldObjectName := fmt.Sprintf("avatar_user_%d", user.ID)
					h.minioStore.DeleteAvatar(oldObjectName)

					objectName, uploadErr := h.minioStore.UploadAvatar(user.ID, file, header.Filename, contentType)
					if uploadErr == nil {
						avatarURL := h.minioStore.GetAvatarURL(objectName)
						if updateErr := h.userUsecase.UpdateUserAvatar(user.ID, avatarURL); updateErr == nil {
							user.AvatarURL = &avatarURL
						}
					}
				}
			}
		}
	}

	apiutils.WriteJSON(w, http.StatusCreated, user)
}

// @Summary Log in a user
// @Description Authenticates a user with email and password, and sets a session_id cookie upon success.
// @Tags users
// @Accept  json
// @Produce json
// @Param   credentials body LoginRequest true "User login credentials"
// @Success 200 {object} map[string]interface{} "Returns success message and user_id"
// @Failure 400 {object} apiutils.ErrorResponse "Invalid request body"
// @Failure 401 {object} apiutils.ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("failed to decode login request")
		apiutils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		apiutils.WriteError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	session, err := h.userUsecase.LoginUser(req.Email, req.Password)
	if err != nil {
		log.Error().Err(err).Msg("failed to login user")
		apiutils.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session["token"].(string),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	apiutils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "login successful",
		"user_id": session["user_id"],
	})
}

// @Summary Log out a user
// @Description Invalidates the current user's session by clearing the session_id cookie.
// @Tags users
// @Produce json
// @Success 200 {object} map[string]string "Returns success message"
// @Failure 400 {object} apiutils.ErrorResponse "No active session to log out"
// @Failure 500 {object} apiutils.ErrorResponse "Internal server error during logout"
// @Security ApiKeyAuth
// @Router /auth/logout [post]
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err != nil {
		apiutils.WriteError(w, http.StatusBadRequest, "no active session")
		return
	}

	if err := h.userUsecase.LogoutUser(session.Value); err != nil {
		log.Error().Err(err).Msg("failed to logout user")
		apiutils.WriteError(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	apiutils.WriteJSON(w, http.StatusOK, map[string]string{"message": "logout successful"})
}

// @Summary Get current user's profile
// @Description Retrieves the profile information for the currently authenticated user.
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 404 {object} apiutils.ErrorResponse "User not found"
// @Security ApiKeyAuth
// @Router /profile/me [get]
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	user, err := h.userUsecase.GetUserProfile(userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user profile")
		apiutils.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, user)
}

// @Summary Update current user's profile
// @Description Updates the profile information (username, password, avatar) for the currently authenticated user.
// @Tags users
// @Accept  json
// @Accept  mpfd
// @Produce json
// @Param   updateRequest body UpdateProfileRequest true "Profile data to update (if using JSON)"
// @Param   username formData string false "New username (if using multipart)"
// @Param   password formData string false "New password (if using multipart)"
// @Param   avatar   formData file   false "New avatar image (if using multipart)"
// @Success 200 {object} models.User
// @Failure 400 {object} apiutils.ErrorResponse "Invalid request body or form data"
// @Failure 401 {object} apiutils.ErrorResponse "User not authenticated"
// @Failure 500 {object} apiutils.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /profile/me [put]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		apiutils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	contentType := r.Header.Get("Content-Type")

	var username, password string

	if strings.HasPrefix(contentType, "multipart/form-data") {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse multipart form")
			apiutils.WriteError(w, http.StatusBadRequest, "failed to parse form")
			return
		}

		username = r.FormValue("username")
		password = r.FormValue("password")
	} else {
		var req UpdateProfileRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error().Err(err).Msg("failed to decode update profile request")
			apiutils.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Username != nil {
			username = *req.Username
		}
		if req.Password != nil {
			password = *req.Password
		}
	}

	user, err := h.userUsecase.GetUserProfile(userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user profile")
		apiutils.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	if username != "" {
		user.Username = username
	}

	if password != "" {
		if err := user.HashPassword(password); err != nil {
			log.Error().Err(err).Msg("failed to hash password")
			apiutils.WriteError(w, http.StatusInternalServerError, "failed to update password")
			return
		}
	}

	if strings.HasPrefix(contentType, "multipart/form-data") {
		file, header, err := r.FormFile("avatar")
		if err == nil {
			defer file.Close()

			contentType := header.Header.Get("Content-Type")
			if h.minioStore.ValidateImageType(contentType) {
				if header.Size <= 5<<20 {
					oldObjectName := fmt.Sprintf("avatar_user_%d", user.ID)
					h.minioStore.DeleteAvatar(oldObjectName)

					objectName, uploadErr := h.minioStore.UploadAvatar(user.ID, file, header.Filename, contentType)
					if uploadErr == nil {
						avatarURL := h.minioStore.GetAvatarURL(objectName)
						user.AvatarURL = &avatarURL
					}
				}
			}
		}
	}

	if err := h.userUsecase.UpdateUser(user); err != nil {
		log.Error().Err(err).Msg("failed to update user profile")
		apiutils.WriteError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/files/")
	if path == "" {
		apiutils.WriteError(w, http.StatusBadRequest, "file path is required")
		return
	}

	object, err := h.minioStore.GetObject(path)
	if err != nil {
		log.Error().Err(err).Msg("failed to get object from MinIO")
		apiutils.WriteError(w, http.StatusNotFound, "file not found")
		return
	}
	defer object.(io.Closer).Close()

	info, err := h.minioStore.GetObjectInfo(path)
	if err != nil {
		log.Error().Err(err).Msg("failed to get object info")
		apiutils.WriteError(w, http.StatusInternalServerError, "failed to get file info")
		return
	}

	w.Header().Set("Content-Type", info.ContentType)
	w.Header().Set("Content-Length", string(rune(info.Size)))
	w.Header().Set("Cache-Control", "public, max-age=3600")

	_, err = io.Copy(w, object)
	if err != nil {
		log.Error().Err(err).Msg("failed to copy file content")
		return
	}
}

// @Summary Request password reset
// @Description Sends a password reset email to the user
// @Tags users
// @Accept  json
// @Produce json
// @Param   request body PasswordResetRequest true "Password reset request"
// @Success 200 {object} map[string]string "Returns success message"
// @Failure 400 {object} apiutils.ErrorResponse "Invalid request body"
// @Failure 500 {object} apiutils.ErrorResponse "Internal server error"
// @Router /auth/forgot-password [post]
func (h *UserHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("failed to decode password reset request")
		apiutils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" {
		apiutils.WriteError(w, http.StatusBadRequest, "email is required")
		return
	}

	if err := h.userUsecase.RequestPasswordReset(req.Email); err != nil {
		log.Error().Err(err).Msg("failed to request password reset")
		apiutils.WriteError(w, http.StatusInternalServerError, "failed to process password reset request")
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "If an account with that email exists, a password reset link has been sent",
	})
}

// @Summary Reset password
// @Description Resets the user's password using a valid reset token
// @Tags users
// @Accept  json
// @Produce json
// @Param   request body PasswordResetConfirmRequest true "Password reset confirmation"
// @Success 200 {object} map[string]string "Returns success message"
// @Failure 400 {object} apiutils.ErrorResponse "Invalid request body"
// @Failure 404 {object} apiutils.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} apiutils.ErrorResponse "Internal server error"
// @Router /auth/reset-password [post]
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req PasswordResetConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("failed to decode password reset confirm request")
		apiutils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		apiutils.WriteError(w, http.StatusBadRequest, "token and new password are required")
		return
	}

	if err := h.userUsecase.ResetPassword(req.Token, req.NewPassword); err != nil {
		if errors.Is(err, named_errors.ErrNotFound) {
			apiutils.WriteError(w, http.StatusNotFound, "invalid or expired token")
			return
		}
		log.Error().Err(err).Msg("failed to reset password")
		apiutils.WriteError(w, http.StatusInternalServerError, "failed to reset password")
		return
	}

	apiutils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Password has been reset successfully",
	})
}
