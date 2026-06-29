package grpc

import (
	"google.golang.org/grpc"

	documentv1 "github.com/steel-feel/prac/proto/document/v1"
	"github.com/steel-feel/prac/internal/port"
)

// NewServer creates and configures a new gRPC server.
func NewServer(docSvc port.DocumentService) *grpc.Server {
	s := grpc.NewServer()

	docServer := NewDocumentServer(docSvc)
	documentv1.RegisterDocumentServiceServer(s, docServer)

	return s
}
