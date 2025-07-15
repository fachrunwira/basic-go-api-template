package files

import "log"

type filesHandler struct {
	AppLogger *log.Logger
}

func NewFilesHandler(log *log.Logger) *filesHandler {
	return &filesHandler{
		AppLogger: log,
	}
}
