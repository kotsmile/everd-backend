package todolist_presentation

import (
	"context"
	"net/http"

	access_domain "github.com/kotsmile/everd-backend/internal/domain/access"
	todolist_domain "github.com/kotsmile/everd-backend/internal/domain/todolist"
	"github.com/kotsmile/everd-backend/internal/util"
)

type TodolistHandler struct {
	util.HTTPHandler
	service todolist_domain.TodolistService
}

func NewTodolistHandler(service todolist_domain.TodolistService) *TodolistHandler {
	return &TodolistHandler{
		service: service,
	}
}

type TodoResponse struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Done    bool   `json:"done"`
}

type GetTodolistResponse struct {
	Todos []TodoResponse `json:"todos"`
}

func (h *TodolistHandler) GetTodolist(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	userID, ok := ctx.Value("userID").(access_domain.UserID)
	if !ok {
		return util.
			NewHTTPError("unauthorized").
			WithStatus(http.StatusUnauthorized).
			WithErrorMessage("user id is not provided")
	}

	todolist, err := h.service.GetTodolist(ctx, userID)
	if err != nil {
		return err
	}

	todolistPF := todolist.PF()

	todolistResponse := GetTodolistResponse{
		Todos: make([]TodoResponse, len(todolistPF.Todos)),
	}

	for i, todo := range todolistPF.Todos {
		todolistResponse.Todos[i] = TodoResponse{
			ID:      todo.ID.Int(),
			Title:   todo.Title,
			Comment: todo.Comment,
			Done:    todo.Done,
		}
	}

	return h.OkJSON(w, todolistResponse)
}

type PostTodoRequest struct {
	Title string `json:"title"`
}

type PostTodoResponse struct{}

func (h *TodolistHandler) PostTodo(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	userID, ok := ctx.Value("userID").(access_domain.UserID)
	if !ok {
		return util.
			NewHTTPError("unauthorized").
			WithStatus(http.StatusUnauthorized).
			WithErrorMessage("user id is not provided")
	}

	var todoRequest PostTodoRequest
	if err := h.ReadJSON(w, r, &todoRequest); err != nil {
		return util.
			NewHTTPError("invalid request").
			WithStatus(http.StatusBadRequest).
			WithError(err)
	}

	if err := h.service.AddTodo(ctx, userID, todoRequest.Title); err != nil {
		return err
	}

	return h.OkJSON(w, PostTodoResponse{})
}
