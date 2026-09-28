package repositories

import (
	"libreria/models"
	"libreria/responses"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetAuditLog() ([]responses.AuditLog, error)
}

type dashboardRepositoryRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepositoryRepository{db: db}
}

func (r *dashboardRepositoryRepository) GetAuditLog() ([]responses.AuditLog, error) {
	var auditLogs []responses.AuditLog
	err := r.db.Model(&models.AuditLog{}).
		Where("audit_logs.route LIKE ? OR audit_logs.route LIKE ? OR audit_logs.route LIKE ? OR audit_logs.route LIKE ? OR audit_logs.route LIKE ?",
			"/api/v1/products%", "/api/v1/brands%", "/api/v1/customers%", "/api/v1/suppliers%", "/api/v1/categories%").
		Where("audit_logs.method IN (?, ?, ?, ?)", "POST", "PUT", "PATCH", "DELETE").
		Order("audit_logs.request_at DESC").
		Limit(5).
		Select(`
            COALESCE(
                NULLIF(audit_logs.user_name, ''),
                CASE 
                    WHEN audit_logs.user_id IS NOT NULL THEN 'Usuario #' || CAST(audit_logs.user_id AS VARCHAR)
                    ELSE 'Sistema'
                END
            ) AS user_name,
            audit_logs.user_id AS user_id,
            CASE 
                WHEN audit_logs.route LIKE '/api/v1/products%' THEN 'Productos'
                WHEN audit_logs.route LIKE '/api/v1/brands%' THEN 'Marcas'
                WHEN audit_logs.route LIKE '/api/v1/customers%' THEN 'Clientes'
                WHEN audit_logs.route LIKE '/api/v1/suppliers%' THEN 'Proveedores'
                WHEN audit_logs.route LIKE '/api/v1/categories%' THEN 'Categorías'
                ELSE ''
            END AS entity,
            CASE 
                WHEN audit_logs.method = 'POST' THEN 'Guardó'
                WHEN audit_logs.method = 'PUT' OR audit_logs.method = 'PATCH' THEN 'Actualizó'
                WHEN audit_logs.method = 'DELETE' THEN 'Eliminó'
                ELSE ''
            END AS action,
            audit_logs.request_at AS request_at
        `).
		Find(&auditLogs).Error

	return auditLogs, err
}
