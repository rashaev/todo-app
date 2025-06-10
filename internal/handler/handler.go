package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rashaev/todo-app/internal/entity"
	"github.com/rashaev/todo-app/internal/usecase"
	"github.com/rashaev/todo-app/pkg/logger"
)

type TodoHandlers struct {
	todoUseCase usecase.TodoUseCase
	logger      logger.Logger
}

func NewTodoHandlers(todoUseCase usecase.TodoUseCase, logger logger.Logger) *TodoHandlers {
	return &TodoHandlers{
		todoUseCase: todoUseCase,
		logger:      logger,
	}
}

func (h *TodoHandlers) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo entity.Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Error("Unable to decode body", "error", err)
		return
	}

	err := h.todoUseCase.CreateTodo(r.Context(), todo.Title, todo.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error("Unable to create todo", "error", err)
		return
	}

	h.logger.Info("CreateTodo", "path", r.URL.Path, "http_method", r.Method, "proto", r.Proto, "remote_addr", r.RemoteAddr)
	w.WriteHeader(http.StatusCreated)
}

func (h *TodoHandlers) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.todoUseCase.GetAllTodos(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error("GetAllTodos", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
	h.logger.Info("GetAllTodos", "path", r.URL.Path, "http_method", r.Method, "proto", r.Proto, "remote_addr", r.RemoteAddr)

}

func (h *TodoHandlers) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Error("Unable to get todo by id", "error", err)
		return
	}
	todo, err := h.todoUseCase.GetTodoByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		h.logger.Error("Todo not found", "error", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
	h.logger.Info("GetTodoByID", "path", r.URL.Path, "http_method", r.Method, "proto", r.Proto, "remote_addr", r.RemoteAddr)

}

func (h *TodoHandlers) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Error("Unable to parse todo id", "error", err)
		return
	}

	var todo entity.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Error("Unable to decode todo during update", "error", err)
		return
	}
	todo.ID = id

	err = h.todoUseCase.UpdateTodo(r.Context(), &todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error("Unable to update todo", "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	h.logger.Info("UpdateTodo", "path", r.URL.Path, "http_method", r.Method, "proto", r.Proto, "remote_addr", r.RemoteAddr, "todo_id", id)

}

func (h *TodoHandlers) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Error("Unable to parse todo id", "error", err)
		return
	}

	err = h.todoUseCase.DeleteTodo(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error("Unable to delete todo", "error", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	h.logger.Info("DeleteTodo", "path", r.URL.Path, "http_method", r.Method, "proto", r.Proto, "remote_addr", r.RemoteAddr, "todo_id", id)
}
