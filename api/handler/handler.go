package handler

import (
	"auth-service/storage/postgres"
	"database/sql"
)

type Handler struct {
	UserRepo *postgres.UserRepo
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		UserRepo: postgres.NewUserRepo(db),
	}
}