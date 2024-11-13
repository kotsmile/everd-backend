package todolist_model

import (
	"fmt"
)

var ErrTodoID = fmt.Errorf("%w: todo id", Err)

type TodoID uint

var NilTodoID TodoID

func NewTodoID(id int) (TodoID, error) {
	if id < 0 {
		return NilTodoID, ErrTodoID
	}

	return TodoID(uint(id)), nil
}

func (id TodoID) Equal(other TodoID) bool {
	return id == other
}

func (id TodoID) Int() int {
	return int(id)
}
