package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthClaims represents JWT claims
type AuthClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash compares a plain text password with a bcrypt hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(user *models.User, cfg *config.AuthConfig) (string, error) {
	expirationTime := time.Now().Add(time.Duration(cfg.AccessTokenExpiry) * time.Second)
	claims := &AuthClaims{
		UserID:  user.ID.String(),
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "url-shortener",
			Subject:   user.ID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ValidateJWT validates a JWT token and returns claims
func ValidateJWT(tokenString string, cfg *config.AuthConfig) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// AuthMiddleware is a middleware that validates JWT token and sets user in context
func (a *Application) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]
		claims, err := ValidateJWT(tokenString, &a.config.Auth)
		if err != nil {
			a.logger.Debug("JWT validation failed", zap.Error(err))
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		// Store claims in request headers for later use (in production, use context)
		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-User-Email", claims.Email)
		r.Header.Set("X-User-IsAdmin", strconv.FormatBool(claims.IsAdmin))
		next.ServeHTTP(w, r)
	}
}

// AdminMiddleware is a middleware that requires admin privileges
func (a *Application) AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return a.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		isAdminStr := r.Header.Get("X-User-IsAdmin")
		isAdmin, _ := strconv.ParseBool(isAdminStr)
		if !isAdmin {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// handleRegister handles POST /api/v1/auth/register
func (a *Application) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error("Failed to decode registration request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Validate
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	// Check if user already exists
	var existingUser models.User
	if err := a.db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}
	// Hash password
	hash, err := HashPassword(req.Password)
	if err != nil {
		a.logger.Error("Failed to hash password", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Create user
	user := models.User{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
		IsActive:     true,
		IsAdmin:      false,
	}
	if err := a.db.DB.Create(&user).Error; err != nil {
		a.logger.Error("Failed to create user", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Generate JWT token
	token, err := GenerateJWT(&user, &a.config.Auth)
	if err != nil {
		a.logger.Error("Failed to generate JWT", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
		"token": token,
	})
}

// handleLogin handles POST /api/v1/auth/login
func (a *Application) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error("Failed to decode login request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Find user
	var user models.User
	if err := a.db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		a.logger.Error("Database error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Check password
	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	a.db.DB.Save(&user)
	// Generate JWT token
	token, err := GenerateJWT(&user, &a.config.Auth)
	if err != nil {
		a.logger.Error("Failed to generate JWT", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
		"token": token,
	})
}

// handleMe handles GET /api/v1/auth/me (returns current user info)
func (a *Application) handleMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var user models.User
	if err := a.db.DB.First(&user, "id = ?", userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"name":       user.Name,
		"is_active":  user.IsActive,
		"is_admin":   user.IsAdmin,
		"created_at": user.CreatedAt,
	})
}
