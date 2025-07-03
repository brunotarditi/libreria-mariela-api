package services

import (
	"libreria/common"
	"libreria/models"
	"libreria/repositories"
	"libreria/responses"

	"gorm.io/gorm"
)

type DashboardService interface {
	GetData() (responses.DashboardResponse, error)
}

type dashboardService struct {
	db            *gorm.DB
	dashboardRepo repositories.DashboardRepository
	supplierOps   *common.GormOperations[models.Supplier]
	customerOps   *common.GormOperations[models.Customer]
	productOps    *common.GormOperations[models.Product]
}

func NewDashboardService(db *gorm.DB, dashboardRepo repositories.DashboardRepository, supplierOps *common.GormOperations[models.Supplier], customerOps *common.GormOperations[models.Customer], productOps *common.GormOperations[models.Product]) DashboardService {
	return &dashboardService{
		db:            db,
		dashboardRepo: dashboardRepo,
		supplierOps:   supplierOps,
		customerOps:   customerOps,
		productOps:    productOps,
	}
}

func (s *dashboardService) GetData() (responses.DashboardResponse, error) {
	options := common.QueryOptions{}
	var response responses.DashboardResponse
	totalCustomers, _ := s.customerOps.Count(options)
	totalProducts, _ := s.productOps.Count(options)
	totalSuppliers, _ := s.supplierOps.Count(options)

	auditLogs, err := s.dashboardRepo.GetAuditLog()
	if err != nil {
		return responses.DashboardResponse{}, err
	}

	response = responses.DashboardResponse{
		TotalProducts:    totalProducts,
		TotalSuppliers:   totalSuppliers,
		TotalClients:     totalCustomers,
		RecentActivities: auditLogs,
	}

	return response, nil
}
