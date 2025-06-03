package repository

import (
	"context"

	"github.com/rashaev/todo-app/internal/entity"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *entity.Todo) error
	List(ctx context.Context) ([]entity.Todo, error)
	GetByID(ctx context.Context, id int64) (entity.Todo, error)
	Update(ctx context.Context, todo *entity.Todo) error
	Delete(ctx context.Context, id int64) error
}
