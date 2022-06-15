package repositories

import "gorm.io/gorm"

type RepositoryBase interface {
	BeginTransaction() *gorm.DB
}

type RepositoryBaseImpl struct {
	db *gorm.DB
}

func (rb *RepositoryBaseImpl) BeginTransaction() *gorm.DB {
	return rb.db.Begin()
}
