package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v5"
	documentv1 "github.com/steel-feel/prac/proto/document/v1"
	"github.com/steel-feel/prac/internal/adapter/primary/http/middleware"
	myhttp "github.com/steel-feel/prac/internal/adapter/primary/http"
	"github.com/steel-feel/prac/internal/adapter/secondary/facilitator"
	"github.com/steel-feel/prac/internal/adapter/secondary/sqlite"
	"github.com/steel-feel/prac/internal/service"
)

func setupTestServer(t *testing.T) (*echo.Echo, string) {
	// Use an in-memory database for testing
	dbPath := ":memory:"
	db, err := sqlite.NewDB(dbPath, os.DirFS("."))
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	fac := facilitator.New(facilitator.Config{})

	healthSvc := service.NewHealthService(db)
	docRepo := sqlite.NewDocumentRepository(db)
	docSvc := service.NewDocumentService(docRepo)
	accessRepo := sqlite.NewAccessLogRepository(db)

	healthHandler := myhttp.NewHealthHandler(healthSvc)
	docHandler := myhttp.NewDocumentHandler(docSvc)

	e := echo.New()
	e.Use(middleware.ObservabilityMiddleware())

	api := e.Group("/api/v1")
	api.GET("/health", healthHandler.Check)

	x402Mw := middleware.X402Middleware(docSvc, fac)
	accessLogMw := middleware.AccessLogMiddleware(accessRepo)
	api.GET("/documents/:id", docHandler.GetDoc, x402Mw, accessLogMw)

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

	var doc documentv1.DocumentResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if doc.Id != "doc-001" {
		t.Errorf("expected doc-001, got %s", doc.Id)
	}
}
