package usecase

import (
	"context"

	"github.com/rashaev/todo-app/internal/entity"
)

type TodoUseCase interface {
	CreateTodo(ctx context.Context, title, description string) error
	GetAllTodos(ctx context.Context) ([]entity.Todo, error)
	GetTodoByID(ctx context.Context, id int64) (entity.Todo, error)
	UpdateTodo(ctx context.Context, todo *entity.Todo) error
	DeleteTodo(ctx context.Context, id int64) error
}
