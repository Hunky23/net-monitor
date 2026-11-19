package main

import "context"

type App struct {
	ctx context.Context
	//userUsecase *usecase.UserUsecase
}

func NewApp() *App {
	return &App{}
}
