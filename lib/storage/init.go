package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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
		s.err = fmt.Errorf("must specify using PublicPath() or StoragePath() first")
	}

	diskPath := filepath.Join(s.path, filepath.Join(dir...))
	s.path = diskPath
	return s
}

func (s *storage) Exists(path string) (bool, error) {
	if s.path == "" {
		return false, fmt.Errorf("must specify using PublicPath() or StoragePath() first")
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

func (s *storage) StoreFile(filename string, source any) error {
	if s.err != nil {
		return s.err
	}

	if s.path == "" {
		return fmt.Errorf("must specify using PublicPath() or StoragePath() first")
	}

	fullPath := filepath.Join(s.path, filename)
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	switch src := source.(type) {
	case *multipart.FileHeader:
		file, err := src.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := io.Copy(f, file); err != nil {
			return err
		}
	case string:
		if _, err := f.WriteString(src); err != nil {
			return err
		}
	case []byte:
		if _, err := f.Write(src); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported source type: %v", src)
	}

	return nil
}

func (s *storage) Get(path string) *storage {
	if s.err != nil {
		return s
	}

	if s.path == "" {
		s.err = fmt.Errorf("must specify using PublicPath() or StoragePath() first")
		return s
	}

	fullPath := filepath.Join(s.path, path)
	s.path = fullPath
	return s
}

func (s *storage) Json() (interface{}, error) {
	if s.err != nil {
		return nil, s.err
	}

	if s.path == "" {
		return nil, fmt.Errorf("must specify using PublicPath() or StoragePath() first")
	}

	file, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("the specify file is not found")
	} else if err != nil {
		return nil, err
	}

	var obj interface{}
	if err := json.Unmarshal(file, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *storage) SaveJsonInChunk(rows int, outputFile string) error {
	if s.err != nil {
		return s.err
	}

	if s.path == "" {
		return fmt.Errorf("must specify using PublicPath() or StoragePath() first")
	}

	if outputFile == "" {
		return fmt.Errorf("outputFile must be specify")
	}

	info, err := os.Stat(s.path)
	if os.IsNotExist(err) {
		return fmt.Errorf("the specify file is not found")
	} else if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("given path is not a file")
	}

	file, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	t, err := decoder.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != '[' {
		return fmt.Errorf("expected JSON array at root")
	}

	var chunk []json.RawMessage
	chunkIndex := 1
	total := 0

	for decoder.More() {
		var obj json.RawMessage
		if err := decoder.Decode(&obj); err != nil {
			return err
		}
		chunk = append(chunk, obj)
		total++

		if len(chunk) == rows {
			if err := s.saveChunk(chunk, chunkIndex, outputFile); err != nil {
				return err
			}
			chunk = chunk[:0]
			chunkIndex++
		}
	}

	if len(chunk) > 0 {
		if err := s.saveChunk(chunk, chunkIndex, outputFile); err != nil {
			return err
		}
	}

	return nil
}

func (s *storage) saveChunk(chunk []json.RawMessage, index int, outFile string) error {
	if len(chunk) == 0 {
		return nil
	}

	outDir := filepath.Dir(s.path)
	_, err := os.Stat(outDir)
	if os.IsNotExist(err) {
		os.MkdirAll(outDir, 0755)
	} else if err != nil {
		return err
	}

	outPath := filepath.Join(outDir, fmt.Sprintf("%s_chunk_%d.json", outFile, index))
	data, err := json.MarshalIndent(chunk, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(outPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s storage) Csv() ([]string, error) {
	if s.err != nil {
		return nil, s.err
	}

	if s.path == "" {
		return nil, fmt.Errorf("must specify using PublicPath() or StoragePath() first")
	}

	file, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("the specify file is not found")
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []string
	scanFile := bufio.NewScanner(file)
	for scanFile.Scan() {
		data = append(data, scanFile.Text())
	}

	if err := scanFile.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (s storage) Url() (*string, error) {
	if s.err != nil {
		return nil, s.err
	}

	if s.path == "" {
		return nil, fmt.Errorf("must specify using PublicPath() or StoragePath() first")
	}

	return &s.path, nil
}

func (s storage) Error() error {
	return s.err
}
