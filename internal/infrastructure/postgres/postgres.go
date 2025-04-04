package postgres

import "gorm.io/gorm"

type database struct {
	db *gorm.DB
}

func NewDB(db *gorm.DB) *database {
	return &database{
		db: db,
	}
}
