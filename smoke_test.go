package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/steel-feel/prac/api/generated"
	mygrpc "github.com/steel-feel/prac/internal/adapter/primary/grpc"
	myhttp "github.com/steel-feel/prac/internal/adapter/primary/http"
	"github.com/steel-feel/prac/internal/adapter/secondary/coinlore"
	"github.com/steel-feel/prac/internal/adapter/secondary/coinpaprika"
	"github.com/steel-feel/prac/internal/adapter/secondary/facilitator"
	"github.com/steel-feel/prac/internal/adapter/secondary/sqlite"
	"github.com/steel-feel/prac/internal/service"
	documentv1 "github.com/steel-feel/prac/proto/document/v1"
)

func setupTestServer(t *testing.T) (*echo.Echo, string) {
	// Use an in-memory database for testing
	dbPath := ":memory:"
	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	if err := sqlite.RunMigrations(db.DB, os.DirFS("."), "migrations", "up"); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	fac := facilitator.New(facilitator.Config{})

	healthSvc := service.NewHealthService(db)
	docRepo := sqlite.NewDocumentRepository(db)
	docSvc := service.NewDocumentService(docRepo)
	accessRepo := sqlite.NewAccessLogRepository(db)
	
	userRepo := sqlite.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, "test-jwt-secret")

	coinpaprikaURL := os.Getenv("COINPAPRIKA_URL")
	if coinpaprikaURL == "" {
		coinpaprikaURL = "https://api.coinpaprika.com"
	}
	coinpaprikaClient := coinpaprika.NewClient(coinpaprikaURL, nil)

	coinloreURL := os.Getenv("COINLORE_URL")
	if coinloreURL == "" {
		coinloreURL = "https://api.coinlore.net"
	}
	coinloreClient := coinlore.NewClient(coinloreURL, nil)

	priceSvc := service.NewPriceService(coinpaprikaClient, coinloreClient)

	healthHandler := myhttp.NewHealthHandler(healthSvc)
	docHandler := myhttp.NewDocumentHandler(docSvc)
	authHandler := myhttp.NewAuthHandler(userSvc)
	priceHandler := myhttp.NewPriceHandler(priceSvc)

	e := echo.New()
	
	handlers := &myhttp.Handlers{
		Health:   healthHandler,
		Document: docHandler,
		Auth:     authHandler,
		Price:    priceHandler,
	}

	myhttp.RegisterRoutes(e, handlers, fac, docSvc, accessRepo, "test-jwt-secret")

	return e, dbPath
}

func TestHealthEndpoint(t *testing.T) {
	e, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "fully operational" {
		t.Errorf("expected status fully operational, got %s", response["status"])
	}
	if response["db_status"] != "ok" {
		t.Errorf("expected db_status ok, got %s", response["db_status"])
	}
}

func TestDocumentEndpoint_NoPayment(t *testing.T) {
	e, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents/doc-001", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusPaymentRequired {
		t.Errorf("expected status 402, got %d", rec.Code)
	}

	if rec.Header().Get("PAYMENT-REQUIRED") == "" {
		t.Error("expected PAYMENT-REQUIRED header")
	}
}

func TestDocumentEndpoint_WithPayment(t *testing.T) {
	e, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/documents/doc-001", nil)
	req.Header.Set("PAYMENT-SIGNATURE", "valid-dummy-signature")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	if rec.Header().Get("PAYMENT-RESPONSE") == "" {
		t.Error("expected PAYMENT-RESPONSE header")
	}

	var doc generated.Document
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if doc.Id != "doc-001" {
		t.Errorf("expected doc-001, got %s", doc.Id)
	}
}

func TestLogin_Success(t *testing.T) {
	e, _ := setupTestServer(t)

	body := `{"username": "admin", "password": "password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	var resp generated.LoginResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	e, _ := setupTestServer(t)

	body := `{"username": "admin", "password": "wrongpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestCreateDocument_NoAuth(t *testing.T) {
	e, _ := setupTestServer(t)

	body := `{"title": "New Doc", "content": "Hello World", "price": 150}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/documents", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestCreateDocument_Success(t *testing.T) {
	e, _ := setupTestServer(t)

	// 1. Login to get token
	body := `{"username": "admin", "password": "password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	var loginResp generated.LoginResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &loginResp)
	token := loginResp.Token

	// 2. Create document using the token
	docBody := `{"title": "Authenticated Doc", "content": "Super secret auth content", "price": 500}`
	reqDoc := httptest.NewRequest(http.MethodPost, "/api/v1/documents", strings.NewReader(docBody))
	reqDoc.Header.Set("Content-Type", "application/json")
	reqDoc.Header.Set("Authorization", "Bearer "+token)
	recDoc := httptest.NewRecorder()
	e.ServeHTTP(recDoc, reqDoc)

	if recDoc.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d. Body: %s", recDoc.Code, recDoc.Body.String())
	}

	var doc generated.Document
	if err := json.Unmarshal(recDoc.Body.Bytes(), &doc); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if doc.Id == "" {
		t.Error("expected non-empty document ID")
	}
	if doc.Title != "Authenticated Doc" {
		t.Errorf("expected title 'Authenticated Doc', got '%s'", doc.Title)
	}

	// 3. Retrieve the created document without payment (should get 402 first)
	reqGet := httptest.NewRequest(http.MethodGet, "/api/v1/documents/"+doc.Id, nil)
	recGet := httptest.NewRecorder()
	e.ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusPaymentRequired {
		t.Errorf("expected status 402 for retrieving doc without payment, got %d", recGet.Code)
	}
}

func TestGrpcDocumentService(t *testing.T) {
	// Initialize database
	dbPath := ":memory:"
	db, err := sqlite.NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer db.Close()

	if err := sqlite.RunMigrations(db.DB, os.DirFS("."), "migrations", "up"); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	docRepo := sqlite.NewDocumentRepository(db)
	docSvc := service.NewDocumentService(docRepo)

	// Create and start gRPC server on a random local port
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	s := mygrpc.NewServer(docSvc)
	go func() {
		if err := s.Serve(lis); err != nil {
			// server might stop on close, ignore error
		}
	}()
	defer s.GracefulStop()

	// Connect gRPC client to the server
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := documentv1.NewDocumentServiceClient(conn)

	// Test ListDocuments
	listResp, err := client.ListDocuments(context.Background(), &documentv1.ListDocumentsRequest{})
	if err != nil {
		t.Fatalf("ListDocuments failed: %v", err)
	}
	if len(listResp.Documents) != 2 {
		t.Errorf("expected 2 documents, got %d", len(listResp.Documents))
	}

	// Test GetDocument
	getResp, err := client.GetDocument(context.Background(), &documentv1.GetDocumentRequest{Id: "doc-001"})
	if err != nil {
		t.Fatalf("GetDocument failed: %v", err)
	}
	if getResp.Id != "doc-001" {
		t.Errorf("expected doc-001, got %s", getResp.Id)
	}
	if getResp.Title != "Getting Started with Tempo" {
		t.Errorf("expected 'Getting Started with Tempo', got %s", getResp.Title)
	}
	if getResp.PriceUsd != 100 {
		t.Errorf("expected price 100, got %d", getResp.PriceUsd)
	}
}

func TestGetEthereumPrice_Integration_Success(t *testing.T) {
	mockPaprikaResponse := `{"quotes": {"USD": {"price": 2045.12}}}`
	mockPaprikaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockPaprikaResponse))
	}))
	defer mockPaprikaServer.Close()

	mockLoreResponse := `[{"price_usd":"2055.12"}]`
	mockLoreServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockLoreResponse))
	}))
	defer mockLoreServer.Close()

	os.Setenv("COINPAPRIKA_URL", mockPaprikaServer.URL)
	os.Setenv("COINLORE_URL", mockLoreServer.URL)
	defer func() {
		os.Unsetenv("COINPAPRIKA_URL")
		os.Unsetenv("COINLORE_URL")
	}()

	e, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/price/eth", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	var resp generated.EthereumPriceResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedPrice := 2050.12 // Average of 2045.12 and 2055.12
	if resp.Price != expectedPrice {
		t.Errorf("expected price %f, got %f", expectedPrice, resp.Price)
	}
}
