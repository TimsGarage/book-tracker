package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bookend/internal/config"
	"bookend/internal/database"
	"bookend/internal/models"
	"bookend/internal/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// In-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Book{}); err != nil {
		t.Fatalf("failed to migrate test db: %v", err)
	}

	cfg := &config.Config{
		Port:      "8080",
		DBPath:    ":memory:",
		GinMode:   "test",
		JWTSecret: "test-jwt-secret-key-12345",
	}

	router := routes.SetupRouter(db, cfg)
	return router, db
}

func TestAuthRegisterAndLogin(t *testing.T) {
	router, _ := setupTestRouter(t)

	// 1. Register a new user
	registerPayload := models.RegisterInput{
		Username: "bookworm",
		Password: "password123",
	}
	body, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created on register, got %d. Body: %s", w.Code, w.Body.String())
	}

	var regResp models.AuthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &regResp); err != nil {
		t.Fatalf("failed to parse register response: %v", err)
	}
	if regResp.Token == "" {
		t.Errorf("expected token in register response, got empty")
	}
	if regResp.User.Username != "bookworm" {
		t.Errorf("expected username 'bookworm', got '%s'", regResp.User.Username)
	}

	// 2. Duplicate registration should fail with 409
	reqDup, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	reqDup.Header.Set("Content-Type", "application/json")
	wDup := httptest.NewRecorder()
	router.ServeHTTP(wDup, reqDup)

	if wDup.Code != http.StatusConflict {
		t.Errorf("expected status 409 Conflict on duplicate register, got %d", wDup.Code)
	}

	// 3. Login with invalid password
	loginBadPass := models.LoginInput{
		Username: "bookworm",
		Password: "wrongpassword",
	}
	badBody, _ := json.Marshal(loginBadPass)
	reqBad, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(badBody))
	reqBad.Header.Set("Content-Type", "application/json")
	wBad := httptest.NewRecorder()
	router.ServeHTTP(wBad, reqBad)

	if wBad.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 Unauthorized for bad password, got %d", wBad.Code)
	}

	// 4. Login with correct credentials
	loginGood := models.LoginInput{
		Username: "bookworm",
		Password: "password123",
	}
	goodBody, _ := json.Marshal(loginGood)
	reqGood, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(goodBody))
	reqGood.Header.Set("Content-Type", "application/json")
	wGood := httptest.NewRecorder()
	router.ServeHTTP(wGood, reqGood)

	if wGood.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK on login, got %d", wGood.Code)
	}

	var loginResp models.AuthResponse
	if err := json.Unmarshal(wGood.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("failed to parse login response: %v", err)
	}
	if loginResp.Token == "" {
		t.Errorf("expected token in login response, got empty")
	}

	// 5. Test /api/v1/auth/me with the login token
	reqMe, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	reqMe.Header.Set("Authorization", fmt.Sprintf("Bearer %s", loginResp.Token))
	wMe := httptest.NewRecorder()
	router.ServeHTTP(wMe, reqMe)

	if wMe.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK on /auth/me, got %d", wMe.Code)
	}
}

func TestBookRoutesProtectedByMiddleware(t *testing.T) {
	router, _ := setupTestRouter(t)

	// 1. Request /api/v1/books without token -> expect 401 Unauthorized
	reqNoAuth, _ := http.NewRequest(http.MethodGet, "/api/v1/books", nil)
	wNoAuth := httptest.NewRecorder()
	router.ServeHTTP(wNoAuth, reqNoAuth)

	if wNoAuth.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 Unauthorized for request without token, got %d", wNoAuth.Code)
	}

	// 2. Request /api/v1/books with invalid token -> expect 401 Unauthorized
	reqBadToken, _ := http.NewRequest(http.MethodGet, "/api/v1/books", nil)
	reqBadToken.Header.Set("Authorization", "Bearer invalid.token.payload")
	wBadToken := httptest.NewRecorder()
	router.ServeHTTP(wBadToken, reqBadToken)

	if wBadToken.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 Unauthorized for invalid token, got %d", wBadToken.Code)
	}

	// 3. Register user and get valid token
	regBody, _ := json.Marshal(models.RegisterInput{
		Username: "author1",
		Password: "password123",
	})
	reqReg, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(regBody))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	router.ServeHTTP(wReg, reqReg)

	var authResp models.AuthResponse
	_ = json.Unmarshal(wReg.Body.Bytes(), &authResp)
	token := authResp.Token

	// 4. Request /api/v1/books with valid token -> expect 200 OK
	reqAuth, _ := http.NewRequest(http.MethodGet, "/api/v1/books", nil)
	reqAuth.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	wAuth := httptest.NewRecorder()
	router.ServeHTTP(wAuth, reqAuth)

	if wAuth.Code != http.StatusOK {
		t.Errorf("expected status 200 OK for valid token, got %d. Body: %s", wAuth.Code, wAuth.Body.String())
	}

	// 5. Create a book with valid token
	createBody, _ := json.Marshal(models.LookupBook{
		Isbn:        "9780134190440",
		Title:       "The Go Programming Language",
		Author:      "Alan A. A. Donovan, Brian W. Kernighan",
		Description: "An authoritative resource for Go programming.",
		Pages:       380,
	})
	reqCreate, _ := http.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewBuffer(createBody))
	reqCreate.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	router.ServeHTTP(wCreate, reqCreate)

	if wCreate.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created for book creation, got %d. Body: %s", wCreate.Code, wCreate.Body.String())
	}

	// Verify book list contains the newly created book
	reqList, _ := http.NewRequest(http.MethodGet, "/api/v1/books", nil)
	reqList.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	wList := httptest.NewRecorder()
	router.ServeHTTP(wList, reqList)

	if wList.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK for book list, got %d", wList.Code)
	}

	var listResp struct {
		Data  []models.Book `json:"data"`
		Count int           `json:"count"`
	}
	if err := json.Unmarshal(wList.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("failed to parse book list: %v", err)
	}

	if listResp.Count != 1 {
		t.Errorf("expected 1 book in list, got %d", listResp.Count)
	}
	if listResp.Data[0].UserId != authResp.User.ID {
		t.Errorf("expected book user_id %d, got %d", authResp.User.ID, listResp.Data[0].UserId)
	}
}

func TestAdminBootupSeeding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file:mem_seed_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Book{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	cfg := &config.Config{
		Port:          "8080",
		DBPath:        ":memory:",
		GinMode:       "test",
		JWTSecret:     "test-secret",
		AdminUsername: "admin",
		AdminPassword: "admin123",
	}

	// 1. Seed admin user
	if err := database.SeedAdminUser(db, cfg); err != nil {
		t.Fatalf("failed to seed admin user: %v", err)
	}

	// Verify admin exists
	var admin models.User
	if err := db.Where("username = ?", "admin").First(&admin).Error; err != nil {
		t.Fatalf("admin user not found: %v", err)
	}

	// Verify admin can login via router
	router := routes.SetupRouter(db, cfg)
	loginBody, _ := json.Marshal(models.LoginInput{
		Username: "admin",
		Password: "admin123",
	})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected admin login to return 200, got %d", w.Code)
	}

	// 2. Calling SeedAdminUser again when users exist should not create duplicates
	if err := database.SeedAdminUser(db, cfg); err != nil {
		t.Fatalf("second call to SeedAdminUser failed: %v", err)
	}
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
}

func TestChangePassword(t *testing.T) {
	router, _ := setupTestRouter(t)

	// 1. Register a user
	regBody, _ := json.Marshal(models.RegisterInput{
		Username: "changepassuser",
		Password: "oldpassword123",
	})
	reqReg, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(regBody))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	router.ServeHTTP(wReg, reqReg)

	var regResp models.AuthResponse
	_ = json.Unmarshal(wReg.Body.Bytes(), &regResp)
	token := regResp.Token

	// 2. Try to change password without auth token -> 401
	chgBody, _ := json.Marshal(models.ChangePasswordInput{
		OldPassword: "oldpassword123",
		NewPassword: "newpassword456",
	})
	reqNoAuth, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewBuffer(chgBody))
	reqNoAuth.Header.Set("Content-Type", "application/json")
	wNoAuth := httptest.NewRecorder()
	router.ServeHTTP(wNoAuth, reqNoAuth)

	if wNoAuth.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for unauthorized change-password, got %d", wNoAuth.Code)
	}

	// 3. Try to change password with wrong old password -> 401
	chgWrongOld, _ := json.Marshal(models.ChangePasswordInput{
		OldPassword: "incorrectpassword",
		NewPassword: "newpassword456",
	})
	reqWrong, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewBuffer(chgWrongOld))
	reqWrong.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	reqWrong.Header.Set("Content-Type", "application/json")
	wWrong := httptest.NewRecorder()
	router.ServeHTTP(wWrong, reqWrong)

	if wWrong.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for wrong old password, got %d", wWrong.Code)
	}

	// 4. Try to change password with too short new password -> 400
	chgShort, _ := json.Marshal(models.ChangePasswordInput{
		OldPassword: "oldpassword123",
		NewPassword: "123",
	})
	reqShort, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewBuffer(chgShort))
	reqShort.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	reqShort.Header.Set("Content-Type", "application/json")
	wShort := httptest.NewRecorder()
	router.ServeHTTP(wShort, reqShort)

	if wShort.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for short new password, got %d", wShort.Code)
	}

	// 5. Change password successfully
	reqSuccess, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewBuffer(chgBody))
	reqSuccess.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	reqSuccess.Header.Set("Content-Type", "application/json")
	wSuccess := httptest.NewRecorder()
	router.ServeHTTP(wSuccess, reqSuccess)

	if wSuccess.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for password change, got %d. Body: %s", wSuccess.Code, wSuccess.Body.String())
	}

	// Also test the PUT /api/v1/auth/password alias
	chgBody2, _ := json.Marshal(models.ChangePasswordInput{
		OldPassword: "newpassword456",
		NewPassword: "evennewerpassword789",
	})
	reqPut, _ := http.NewRequest(http.MethodPut, "/api/v1/auth/password", bytes.NewBuffer(chgBody2))
	reqPut.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	reqPut.Header.Set("Content-Type", "application/json")
	wPut := httptest.NewRecorder()
	router.ServeHTTP(wPut, reqPut)

	if wPut.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for PUT /password, got %d. Body: %s", wPut.Code, wPut.Body.String())
	}

	// 6. Previous passwords should now fail login
	oldLogin, _ := json.Marshal(models.LoginInput{
		Username: "changepassuser",
		Password: "oldpassword123",
	})
	reqOldLogin, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(oldLogin))
	reqOldLogin.Header.Set("Content-Type", "application/json")
	wOldLogin := httptest.NewRecorder()
	router.ServeHTTP(wOldLogin, reqOldLogin)

	if wOldLogin.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 logging in with old password, got %d", wOldLogin.Code)
	}

	// 7. Latest password should succeed
	newLogin, _ := json.Marshal(models.LoginInput{
		Username: "changepassuser",
		Password: "evennewerpassword789",
	})
	reqNewLogin, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(newLogin))
	reqNewLogin.Header.Set("Content-Type", "application/json")
	wNewLogin := httptest.NewRecorder()
	router.ServeHTTP(wNewLogin, reqNewLogin)

	if wNewLogin.Code != http.StatusOK {
		t.Fatalf("expected 200 logging in with new password, got %d", wNewLogin.Code)
	}
}
