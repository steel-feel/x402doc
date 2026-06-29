package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/steel-feel/prac/api/generated"
	mygrpc "github.com/steel-feel/prac/internal/adapter/primary/grpc"
	myhttp "github.com/steel-feel/prac/internal/adapter/primary/http"
	"github.com/steel-feel/prac/internal/adapter/primary/http/middleware"
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

	var doc generated.Document
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if doc.Id != "doc-001" {
		t.Errorf("expected doc-001, got %s", doc.Id)
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
