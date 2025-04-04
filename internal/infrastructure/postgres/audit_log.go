package postgres

import (
	"context"
	"operation-borderless/internal/domain/model"
)

func (d *database) CreateAuditLogs(ctx context.Context, auditLog model.AuditLog) error {

	err := d.db.WithContext(ctx).Create(&auditLog).Error
	if err != nil {
		return err
	}

	return nil
}
