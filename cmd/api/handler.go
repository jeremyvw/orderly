package main

import (
	authhandler "orderly/internal/handler/http/auth"
)

type handlerSet struct {
	auth *authhandler.Handler
}

func initHandler(uc usecaseSet) handlerSet {
	return handlerSet{
		auth: authhandler.New(uc.auth),
	}
}
