package main

import (
	"flag"
	"fmt"
	"os"

	"net"

	"github.com/labstack/echo/v5"

	"github.com/steel-feel/prac/internal/adapter/primary/http"
	"github.com/steel-feel/prac/internal/adapter/primary/http/middleware"
	mygrpc "github.com/steel-feel/prac/internal/adapter/primary/grpc"
	"github.com/steel-feel/prac/internal/adapter/secondary/facilitator"
	"github.com/steel-feel/prac/internal/adapter/secondary/sqlite"
	"github.com/steel-feel/prac/internal/service"
)

func main() {
	port := flag.String("port", "9000", "port that will be exposed")
	grpcPort := flag.String("grpc-port", "9001", "port for gRPC server")
	dbPath := flag.String("db-path", "data.db", "path to sqlite database")
	flag.Parse()

	// 1. Initialize Secondary Adapters (Database & Facilitator)
	db, err := sqlite.NewDB(*dbPath)
	if err != nil {
		fmt.Printf("failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	facCfg := facilitator.Config{
		RPCURL:       "https://testnet.tempo.xyz",
		ChainID:      "tempo-testnet",
		USDCContract: "0x123", // Dummy
		PayTo:        "0x456", // Dummy
	}
	fac := facilitator.New(facCfg)

	// 2. Initialize Services (Domain logic)
	healthSvc := service.NewHealthService(db)
	docRepo := sqlite.NewDocumentRepository(db)
	docSvc := service.NewDocumentService(docRepo)
	accessRepo := sqlite.NewAccessLogRepository(db)

	// 3. Initialize Primary Adapters (HTTP Handlers & gRPC)
	healthHandler := http.NewHealthHandler(healthSvc)
	docHandler := http.NewDocumentHandler(docSvc)
	handlers := &http.Handlers{
		Health:   healthHandler,
		Document: docHandler,
	}
	
	grpcServer := mygrpc.NewServer(docSvc)

	// 4. Setup Echo and Register Routes
	e := echo.New()
	e.Use(middleware.ObservabilityMiddleware())
	http.RegisterRoutes(e, handlers, fac, docSvc, accessRepo)

	// 5. Start gRPC Server
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *grpcPort))
		if err != nil {
			fmt.Printf("failed to listen on gRPC port: %v\n", err)
			return
		}
		fmt.Printf("Starting gRPC server on :%s\n", *grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Printf("failed to serve gRPC: %v\n", err)
		}
	}()

	// 6. Start HTTP Server
	if err := e.Start(fmt.Sprintf(":%s", *port)); err != nil {
		e.Logger.Error("server error", "error", err)
	}
}
