package usecase

import (
	"context"
	"time"

	"github.com/rashaev/todo-app/internal/entity"
	"github.com/rashaev/todo-app/internal/repository"
)

type todoUseCase struct {
	todoRepo repository.TodoRepository
}

func NewTodoUseCase(r repository.TodoRepository) TodoUseCase {
	return &todoUseCase{todoRepo: r}
}

func (uc *todoUseCase) CreateTodo(ctx context.Context, title, description string) error {
	todo := &entity.Todo{
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return uc.todoRepo.Create(ctx, todo)
}

func (uc *todoUseCase) GetAllTodos(ctx context.Context) ([]entity.Todo, error) {
	return uc.todoRepo.List(ctx)
}

func (uc *todoUseCase) GetTodoByID(ctx context.Context, id int64) (entity.Todo, error) {
	return uc.todoRepo.GetByID(ctx, id)
}

func (uc *todoUseCase) UpdateTodo(ctx context.Context, todo *entity.Todo) error {
	todo.UpdatedAt = time.Now()
	return uc.todoRepo.Update(ctx, todo)
}

func (uc *todoUseCase) DeleteTodo(ctx context.Context, id int64) error {
	return uc.todoRepo.Delete(ctx, id)
}

func (uc *todoUseCase) MarkTodoDone(ctx context.Context, id int64) error {
	return uc.todoRepo.MarkDone(ctx, id)
}
