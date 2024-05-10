package gorm

import (
	"context"
	"gorm.io/gorm"
)

type Repository[T any] ChainDao

func NewRepository[T any](ctx context.Context, db *gorm.DB) *Repository[T] {
	return (*Repository[T])(NewChainDao(ctx, db))
}

func (r *Repository[T]) Create(t *T) error {
	return r.DB.Create(t).Error
}

func (r *Repository[T]) Retrieve(id int) (*T, error) {
	var t T
	err := r.DB.First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository[T]) Update(t *T) error {
	return r.DB.Updates(&t).Error
}

func (r *Repository[T]) Delete(id int) error {
	//(*T)(nil)
	var t T
	return r.DB.Delete(&t, id).Error
}

func GetById[T any](db *gorm.DB, id int) (*T, error) {
	t := new(T)
	err := db.First(t, id).Error
	return t, err
}

type Repository2[T any] gorm.DB

func (db *Repository2[T]) GetById(id int) (*T, error) {
	t := new(T)
	err := (*gorm.DB)(db).First(t, id).Error
	return t, err
}

type Repository3[T any] struct {
	gorm.DB
}

func (db *Repository3[T]) GetById(id int) (*T, error) {
	t := new(T)
	err := db.First(t, id).Error
	return t, err
}
