package service

import (
	"context"
	"time"
	"github.com/steel-feel/prac/internal/domain"
	"github.com/steel-feel/prac/internal/port"
)

type healthService struct {
	db        pingable
	startTime time.Time
}

// Interface for what we need from the DB to check health
type pingable interface {
	PingContext(ctx context.Context) error
}

// NewHealthService creates a new health service.
func NewHealthService(db pingable) port.HealthService {
	return &healthService{
		db:        db,
		startTime: time.Now(),
	}
}

func (s *healthService) Check(ctx context.Context) (*domain.HealthStatus, error) {
	dbStatus := "ok"
	if err := s.db.PingContext(ctx); err != nil {
		dbStatus = "error"
	}

	uptime := time.Since(s.startTime).String()

	status := "fully operational"
	if dbStatus != "ok" {
		status = "degraded"
	}

	return &domain.HealthStatus{
		Status:   status,
		DBStatus: dbStatus,
		Uptime:   uptime,
	}, nil
}
