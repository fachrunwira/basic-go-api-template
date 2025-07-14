package storage

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fachrunwira/basic-go-api-template/lib/logger"
)

type storage struct {
	root string
	app  string
	path string
	err  error
}

var AppLogger *log.Logger = logger.SetLogger("./storage/log/app.log")

func Init() *storage {
	dir, err := os.Getwd()
	if err != nil {
		AppLogger.Fatalf("Failed to current path: %v\n", err)
	}

	return &storage{
		root: dir,
		app:  filepath.Join(dir, "storage", "app"),
	}
}
func (s *storage) Directory(dir string) *storage {
	if s.err != nil {
		return s
	}

	diskPath := filepath.Join(s.app, "public", dir)
	if _, err := os.Stat(diskPath); os.IsNotExist(err) {
		if err := os.MkdirAll(diskPath, os.ModePerm); err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	s.path = diskPath
	return s
}

func (s *storage) CheckFileExists(path string) (bool, error) {
	if s.path == "" {
		s.path = filepath.Join(s.app, "public", path)
	}

	if s.err != nil {
		return false, s.err
	}

	info, err := os.Stat(s.path)
	if os.IsNotExist(err) {
		return false, s.err
	} else if err != nil {
		return false, s.err
	}

	return !info.IsDir(), nil
}

func (s *storage) CheckDirectoryExists(path string) (bool, error) {
	if s.path == "" {
		s.path = filepath.Join(s.app, "public", path)
	}

	if s.err != nil {
		return false, s.err
	}

	info, err := os.Stat(s.path)
	if os.IsNotExist(err) {
		return false, s.err
	} else if err != nil {
		return false, s.err
	}

	return info.IsDir(), nil
}
