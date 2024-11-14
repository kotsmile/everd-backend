package todolist_model

import (
	"fmt"

	access_domain "github.com/kotsmile/everd-backend/internal/app/domain/access"
)

var ErrNotFound = fmt.Errorf("%w: todo is not found", Err)

type Todolist struct {
	userID access_domain.UserID
	todos  []Todo
}

func NewTodolistEmpty(userID access_domain.UserID) *Todolist {
	return &Todolist{
		userID: userID,
		todos:  []Todo{},
	}
}

func NewTodolist(userID access_domain.UserID, todos []Todo) (*Todolist, error) {
	todolist := &Todolist{
		userID: userID,
		todos:  todos,
	}

	if err := todolist.Validate(); err != nil {
		return nil, err
	}

	return todolist, nil
}

type TodolistPF struct {
	UserID access_domain.UserID
	Todos  []TodoPF
}

func (l *Todolist) PF() TodolistPF {
	todoPFs := make([]TodoPF, len(l.todos))
	for i, todo := range l.todos {
		todoPFs[i] = todo.PF()
	}
	return TodolistPF{
		UserID: l.userID,
		Todos:  todoPFs,
	}
}

func (l *Todolist) Validate() error {
	for _, todo := range l.todos {
		if err := todo.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (l *Todolist) AddTodo(id TodoID, title string) {
	todo := NewTodo(id, title)
	l.todos = append(l.todos, todo)
}

func (l *Todolist) CompleteTodo(todoID TodoID) error {
	for _, todo := range l.todos {
		if todo.id.Equal(todoID) {
			if err := todo.Complete(); err != nil {
				return err
			}
			return nil
		}
	}

	return ErrNotFound
}

func (l *Todolist) UncompleteTodo(todoID TodoID) error {
	for i, todo := range l.todos {
		if todo.id.Equal(todoID) {
			if err := todo.Uncomplete(); err != nil {
				return err
			}

			l.todos[i] = todo
			return nil
		}
	}

	return ErrNotFound
}

func (l *Todolist) ChangeTitle(todoID TodoID, title string) error {
	for _, todo := range l.todos {
		if todo.id.Equal(todoID) {
			todo.ChangeTitle(title)
			return nil
		}
	}

	return ErrNotFound
}

func (l *Todolist) ChangeComment(todoID TodoID, comment string) error {
	for _, todo := range l.todos {
		if todo.id.Equal(todoID) {
			todo.ChangeComment(comment)
			return nil
		}
	}

	return ErrNotFound
}
