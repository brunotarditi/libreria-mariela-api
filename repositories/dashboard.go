package repositories

import (
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
	err := r.db.Joins("INNER JOIN users u ON audit_logs.user_id = u.id").
		Where("audit_logs.route IN (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			"/api/v1/suppliers", "/api/v1/suppliers/:id",
			"/api/v1/brands", "/api/v1/brands/list", "/api/v1/brands/:id",
			"/api/v1/products", "/api/v1/products/import", "/api/v1/products/:id",
			"/api/v1/customers", "/api/v1/customers/:id",
			"/api/v1/categories", "/api/v1/categories/list", "/api/v1/categories/:id").
		Where("audit_logs.method IN (?, ?, ?)", "POST", "PUT", "DELETE").
		Order("audit_logs.request_at DESC").
		Limit(5).
		Select(`
            u.username AS user_name,
            CASE 
                WHEN audit_logs.route = '/api/v1/products/import' OR audit_logs.route = '/api/v1/products' OR audit_logs.route = '/api/v1/products/:id' THEN 'Productos'
                WHEN audit_logs.route = '/api/v1/brands' OR audit_logs.route = '/api/v1/brands/list' OR audit_logs.route = '/api/v1/brands/:id' THEN 'Marcas'
                WHEN audit_logs.route = '/api/v1/customers' OR audit_logs.route = '/api/v1/customers/:id' THEN 'Clientes'
                WHEN audit_logs.route = '/api/v1/suppliers' OR audit_logs.route = '/api/v1/suppliers/:id' THEN 'Proveedores'
                WHEN audit_logs.route = '/api/v1/categories' OR audit_logs.route = '/api/v1/categories/list' OR audit_logs.route = '/api/v1/categories/:id' THEN 'Categorías'
                ELSE ''
            END AS entity,
            CASE 
                WHEN audit_logs.method = 'POST' THEN 'Guardó'
                WHEN audit_logs.method = 'PUT' THEN 'Actualizó'
                WHEN audit_logs.method = 'DELETE' THEN 'Eliminó'
                ELSE ''
            END AS action,
            audit_logs.request_at AS request_at
        `).
		Find(&auditLogs).Error

	return auditLogs, err
}
