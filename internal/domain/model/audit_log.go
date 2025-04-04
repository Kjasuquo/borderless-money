package model

import (
	"operation-borderless/internal/domain/dto"
	"time"
)

type AuditLog struct {
	Models
	UserID    string
	IPAddress string
	Device    string
	Country   string
	Browser   string
	Timestamp time.Time
}

func (a *AuditLog) ToAuditLogDTO() dto.AuditLog {

	out := dto.AuditLog{
		Models: dto.Models{
			ID:        a.ID,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		},
		UserID:    a.UserID,
		IPAddress: a.IPAddress,
		Device:    a.Device,
		Country:   a.Country,
		Browser:   a.Browser,
		Timestamp: a.Timestamp,
	}

	return out
}

func FromAuditLogDTO(a dto.AuditLog) AuditLog {
	response := AuditLog{
		Models: Models{
			ID:        a.ID,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		},
		UserID:    a.UserID,
		IPAddress: a.IPAddress,
		Device:    a.Device,
		Country:   a.Country,
		Browser:   a.Browser,
		Timestamp: a.Timestamp,
	}

	return response
}
