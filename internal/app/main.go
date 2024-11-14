package app

import (
	"github.com/gorilla/mux"
	access_handler "github.com/kotsmile/everd-backend/internal/app/domain/access/handler"
	todolist_domain "github.com/kotsmile/everd-backend/internal/app/domain/todolist"
	todolist_handler "github.com/kotsmile/everd-backend/internal/app/domain/todolist/handler"
	todolist_infrastructure "github.com/kotsmile/everd-backend/internal/app/domain/todolist/infrastructure"
	"github.com/kotsmile/everd-backend/internal/app/infrastructure/storage"
	"github.com/kotsmile/everd-backend/internal/util"
)

func run() {
	logger := util.NewLogger(false)
	apiHelper := util.NewApiHelper(logger)

	// repositories
	todoRepo := todolist_infrastructure.NewPostrgesTodoRepository(nil)
	todolistRepo := todolist_infrastructure.NewPostrgesTodolistRepository(nil)

	// infrastructure
	txFactory := storage.NewSQLTransactionFactory(nil)

	// services
	todolistService := todolist_domain.NewTodoService(
		txFactory,
		todolistRepo,
		todoRepo,
	)

	r := mux.NewRouter()
	access := access_handler.NewAccessHandler(apiHelper)

	todolist := todolist_handler.NewTodolistHandler(todolistService, apiHelper)
	r.HandleFunc("/todolist", apiHelper.Wrapper(todolist.GetTodolist, access.AuthMiddlerware)).Methods("GET")
	r.HandleFunc("/todolist/todo", apiHelper.Wrapper(todolist.PostTodo, access.AuthMiddlerware)).Methods("POST")
}
