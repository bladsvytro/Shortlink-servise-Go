package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/database"
	"url-shortener/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Application represents the main application
type Application struct {
	config      *config.Config
	logger      *zap.Logger
	server      *http.Server
	db          *database.Database
	rateLimiter *RateLimiter
}

// New creates a new Application instance
func New(cfg *config.Config, logger *zap.Logger) (*Application, error) {
	// Initialize database
	db, err := database.New(cfg.Database, logger)
	if err != nil {
		return nil, err
	}

	app := &Application{
		config: cfg,
		logger: logger,
		db:     db,
	}

	// Initialize HTTP server
	app.server = &http.Server{
		Addr:         cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
		Handler:      app.setupRouter(),
	}

	return app, nil
}

// Start starts the application
func (a *Application) Start() error {
	a.logger.Info("Starting application",
		zap.Int("port", a.config.Server.Port),
		zap.String("environment", a.config.Server.Env),
	)

	// Run database migrations
	if err := a.db.Migrate(); err != nil {
		return err
	}

	// Start HTTP server
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	return nil
}

// Stop gracefully stops the application
func (a *Application) Stop(ctx context.Context) error {
	a.logger.Info("Stopping application...")

	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	// Close database connection
	if err := a.db.Close(); err != nil {
		return err
	}

	// Stop rate limiter cleanup goroutine
	if a.rateLimiter != nil {
		a.rateLimiter.Stop()
	}

	return nil
}

// Router returns the HTTP handler for testing purposes
func (a *Application) Router() http.Handler {
	return a.server.Handler
}

// DB returns the database connection for testing purposes
func (a *Application) DB() *database.Database {
	return a.db
}

// setupRouter sets up the HTTP router
func (a *Application) setupRouter() http.Handler {
	mux := http.NewServeMux()

	// Serve static files from web directory
	fs := http.FileServer(http.Dir("web"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	// Serve index.html for root path
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// If the request is for the root path, serve index.html
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "web/index.html")
			return
		}
		// For other paths, try to serve static files, otherwise fall through to API
		// This allows direct access to static files if they exist
		fs.ServeHTTP(w, r)
	})

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		a.logger.Debug("Health check requested", zap.String("path", r.URL.Path))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Authentication endpoints
	mux.HandleFunc("POST /api/v1/auth/register", a.handleRegister)
	mux.HandleFunc("POST /api/v1/auth/login", a.handleLogin)
	mux.HandleFunc("GET /api/v1/auth/me", a.AuthMiddleware(a.handleMe))

	// URL shortening endpoints
	mux.HandleFunc("POST /api/v1/links", a.AuthMiddleware(a.handleCreateLink))
	mux.HandleFunc("GET /api/v1/links", a.AuthMiddleware(a.handleListLinks))
	mux.HandleFunc("GET /{code}", a.handleRedirect)
	mux.HandleFunc("GET /api/v1/links/{code}/stats", a.handleLinkStats)

	// User statistics endpoint
	mux.HandleFunc("GET /api/v1/stats", a.AuthMiddleware(a.handleUserStats))

	// API key management endpoints
	mux.HandleFunc("GET /api/v1/api-keys", a.AuthMiddleware(a.handleListAPIKeys))
	mux.HandleFunc("POST /api/v1/api-keys", a.AuthMiddleware(a.handleCreateAPIKey))
	mux.HandleFunc("DELETE /api/v1/api-keys/{id}", a.AuthMiddleware(a.handleDeleteAPIKey))

	// Domain management endpoints
	mux.HandleFunc("GET /api/v1/domains", a.AuthMiddleware(a.handleListDomains))
	mux.HandleFunc("POST /api/v1/domains", a.AuthMiddleware(a.handleCreateDomain))
	mux.HandleFunc("GET /api/v1/domains/{id}", a.AuthMiddleware(a.handleGetDomain))
	mux.HandleFunc("DELETE /api/v1/domains/{id}", a.AuthMiddleware(a.handleDeleteDomain))
	mux.HandleFunc("POST /api/v1/domains/{id}/verify", a.AuthMiddleware(a.handleVerifyDomain))

	// Admin endpoints
	mux.HandleFunc("GET /admin/users", a.AdminMiddleware(a.handleAdminListUsers))
	mux.HandleFunc("GET /admin/links", a.AdminMiddleware(a.handleAdminListLinks))
	mux.HandleFunc("GET /admin/stats", a.AdminMiddleware(a.handleAdminStats))

	return mux
}

// handleCreateLink handles POST /api/v1/links
func (a *Application) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL         string     `json:"url"`
		CustomCode  string     `json:"custom_code,omitempty"`
		Title       string     `json:"title,omitempty"`
		Description string     `json:"description,omitempty"`
		ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error("Failed to decode request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate URL
	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		req.URL = "https://" + req.URL
	}

	// Determine user ID from authentication (if any)
	userID := uuid.Nil
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = parsed
		}
	}

	// Generate short code
	shortCode := req.CustomCode
	customCodeProvided := shortCode != ""

	if customCodeProvided {
		// Validate custom code
		if !a.config.Shortener.AllowCustomCodes {
			http.Error(w, "Custom codes are not allowed", http.StatusBadRequest)
			return
		}

		// Check length
		if len(shortCode) > a.config.Shortener.MaxCustomCodeLength {
			http.Error(w, fmt.Sprintf("Custom code too long (max %d characters)", a.config.Shortener.MaxCustomCodeLength), http.StatusBadRequest)
			return
		}

		// Validate characters: alphanumeric, underscore, hyphen
		for _, ch := range shortCode {
			if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '-') {
				http.Error(w, "Custom code can only contain letters, numbers, underscores and hyphens", http.StatusBadRequest)
				return
			}
		}

		// Check if user is authenticated (custom codes only for registered users)
		if userID == uuid.Nil {
			http.Error(w, "Custom codes are only available for registered users", http.StatusUnauthorized)
			return
		}
	} else {
		// Generate random short code
		shortCode = GenerateShortCode(a.config.Shortener.ShortCodeLength)
	}

	// Check if short code already exists
	var existingLink models.Link
	if err := a.db.DB.Where("short_code = ?", shortCode).First(&existingLink).Error; err == nil {
		http.Error(w, "Short code already exists", http.StatusConflict)
		return
	}

	// Create link
	link := models.Link{
		ShortCode:   shortCode,
		OriginalURL: req.URL,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		ExpiresAt:   req.ExpiresAt,
		IsActive:    true,
		ClickCount:  0,
	}

	if err := a.db.DB.Create(&link).Error; err != nil {
		a.logger.Error("Failed to create link", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build response
	resp := struct {
		ShortCode   string    `json:"short_code"`
		ShortURL    string    `json:"short_url"`
		OriginalURL string    `json:"original_url"`
		CreatedAt   time.Time `json:"created_at"`
	}{
		ShortCode:   link.ShortCode,
		ShortURL:    link.GetShortURL(a.config.Shortener.BaseURL),
		OriginalURL: link.OriginalURL,
		CreatedAt:   link.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// handleRedirect handles GET /{code}
func (a *Application) handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	var link models.Link
	if err := a.db.DB.Where("short_code = ?", code).First(&link).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Link not found", http.StatusNotFound)
			return
		}
		a.logger.Error("Database error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if link is accessible
	if !link.CanBeAccessed() {
		http.Error(w, "Link is inactive or expired", http.StatusGone)
		return
	}

	// Increment click count (async)
	go func() {
		link.IncrementClickCount()
		a.db.DB.Model(&link).Updates(map[string]interface{}{
			"click_count":     link.ClickCount,
			"last_clicked_at": link.LastClickedAt,
		})
	}()

	// Redirect
	http.Redirect(w, r, link.OriginalURL, a.config.Shortener.DefaultRedirectCode)
}

// handleLinkStats handles GET /api/v1/links/{code}/stats
func (a *Application) handleLinkStats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	var link models.Link
	if err := a.db.DB.Where("short_code = ?", code).First(&link).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Link not found", http.StatusNotFound)
			return
		}
		a.logger.Error("Database error", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := struct {
		ShortCode     string     `json:"short_code"`
		OriginalURL   string     `json:"original_url"`
		ClickCount    int64      `json:"click_count"`
		LastClickedAt *time.Time `json:"last_clicked_at,omitempty"`
		CreatedAt     time.Time  `json:"created_at"`
		ExpiresAt     *time.Time `json:"expires_at,omitempty"`
		IsActive      bool       `json:"is_active"`
	}{
		ShortCode:     link.ShortCode,
		OriginalURL:   link.OriginalURL,
		ClickCount:    link.ClickCount,
		LastClickedAt: link.LastClickedAt,
		CreatedAt:     link.CreatedAt,
		ExpiresAt:     link.ExpiresAt,
		IsActive:      link.IsActive,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleUserStats handles GET /api/v1/stats
func (a *Application) handleUserStats(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authentication header (set by AuthMiddleware)
	userID := uuid.Nil
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = parsed
		}
	}
	
	if userID == uuid.Nil {
		a.logger.Error("User ID not found in headers")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Calculate total links for this user
	var totalLinks int64
	if err := a.db.DB.Model(&models.Link{}).Where("user_id = ?", userID).Count(&totalLinks).Error; err != nil {
		a.logger.Error("Database error counting links", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Calculate total clicks for this user
	var totalClicks int64
	if err := a.db.DB.Model(&models.Link{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(click_count), 0)").
		Scan(&totalClicks).Error; err != nil {
		a.logger.Error("Database error summing clicks", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Calculate today's clicks (clicks where last_clicked_at >= today)
	var todayClicks int64
	today := time.Now().Truncate(24 * time.Hour)
	if err := a.db.DB.Model(&models.Link{}).
		Where("user_id = ? AND last_clicked_at >= ?", userID, today).
		Select("COALESCE(SUM(click_count), 0)").
		Scan(&todayClicks).Error; err != nil {
		a.logger.Error("Database error summing today's clicks", zap.Error(err))
		// If error, just set to 0
		todayClicks = 0
	}

	// Prepare response
	resp := struct {
		TotalLinks  int64 `json:"total_links"`
		TotalClicks int64 `json:"total_clicks"`
		TodayClicks int64 `json:"today_clicks"`
	}{
		TotalLinks:  totalLinks,
		TotalClicks: totalClicks,
		TodayClicks: todayClicks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleListLinks handles GET /api/v1/links
func (a *Application) handleListLinks(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authentication header
	userID := uuid.Nil
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = parsed
		}
	}
	
	if userID == uuid.Nil {
		a.logger.Error("User ID not found in headers")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse query parameters
	limit := 50
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Fetch links for this user
	var links []models.Link
	if err := a.db.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&links).Error; err != nil {
		a.logger.Error("Database error fetching links", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build response
	type linkResponse struct {
		ID           string     `json:"id"`
		ShortCode    string     `json:"short_code"`
		OriginalURL  string     `json:"original_url"`
		Title        string     `json:"title,omitempty"`
		Description  string     `json:"description,omitempty"`
		ClickCount   int64      `json:"click_count"`
		LastClickedAt *time.Time `json:"last_clicked_at,omitempty"`
		CreatedAt    time.Time  `json:"created_at"`
		ExpiresAt    *time.Time `json:"expires_at,omitempty"`
		IsActive     bool       `json:"is_active"`
		ShortURL     string     `json:"short_url"`
	}

	resp := make([]linkResponse, len(links))
	for i, link := range links {
		resp[i] = linkResponse{
			ID:            link.ID.String(),
			ShortCode:     link.ShortCode,
			OriginalURL:   link.OriginalURL,
			Title:         link.Title,
			Description:   link.Description,
			ClickCount:    link.ClickCount,
			LastClickedAt: link.LastClickedAt,
			CreatedAt:     link.CreatedAt,
			ExpiresAt:     link.ExpiresAt,
			IsActive:      link.IsActive,
			ShortURL:      link.GetShortURL(a.config.Shortener.BaseURL),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleListAPIKeys handles GET /api/v1/api-keys
func (a *Application) handleListAPIKeys(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authentication header
	userID := uuid.Nil
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = parsed
		}
	}
	
	if userID == uuid.Nil {
		a.logger.Error("User ID not found in headers")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch API keys for this user
	var apiKeys []models.APIKey
	if err := a.db.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&apiKeys).Error; err != nil {
		a.logger.Error("Database error fetching API keys", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build response (exclude key hash for security)
	type apiKeyResponse struct {
		ID         string     `json:"id"`
		Name       string     `json:"name"`
		LastUsedAt *time.Time `json:"last_used_at,omitempty"`
		ExpiresAt  *time.Time `json:"expires_at,omitempty"`
		RateLimit  int        `json:"rate_limit"`
		CreatedAt  time.Time  `json:"created_at"`
		IsExpired  bool       `json:"is_expired"`
	}

	resp := make([]apiKeyResponse, len(apiKeys))
	for i, key := range apiKeys {
		resp[i] = apiKeyResponse{
			ID:         key.ID.String(),
			Name:       key.Name,
			LastUsedAt: key.LastUsedAt,
			ExpiresAt:  key.ExpiresAt,
			RateLimit:  key.RateLimit,
			CreatedAt:  key.CreatedAt,
			IsExpired:  key.IsExpired(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleCreateAPIKey handles POST /api/v1/api-keys
func (a *Application) handleCreateAPIKey(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authentication header
	userID := uuid.Nil
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = parsed
		}
	}
	
	if userID == uuid.Nil {
		a.logger.Error("User ID not found in headers")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name      string     `json:"name"`
		ExpiresAt *time.Time `json:"expires_at,omitempty"`
		RateLimit int        `json:"rate_limit,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error("Failed to decode request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Generate API key (raw key)
	rawKey := uuid.New().String() + uuid.New().String() // 64 chars
	keyHash := sha256Hash(rawKey)

	// Set default rate limit
	if req.RateLimit <= 0 {
		req.RateLimit = 1000
	}

	// Create API key record
	apiKey := models.APIKey{
		UserID:    userID,
		KeyHash:   keyHash,
		Name:      req.Name,
		ExpiresAt: req.ExpiresAt,
		RateLimit: req.RateLimit,
	}

	if err := a.db.DB.Create(&apiKey).Error; err != nil {
		a.logger.Error("Failed to create API key", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Response includes the raw key (only shown once)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         apiKey.ID.String(),
		"name":       apiKey.Name,
		"key":        rawKey, // Only time the raw key is returned
		"expires_at": apiKey.ExpiresAt,
		"rate_limit": apiKey.RateLimit,
		"created_at": apiKey.CreatedAt,
	})
}

// handleDeleteAPIKey handles DELETE /api/v1/api-keys/{id}
func (a *Application) handleDeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authentication header
	userID := uuid.Nil
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = parsed
		}
	}
	
	if userID == uuid.Nil {
		a.logger.Error("User ID not found in headers")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	keyID := r.PathValue("id")
	if keyID == "" {
		http.Error(w, "API key ID is required", http.StatusBadRequest)
		return
	}

	// Parse UUID
	parsedKeyID, err := uuid.Parse(keyID)
	if err != nil {
		http.Error(w, "Invalid API key ID", http.StatusBadRequest)
		return
	}

	// Delete only if belongs to user
	result := a.db.DB.Where("id = ? AND user_id = ?", parsedKeyID, userID).Delete(&models.APIKey{})
	if result.Error != nil {
		a.logger.Error("Database error deleting API key", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "API key not found or access denied", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleListDomains handles GET /api/v1/domains
func (a *Application) handleListDomains(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authentication header
	userID := uuid.Nil
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = parsed
		}
	}
	
	if userID == uuid.Nil {
		a.logger.Error("User ID not found in headers")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch domains for this user
	var domains []models.Domain
	if err := a.db.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&domains).Error; err != nil {
		a.logger.Error("Database error fetching domains", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build response
	type domainResponse struct {
		ID               string     `json:"id"`
		DomainName       string     `json:"domain_name"`
		IsVerified       bool       `json:"is_verified"`
		IsActive         bool       `json:"is_active"`
		VerifiedAt       *time.Time `json:"verified_at,omitempty"`
		CreatedAt        time.Time  `json:"created_at"`
		CanBeUsed        bool       `json:"can_be_used"`
		VerificationToken string    `json:"verification_token,omitempty"`
	}

	resp := make([]domainResponse, len(domains))
	for i, domain := range domains {
		resp[i] = domainResponse{
			ID:                domain.ID.String(),
			DomainName:        domain.DomainName,
			IsVerified:        domain.IsVerified,
			IsActive:          domain.IsActive,
			VerifiedAt:        domain.VerifiedAt,
			CreatedAt:         domain.CreatedAt,
			CanBeUsed:         domain.CanBeUsed(),
			VerificationToken: domain.VerificationToken,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleCreateDomain handles POST /api/v1/domains
func (a *Application) handleCreateDomain(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authentication header
	userID := uuid.Nil
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = parsed
		}
	}
	
	if userID == uuid.Nil {
		a.logger.Error("User ID not found in headers")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		DomainName string `json:"domain_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.logger.Error("Failed to decode request", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.DomainName == "" {
		http.Error(w, "Domain name is required", http.StatusBadRequest)
		return
	}

	// Basic domain validation
	if !strings.Contains(req.DomainName, ".") {
		http.Error(w, "Invalid domain name", http.StatusBadRequest)
		return
	}

	// Check if domain already exists (globally)
	var existingDomain models.Domain
	if err := a.db.DB.Where("domain_name = ?", req.DomainName).First(&existingDomain).Error; err == nil {
		http.Error(w, "Domain already registered", http.StatusConflict)
		return
	}

	// Generate verification token
	token := uuid.New().String()

	// Create domain record
	domain := models.Domain{
		DomainName:        req.DomainName,
		UserID:            userID,
		IsVerified:        false,
		IsActive:          true,
		VerificationToken: token,
	}

	if err := a.db.DB.Create(&domain).Error; err != nil {
		a.logger.Error("Failed to create domain", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":                 domain.ID.String(),
		"domain_name":        domain.DomainName,
		"is_verified":        domain.IsVerified,
		"is_active":          domain.IsActive,
		"verification_token": domain.VerificationToken,
		"created_at":         domain.CreatedAt,
		"verification_instructions": "Add a TXT record to your DNS with the value: " + token,
	})
}

// handleGetDomain handles GET /api/v1/domains/{id}
func (a *Application) handleGetDomain(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authentication header
	userID := uuid.Nil
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = parsed
		}
	}
	
	if userID == uuid.Nil {
		a.logger.Error("User ID not found in headers")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	domainID := r.PathValue("id")
	if domainID == "" {
		http.Error(w, "Domain ID is required", http.StatusBadRequest)
		return
	}

	// Parse UUID
	parsedDomainID, err := uuid.Parse(domainID)
	if err != nil {
		http.Error(w, "Invalid domain ID", http.StatusBadRequest)
		return
	}

	// Fetch domain
	var domain models.Domain
	if err := a.db.DB.Where("id = ? AND user_id = ?", parsedDomainID, userID).First(&domain).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Domain not found or access denied", http.StatusNotFound)
			return
		}
		a.logger.Error("Database error fetching domain", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":                 domain.ID.String(),
		"domain_name":        domain.DomainName,
		"is_verified":        domain.IsVerified,
		"is_active":          domain.IsActive,
		"verified_at":        domain.VerifiedAt,
		"created_at":         domain.CreatedAt,
		"can_be_used":        domain.CanBeUsed(),
		"verification_token": domain.VerificationToken,
	})
}

// handleDeleteDomain handles DELETE /api/v1/domains/{id}
func (a *Application) handleDeleteDomain(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authentication header
	userID := uuid.Nil
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = parsed
		}
	}
	
	if userID == uuid.Nil {
		a.logger.Error("User ID not found in headers")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	domainID := r.PathValue("id")
	if domainID == "" {
		http.Error(w, "Domain ID is required", http.StatusBadRequest)
		return
	}

	// Parse UUID
	parsedDomainID, err := uuid.Parse(domainID)
	if err != nil {
		http.Error(w, "Invalid domain ID", http.StatusBadRequest)
		return
	}

	// Delete only if belongs to user
	result := a.db.DB.Where("id = ? AND user_id = ?", parsedDomainID, userID).Delete(&models.Domain{})
	if result.Error != nil {
		a.logger.Error("Database error deleting domain", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Domain not found or access denied", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleVerifyDomain handles POST /api/v1/domains/{id}/verify
func (a *Application) handleVerifyDomain(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authentication header
	userID := uuid.Nil
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = parsed
		}
	}
	
	if userID == uuid.Nil {
		a.logger.Error("User ID not found in headers")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	domainID := r.PathValue("id")
	if domainID == "" {
		http.Error(w, "Domain ID is required", http.StatusBadRequest)
		return
	}

	// Parse UUID
	parsedDomainID, err := uuid.Parse(domainID)
	if err != nil {
		http.Error(w, "Invalid domain ID", http.StatusBadRequest)
		return
	}

	// Fetch domain
	var domain models.Domain
	if err := a.db.DB.Where("id = ? AND user_id = ?", parsedDomainID, userID).First(&domain).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Domain not found or access denied", http.StatusNotFound)
			return
		}
		a.logger.Error("Database error fetching domain", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// For MVP, we'll just mark as verified (in production, would check DNS)
	domain.Verify()
	if err := a.db.DB.Save(&domain).Error; err != nil {
		a.logger.Error("Failed to verify domain", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":           domain.ID.String(),
		"domain_name":  domain.DomainName,
		"is_verified":  domain.IsVerified,
		"verified_at":  domain.VerifiedAt,
		"message":      "Domain verified successfully",
	})
}

// handleAdminListUsers handles GET /admin/users
func (a *Application) handleAdminListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limit := 50
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Fetch users
	var users []models.User
	if err := a.db.DB.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		a.logger.Error("Database error fetching users", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build response
	type userResponse struct {
		ID         string     `json:"id"`
		Email      string     `json:"email"`
		Name       string     `json:"name"`
		IsActive   bool       `json:"is_active"`
		IsAdmin    bool       `json:"is_admin"`
		LastLoginAt *time.Time `json:"last_login_at,omitempty"`
		CreatedAt  time.Time  `json:"created_at"`
		LinkCount  int64      `json:"link_count"`
	}

	resp := make([]userResponse, len(users))
	for i, user := range users {
		// Count links for this user
		var linkCount int64
		a.db.DB.Model(&models.Link{}).Where("user_id = ?", user.ID).Count(&linkCount)

		resp[i] = userResponse{
			ID:          user.ID.String(),
			Email:       user.Email,
			Name:        user.Name,
			IsActive:    user.IsActive,
			IsAdmin:     user.IsAdmin,
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
			LinkCount:   linkCount,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleAdminListLinks handles GET /admin/links
func (a *Application) handleAdminListLinks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limit := 50
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Fetch links with user info
	var links []models.Link
	if err := a.db.DB.
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&links).Error; err != nil {
		a.logger.Error("Database error fetching links", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build response
	type linkResponse struct {
		ID           string     `json:"id"`
		ShortCode    string     `json:"short_code"`
		OriginalURL  string     `json:"original_url"`
		UserEmail    string     `json:"user_email"`
		ClickCount   int64      `json:"click_count"`
		LastClickedAt *time.Time `json:"last_clicked_at,omitempty"`
		CreatedAt    time.Time  `json:"created_at"`
		ExpiresAt    *time.Time `json:"expires_at,omitempty"`
		IsActive     bool       `json:"is_active"`
	}

	resp := make([]linkResponse, len(links))
	for i, link := range links {
		resp[i] = linkResponse{
			ID:            link.ID.String(),
			ShortCode:     link.ShortCode,
			OriginalURL:   link.OriginalURL,
			UserEmail:     link.User.Email,
			ClickCount:    link.ClickCount,
			LastClickedAt: link.LastClickedAt,
			CreatedAt:     link.CreatedAt,
			ExpiresAt:     link.ExpiresAt,
			IsActive:      link.IsActive,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleAdminStats handles GET /admin/stats
func (a *Application) handleAdminStats(w http.ResponseWriter, r *http.Request) {
	// Get total counts
	var totalUsers, totalLinks, totalDomains, totalClicks int64

	if err := a.db.DB.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		a.logger.Error("Database error counting users", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := a.db.DB.Model(&models.Link{}).Count(&totalLinks).Error; err != nil {
		a.logger.Error("Database error counting links", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := a.db.DB.Model(&models.Domain{}).Count(&totalDomains).Error; err != nil {
		a.logger.Error("Database error counting domains", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := a.db.DB.Model(&models.Link{}).Select("COALESCE(SUM(click_count), 0)").Scan(&totalClicks).Error; err != nil {
		a.logger.Error("Database error summing clicks", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get today's clicks
	var todayClicks int64
	today := time.Now().Truncate(24 * time.Hour)
	if err := a.db.DB.Model(&models.Link{}).
		Where("last_clicked_at >= ?", today).
		Select("COALESCE(SUM(click_count), 0)").
		Scan(&todayClicks).Error; err != nil {
		todayClicks = 0
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_users":   totalUsers,
		"total_links":   totalLinks,
		"total_domains": totalDomains,
		"total_clicks":  totalClicks,
		"today_clicks":  todayClicks,
		"timestamp":     time.Now(),
	})
}

// sha256Hash returns SHA-256 hash of a string
func sha256Hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// GenerateShortCode generates a random alphanumeric short code of given length
func GenerateShortCode(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	// Generate random bytes
	randBytes := make([]byte, length)
	_, err := rand.Read(randBytes)
	if err != nil {
		// Fallback to time-based pseudo-random (should rarely happen)
		for i := range b {
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		}
	} else {
		for i := range b {
			b[i] = charset[int(randBytes[i])%len(charset)]
		}
	}
	return string(b)
}
