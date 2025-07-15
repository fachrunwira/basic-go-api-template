package storage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fachrunwira/basic-go-api-template/lib/logger"
)

type storage struct {
	root string
	path string
	err  error
}

var AppLogger *log.Logger = logger.SetLogger("./storage/log/app.log")

func Init() *storage {
	dir, err := os.Getwd()
	if err != nil {
		AppLogger.Fatalf("error while getting path: %v", err)
	}

	return &storage{
		root: dir,
	}
}

func (s *storage) PublicPath() *storage {
	if s.err != nil {
		return s
	}

	path := filepath.Join(s.root, "public", "storage")
	s.path = path
	return s
}

func (s *storage) StoragePath() *storage {
	if s.err != nil {
		return s
	}

	s.path = filepath.Join(s.root, "storage", "app")
	return s
}

func (s *storage) Directory(dir ...string) *storage {
	if s.err != nil {
		return s
	}

	if len(dir) == 0 {
		s.err = fmt.Errorf("please input the directory path")
		return s
	}

	if s.path == "" {
		s.err = fmt.Errorf("must specified using PublicPath() or StoragePath() first")
	}

	diskPath := filepath.Join(s.path, filepath.Join(dir...))
	s.path = diskPath
	return s
}

func (s storage) Exists(path string) (bool, error) {
	if s.path == "" {
		return false, fmt.Errorf("must specified using PublicPath() or StoragePath() first")
	}

	if s.err != nil {
		return false, s.err
	}

	fullPath := filepath.Join(s.path, path)

	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (s *storage) SaveFile(filename string) *storage {
	if s.err != nil {
		return s
	}

	if s.path == "" {
		s.err = fmt.Errorf("must specified using PublicPath() or StoragePath() first")
		return s
	}

	fullPath := filepath.Join(s.path, filename)
	f, err := os.Create(fullPath)
	if err != nil {
		s.err = err
		return s
	}
	defer f.Close()

	f.WriteString("dummy content")

	fmt.Printf("Saved file at: %v", fullPath)
	return s
}

func (s *storage) Get(path string) *storage {
	if s.err != nil {
		return s
	}

	if s.path == "" {
		s.err = fmt.Errorf("must specified using PublicPath() or StoragePath() first")
		return s
	}

	fullPath := filepath.Join(s.path, path)
	s.path = fullPath
	return s
}

func (s storage) Csv() ([]string, error) {
	if s.err != nil {
		return nil, s.err
	}

	if s.path == "" {
		return nil, fmt.Errorf("must specified using PublicPath() or StoragePath() first")
	}

	file, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("the specified file is not found")
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []string
	scanFile := bufio.NewScanner(file)
	for scanFile.Scan() {
		data = append(data, scanFile.Text())
	}

	return data, nil
}

func (s storage) Url() (*string, error) {
	if s.err != nil {
		return nil, s.err
	}

	if s.path == "" {
		return nil, fmt.Errorf("must specified using PublicPath() or StoragePath() first")
	}

	return &s.path, nil
}

func (s *storage) Error() error {
	return s.err
}
