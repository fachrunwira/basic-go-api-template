package auth

import "log"

type authHandler struct {
	AppLogger *log.Logger
}

func NewAuthHandler(logger *log.Logger) *authHandler {
	return &authHandler{
		AppLogger: logger,
	}
}
