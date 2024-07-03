package services

import "context"

type TenantService struct {
	tenantRepository TenantRepository
}

type TenantRepository interface {
	Create(ctx context.Context, tenantId string) error
}

func NewTeanantService(tenantRepository TenantRepository) *TenantService {
	return &TenantService{
		tenantRepository: tenantRepository,
	}
}

func (t *TenantService) Create(ctx context.Context, tenantId string) error {
	return t.tenantRepository.Create(ctx, tenantId)
}
