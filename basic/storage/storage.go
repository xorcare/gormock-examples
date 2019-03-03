package storage

import (
	"github.com/xorcare/gormock/linker"
	"time"
)

// Model is gorm database model.
type Model struct {
	ID        uint64     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`
}

// Storage model.
type Storage struct {
	db *linker.DB
}

// New returns instance.
func New(db *linker.DB) *Storage {
	return &Storage{db: db}
}

// FindByID method returns one model from the table is
// filtered by the passed primary key.
func (repo Storage) FindByID(id uint64) (model *Model, err error) {
	defer func() {
		if err != nil {
			model = nil
		}
	}()
	model = &Model{ID: id}
	return model, repo.db.First(model).Error
}

// FindAll method returns all models from the table.
func (repo Storage) FindAll() (models []Model, err error) {
	return models, repo.db.Find(&models).Error
}

// Count method returns number of all models in the table.
func (repo Storage) Count() (result int64, err error) {
	err = repo.db.Model(Model{}).Count(&result).Error
	return result, err
}

// DeleteByID method deletes one model from the table is
// filtered by the passed primary key.
func (repo Storage) DeleteByID(id uint64) error {
	return repo.db.Where("id = ?", id).Delete(Model{}).Error
}

// DeleteAll method deletes all models from the table.
func (repo Storage) DeleteAll() error {
	return repo.db.Delete(Model{}).Error
}

// Save method saves an passed model to the table
func (repo Storage) Save(model *Model) error {
	return repo.db.Save(model).Error
}
