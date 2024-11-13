package app

import (
	todolist_domain "github.com/kotsmile/everd-backend/internal/app/domain/todolist"
	todolist_infrastructure "github.com/kotsmile/everd-backend/internal/app/domain/todolist/infrastructure"
	todolist_presentation "github.com/kotsmile/everd-backend/internal/app/domain/todolist/presentation"
	"github.com/kotsmile/everd-backend/internal/app/infrastructure/storage"
	"github.com/kotsmile/everd-backend/internal/util"
)

func run() {
	todoRepo := todolist_infrastructure.NewPostrgesTodoRepository(nil)
	todolistRepo := todolist_infrastructure.NewPostrgesTodolistRepository(nil)

	txFactory := storage.NewSQLTransactionFactory(nil)

	todolistService := todolist_domain.NewTodoService(
		txFactory,
		todolistRepo,
		todoRepo,
	)

	util.NewLogger(false)
	todolist_presentation.NewTodolistHandler(todolistService)
}
