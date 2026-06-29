package domain

// HealthStatus represents the health of the service.
type HealthStatus struct {
	Status   string `json:"status"`
	DBStatus string `json:"db_status"`
	Uptime   string `json:"uptime"`
}
