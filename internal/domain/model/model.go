package model

import (
	"errors"
	"gorm.io/gorm"
	"time"

	"github.com/google/uuid"
)

type Models struct {
	ID        string    `sql:"type:uuid; default:uuid_generate_v4();size:100; not null"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (u *Models) BeforeCreate(tx *gorm.DB) (err error) {

	u.ID = uuid.New().String()
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	if u.ID == "" {
		err = errors.New("can't save invalid data")
	}
	return
}
