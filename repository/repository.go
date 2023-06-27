package repository

import (
	"github.com/jinzhu/gorm"
)

type Repository struct {
	db *gorm.DB
}

var _defaultRepository *Repository

func GetRepository() *Repository {
	return _defaultRepository
}

func InitRepository(db *gorm.DB) error {
	_defaultRepository = NewRepository(db)
	return nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (p *Repository) GetDB() *gorm.DB {
	return p.db
}

type TRepository struct {
	*Repository
}

func NewTRepository(db *gorm.DB) *TRepository {
	return &TRepository{Repository: NewRepository(db)}
}

func (p *TRepository) Begin() error {

	p.Repository.db = p.Repository.db.Begin()

	return p.Repository.db.Error
}

func (p *TRepository) Commit() error {
	return p.Repository.db.Commit().Error
}

func (p *TRepository) Rollback() error {
	return p.Repository.db.Rollback().Error
}
