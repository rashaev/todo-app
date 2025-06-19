package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/rashaev/todo-app/internal/entity"
	"github.com/rashaev/todo-app/internal/repository"
)

type todoPostgresRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) repository.TodoRepository {
	return &todoPostgresRepository{db: db}
}

func (r *todoPostgresRepository) Create(ctx context.Context, todo *entity.Todo) error {
	query := `INSERT INTO todos (title, description, created_at, updated_at, completed) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		todo.Title,
		todo.Description,
		todo.CreatedAt,
		todo.UpdatedAt,
		todo.Completed,
	).Scan(&todo.ID)

	return err
}

func (r *todoPostgresRepository) List(ctx context.Context) ([]entity.Todo, error) {
	query := `SELECT id, title, description, created_at, updated_at, completed FROM todos`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []entity.Todo
	for rows.Next() {
		var todo entity.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.Completed,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *todoPostgresRepository) GetByID(ctx context.Context, id int64) (entity.Todo, error) {
	query := `SELECT id, title, description, completed, created_at, updated_at 
	          FROM todos WHERE id = $1`
	var todo entity.Todo
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	return todo, err
}

func (r *todoPostgresRepository) Update(ctx context.Context, todo *entity.Todo) error {
	query := `UPDATE todos SET title=$1, description=$2, updated_at=$3 WHERE id=$4`
	_, err := r.db.ExecContext(
		ctx,
		query,
		todo.Title,
		todo.Description,
		time.Now(),
		todo.ID,
	)
	return err
}

func (r *todoPostgresRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM todos WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *todoPostgresRepository) MarkDone(ctx context.Context, id int64) error {
	query := `UPDATE todos SET completed=true WHERE id=$1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
