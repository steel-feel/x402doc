package grpc

import (
	"context"

	documentv1 "github.com/steel-feel/prac/proto/document/v1"
	"github.com/steel-feel/prac/internal/port"
)

type DocumentServer struct {
	documentv1.UnimplementedDocumentServiceServer
	svc port.DocumentService
}

func NewDocumentServer(svc port.DocumentService) *DocumentServer {
	return &DocumentServer{svc: svc}
}

func (s *DocumentServer) GetDocument(ctx context.Context, req *documentv1.GetDocumentRequest) (*documentv1.DocumentResponse, error) {
	doc, err := s.svc.GetDocument(ctx, req.Id)
	if err != nil {
		return nil, err // In a real app, translate domain errors to gRPC status codes
	}
	if doc == nil {
		return nil, nil // Return Not Found status
	}

	return &documentv1.DocumentResponse{
		Id:       doc.ID,
		Title:    doc.Title,
		Content:  doc.Content,
		PriceUsd: doc.PriceUSD,
	}, nil
}

func (s *DocumentServer) ListDocuments(ctx context.Context, req *documentv1.ListDocumentsRequest) (*documentv1.ListDocumentsResponse, error) {
	docs, err := s.svc.ListDocuments(ctx)
	if err != nil {
		return nil, err
	}

	var res []*documentv1.DocumentResponse
	for _, doc := range docs {
		res = append(res, &documentv1.DocumentResponse{
			Id:       doc.ID,
			Title:    doc.Title,
			Content:  doc.Content,
			PriceUsd: doc.PriceUSD,
		})
	}

	return &documentv1.ListDocumentsResponse{Documents: res}, nil
}
