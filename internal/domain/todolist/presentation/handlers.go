package todolist_presentation

import (
	"context"
	"encoding/json"
	"net/http"

	access_domain "github.com/kotsmile/everd-backend/internal/domain/access"
	todolist_domain "github.com/kotsmile/everd-backend/internal/domain/todolist"
)

type TodolistHandler struct {
	service todolist_domain.TodolistService
}

func NewTodolistHandler(service todolist_domain.TodolistService) *TodolistHandler {
	return &TodolistHandler{
		service: service,
	}
}

type HttpError struct{}

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)

	return err
}

func OkJSON(w http.ResponseWriter, data any) {
	if err := WriteJSON(w, http.StatusOK, JsonResponse{
		Data: data,
	}); err != nil {
		// TODO: log error
	}
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload JsonResponse

	payload.Error = true
	payload.Message = err.Error()

	if err := WriteJSON(w, statusCode, payload); err != nil {
		// TODO: log error
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

func (h *TodolistHandler) GetTodolist(ctx context.Context, w http.ResponseWriter, r *http.Request) *HttpError {
	userID, ok := ctx.Value("userID").(access_domain.UserID)
	if !ok {
		return &HttpError{}
	}

	todolist, err := h.service.GetTodolist(ctx, userID)
	if err != nil {
		return &HttpError{}
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

	OkJSON(w, todolistResponse)

	return nil
}

type PostTodoRequest struct {
	Title string `json:"title"`
}

type PostTodoResponse struct{}

func (h *TodolistHandler) PostTodo(ctx context.Context, w http.ResponseWriter, r *http.Request) *HttpError {
	userID, ok := ctx.Value("userID").(access_domain.UserID)
	if !ok {
		return &HttpError{}
	}

	var todoRequest PostTodoRequest
	if err := ReadJSON(w, r, &todoRequest); err != nil {
		return &HttpError{}
	}

	if err := h.service.AddTodo(ctx, userID, todoRequest.Title); err != nil {
		return &HttpError{}
	}

	OkJSON(w, PostTodoResponse{})
	return nil
}
