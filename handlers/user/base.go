package user

import "log"

type userHandler struct {
	AppLogger *log.Logger
}

func NewUserHandler(logger *log.Logger) *userHandler {
	return &userHandler{
		AppLogger: logger,
	}
}
