package todolist_domain

import (
	"context"
	"errors"
	"fmt"

	access_domain "github.com/kotsmile/everd-backend/internal/domain/access"
	todolist_model "github.com/kotsmile/everd-backend/internal/domain/todolist/model"
	"github.com/kotsmile/everd-backend/internal/util"
)

//go:generate mockgen -package todolist_domain -source=./service.go -destination=./mock.go *

var ErrTodolistNotFound = fmt.Errorf("%w: todolist not found", todolist_model.Err)

type TodoRepository interface {
	NextID(ctx context.Context, tx util.Transaction) (todolist_model.TodoID, error)
}

type TodolistRepository interface {
	Get(ctx context.Context, userID access_domain.UserID, tx util.Transaction) (*todolist_model.Todolist, error)
	Save(ctx context.Context, todolist *todolist_model.Todolist, tx util.Transaction) error
}

type TodolistService struct {
	txFactory    util.TransactionFactory
	todolistRepo TodolistRepository
	todoRepo     TodoRepository
}

func NewTodoService(
	txFactory util.TransactionFactory,
	todolistRepo TodolistRepository,
	todoRepo TodoRepository,
) *TodolistService {
	return &TodolistService{
		txFactory:    txFactory,
		todolistRepo: todolistRepo,
		todoRepo:     todoRepo,
	}
}

func (s *TodolistService) GetTodolist(
	ctx context.Context,
	userID access_domain.UserID,
) (*todolist_model.Todolist, error) {
	return s.getOrCreateTodolist(ctx, userID, nil)
}

func (s *TodolistService) AddTodo(
	ctx context.Context,
	userID access_domain.UserID,
	title string,
) error {
	return s.txFactory.WithTransaction(ctx, func(tx util.Transaction) error {
		list, err := s.getOrCreateTodolist(ctx, userID, tx)
		if err != nil {
			return err
		}

		todoID, err := s.todoRepo.NextID(ctx, tx)
		if err != nil {
			return err
		}

		list.AddTodo(todoID, title)
		if err := list.Validate(); err != nil {
			return err
		}

		if err := s.todolistRepo.Save(ctx, list, tx); err != nil {
			return err
		}

		return nil
	})
}

func (s *TodolistService) CompleteTodo(
	ctx context.Context,
	userID access_domain.UserID,
	todoID todolist_model.TodoID,
) error {
	return s.txFactory.WithTransaction(ctx, func(tx util.Transaction) error {
		list, err := s.getOrCreateTodolist(ctx, userID, tx)
		if err != nil {
			return err
		}

		if err := list.CompleteTodo(todoID); err != nil {
			return err
		}

		if err := list.Validate(); err != nil {
			return err
		}

		if err := s.todolistRepo.Save(ctx, list, tx); err != nil {
			return err
		}

		return nil
	})
}

func (s *TodolistService) UncompleteTodo(
	ctx context.Context,
	userID access_domain.UserID,
	todoID todolist_model.TodoID,
) error {
	return s.txFactory.WithTransaction(ctx, func(tx util.Transaction) error {
		list, err := s.getOrCreateTodolist(ctx, userID, tx)
		if err != nil {
			return err
		}

		if err := list.UncompleteTodo(todoID); err != nil {
			return err
		}

		if err := list.Validate(); err != nil {
			return err
		}

		if err := s.todolistRepo.Save(ctx, list, tx); err != nil {
			return err
		}

		return nil
	})
}

func (s *TodolistService) ChangeComment(
	ctx context.Context,
	userID access_domain.UserID,
	todoID todolist_model.TodoID,
	comment string,
) error {
	return s.txFactory.WithTransaction(ctx, func(tx util.Transaction) error {
		list, err := s.getOrCreateTodolist(ctx, userID, tx)
		if err != nil {
			return err
		}

		if err := list.ChangeComment(todoID, comment); err != nil {
			return err
		}

		if err := list.Validate(); err != nil {
			return err
		}

		if err := s.todolistRepo.Save(ctx, list, tx); err != nil {
			return err
		}

		return nil
	})
}

func (s *TodolistService) getOrCreateTodolist(
	ctx context.Context,
	userID access_domain.UserID,
	tx util.Transaction,
) (*todolist_model.Todolist, error) {
	list, err := s.todolistRepo.Get(ctx, userID, tx)
	if err != nil {
		if errors.Is(err, ErrTodolistNotFound) {
			list = todolist_model.NewTodolistEmpty(userID)

			if err := s.todolistRepo.Save(ctx, list, tx); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return list, nil
}
