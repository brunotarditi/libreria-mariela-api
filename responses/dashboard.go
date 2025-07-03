package responses

import "time"

type DashboardResponse struct {
	TotalProducts    int64      `json:"total_products"`
	TotalSuppliers   int64      `json:"total_suppliers"`
	TotalClients     int64      `json:"total_clients"`
	RecentActivities []AuditLog `json:"recent_activities"`
}

type AuditLog struct {
	UserName  string    `json:"user_name"`
	Entity    string    `json:"entity"`
	Action    string    `json:"action"`
	RequestAt time.Time `json:"request_at"`
}
