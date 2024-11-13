package todlist_infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	access_domain "github.com/kotsmile/everd-backend/internal/domain/access"
	todolist_domain "github.com/kotsmile/everd-backend/internal/domain/todolist"
	todolist_model "github.com/kotsmile/everd-backend/internal/domain/todolist/model"
	"github.com/kotsmile/everd-backend/internal/infrastructure/storage"
	"github.com/kotsmile/everd-backend/internal/util"
)

type PostrgesTodoRepository struct {
	db *sql.DB
}

var _ todolist_domain.TodoRepository = (*PostrgesTodoRepository)(nil)

type TodoDTO struct {
	ID        int
	Title     string
	Comment   string
	Done      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func toTodoDTO(todo todolist_model.Todo) TodoDTO {
	todoPF := todo.PF()
	return TodoDTO{
		ID:        todoPF.ID.Int(),
		Title:     todoPF.Title,
		Comment:   todoPF.Comment,
		Done:      todoPF.Done,
		CreatedAt: todoPF.CreatedAt,
		UpdatedAt: todoPF.UpdatedAt,
	}
}

func fromTodoDTO(todoDTO TodoDTO) (todolist_model.Todo, error) {
	todoID, err := todolist_model.NewTodoID(todoDTO.ID)
	if err != nil {
		return todolist_model.Todo{}, err
	}

	todo, err := todolist_model.NewTodoFromDB(
		todoID,
		todoDTO.Title,
		todoDTO.Comment,
		todoDTO.Done,
		todoDTO.CreatedAt,
		todoDTO.UpdatedAt,
	)
	if err != nil {
		return todolist_model.Todo{}, err
	}

	return todo, nil
}

func (r *PostrgesTodoRepository) NextID(
	ctx context.Context,
	tx util.Transaction,
) (todolist_model.TodoID, error) {
	exec, err := storage.GetDBExecutor(tx, r.db)
	if err != nil {
		return todolist_model.NilTodoID, err
	}

	var id int
	if err := exec.QueryRow("select max(id) from todos").Scan(&id); err != nil {
		return todolist_model.NilTodoID, err
	}

	todoID, err := todolist_model.NewTodoID(id + 1)
	if err != nil {
		return todolist_model.NilTodoID, err
	}

	return todoID, nil
}

type PostrgesTodolistRepository struct {
	db *sql.DB
}

var _ todolist_domain.TodolistRepository = (*PostrgesTodolistRepository)(nil)

func (r *PostrgesTodolistRepository) Get(
	ctx context.Context,
	userID access_domain.UserID,
	tx util.Transaction,
) (*todolist_model.Todolist, error) {
	exec, err := storage.GetDBExecutor(tx, r.db)
	if err != nil {
		return nil, err
	}

	rows, err := exec.Query(`select todos.* from todolist 
	                         left join todos 
	                         on todolist.todo_id = todos.id 
	                         where user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []todolist_model.Todo

	for rows.Next() {
		var todoDTO TodoDTO
		if err := rows.Scan(
			&todoDTO.ID,
			&todoDTO.Title,
			&todoDTO.Comment,
			&todoDTO.Done,
			&todoDTO.CreatedAt,
			&todoDTO.UpdatedAt,
		); err != nil {
			return nil, err
		}

		todo, err := fromTodoDTO(todoDTO)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	todolist, err := todolist_model.NewTodolist(userID, todos)
	if err != nil {
		return nil, err
	}

	return todolist, nil
}

func (r *PostrgesTodolistRepository) Save(
	ctx context.Context,
	todolist *todolist_model.Todolist,
	tx util.Transaction,
) error {
	exec, commit, rollaback, err := storage.GetTxOrCreateTx(ctx, tx, r.db)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = errors.Join(rollaback(), err)
		} else {
			err = commit()
		}
	}()

	// insert all todos
	stmtTodoInsert, err := exec.Prepare(`insert into todos
		(id, title, comment, done, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return err
	}
	defer stmtTodoInsert.Close()

	// isnert all todos
	todolistPF := todolist.PF()
	for _, todo := range todolistPF.Todos {
		if _, err := stmtTodoInsert.Exec(
			todo.ID,
			todo.Title,
			todo.Comment,
			todo.Done,
			todo.CreatedAt,
			todo.UpdatedAt,
		); err != nil {
			return err
		}
	}

	// delete all todos from todolist for user
	if _, err := exec.Exec(
		`delete from todolist where user_id = $1`,
		todolistPF.UserID,
	); err != nil {
		return err
	}

	// insert all todos from todolist for user
	stmtTodolistInsert, err := exec.Prepare(`insert into todolist
		(user_id, todo_id)
		values ($1, $2)`)
	if err != nil {
		return err
	}
	defer stmtTodolistInsert.Close()

	for _, todo := range todolistPF.Todos {
		if _, err := stmtTodolistInsert.Exec(
			todolistPF.UserID,
			todo.ID,
		); err != nil {
			return err
		}
	}

	return nil
}
