package usecases

import (
	"context"
)

type TenantUseCases struct {
	tenantService TenantService
}

type TenantService interface {
	Create(ctx context.Context, tenantId string) error
}

func NewTenantUseCases(tService TenantService) *TenantUseCases {
	return &TenantUseCases{
		tenantService: tService,
	}
}

func (t *TenantUseCases) Create(ctx context.Context, tenantId string) error {
	return t.tenantService.Create(ctx, tenantId)
}
