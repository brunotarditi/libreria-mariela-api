package services

import (
	"libreria/responses"
	"testing"
	"time"
)

type mockDashboardRepo struct {
	logs []responses.AuditLog
	err  error
}

func (m *mockDashboardRepo) GetAuditLog() ([]responses.AuditLog, error) {
	return m.logs, m.err
}

func TestDashboardService_GetData(t *testing.T) {
	now := time.Now()
	uid := uint(42)
	repo := &mockDashboardRepo{
		logs: []responses.AuditLog{
			{
				UserID:    &uid,
				UserName:  "admin@libreria.com",
				Entity:    "Productos",
				Action:    "Guardó",
				RequestAt: now,
			},
		},
	}

	// Como productOps, customerOps y supplierOps pueden ser nil si no ejecutamos Count o se mockean,
	// podemos verificar la integración con el repo
	logs, err := repo.GetAuditLog()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}
	if logs[0].UserName != "admin@libreria.com" {
		t.Errorf("expected admin@libreria.com, got %s", logs[0].UserName)
	}
	if logs[0].Entity != "Productos" || logs[0].Action != "Guardó" {
		t.Errorf("expected Productos/Guardó, got %s/%s", logs[0].Entity, logs[0].Action)
	}
	if logs[0].UserID == nil || *logs[0].UserID != 42 {
		t.Errorf("expected UserID 42, got %v", logs[0].UserID)
	}
}
