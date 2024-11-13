package todolist_model

import (
	"errors"
	"fmt"
	"time"
)

var (
	Err                 = errors.New("todolist")
	ErrIsCompleted      = fmt.Errorf("%w: todo is completed", Err)
	ErrNotCompleted     = fmt.Errorf("%w: todo is not completed", Err)
	ErrTitleIsEmpty     = fmt.Errorf("%w: title is empty", Err)
	ErrTitleIsTooLong   = fmt.Errorf("%w: title is too long", Err)
	ErrCommentIsTooLong = fmt.Errorf("%w: comment is too long", Err)
)

const (
	MaxTitleLength   = 100
	MaxCommentLength = 1000
)

type Todo struct {
	id TodoID

	title   string
	comment string
	done    bool

	createdAt time.Time
	updatedAt time.Time
}

func NewTodo(id TodoID, title string) Todo {
	return Todo{
		id:        id,
		title:     title,
		comment:   "",
		done:      false,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
}

func NewTodoFromDB(
	id TodoID,
	title string,
	comment string,
	done bool,
	createdAt time.Time,
	updatedAt time.Time,
) (Todo, error) {
	todo := Todo{
		id:        id,
		title:     title,
		comment:   comment,
		done:      done,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	if err := todo.Validate(); err != nil {
		return Todo{}, err
	}

	return todo, nil
}

type TodoPF struct {
	ID        TodoID
	Title     string
	Comment   string
	Done      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Todo) PF() TodoPF {
	return TodoPF{
		ID:        t.id,
		Title:     t.title,
		Comment:   t.comment,
		Done:      t.done,
		CreatedAt: t.createdAt,
		UpdatedAt: t.updatedAt,
	}
}

func (t *Todo) Validate() error {
	if t.title == "" {
		return ErrTitleIsEmpty
	}

	if len(t.title) > MaxTitleLength {
		return ErrTitleIsTooLong
	}

	if t.comment != "" && len(t.comment) > MaxCommentLength {
		return ErrCommentIsTooLong
	}

	return nil
}

func (t *Todo) UpdateTitle(title string) {
	t.title = title
	t.updatedAt = time.Now()
}

func (t *Todo) ChangeComment(comment string) {
	t.comment = comment
	t.updatedAt = time.Now()
}

func (t *Todo) Complete() error {
	if t.done {
		return ErrIsCompleted
	}

	t.done = true
	t.updatedAt = time.Now()

	return nil
}

func (t *Todo) Uncomplete() error {
	if !t.done {
		return ErrNotCompleted
	}

	t.done = false
	t.updatedAt = time.Now()

	return nil
}
