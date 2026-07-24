package main

import (
	"database/sql"

	"orderly/internal/pkg/token"
	userrepo "orderly/internal/repo/postgres/user"
	authusecase "orderly/internal/usecase/auth"
)

type usecaseSet struct {
	auth *authusecase.Usecase
}

func initUsecase(db *sql.DB) usecaseSet {
	userRepo := userrepo.NewUserRepo(db)

	return usecaseSet{
		auth: authusecase.New(userRepo, token.Generate),
	}
}
