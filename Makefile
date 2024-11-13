MIGRATION_DIR=migrations
GOOSE_CMD=goose

migrate-create:
ifndef NAME
	$(error NAME is required for creating a new migration)
endif
	$(GOOSE_CMD) -dir $(MIGRATION_DIR) create $(NAME) sql
