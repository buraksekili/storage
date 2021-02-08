package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Storage interface {
	Save(data []byte, ext, imgName string) error
	// Save(data bytes.Buffer, ext, imgName string) (int32, error)
}

type LocalImgStorage struct {
	path string
	l    *log.Logger
}

func NewLocalImgStorage(p string, l *log.Logger) *LocalImgStorage {
	return &LocalImgStorage{p, l}
}

func (lis *LocalImgStorage) Save(chunk []byte, ext, imgName string) error {
	fullImgName := fmt.Sprintf("%s.%s", imgName, ext)
	fp := filepath.Join(lis.path, fullImgName)

	f, err := os.OpenFile(fp, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("[ERROR] cannot open file to path %s: %v", fp, err)
	}
	defer f.Close()

	if _, err = f.Write(chunk); err != nil {
		return fmt.Errorf("[ERROR] cannot open file to path %s: %v", fp, err)
	}

	return nil
}

func validate(imgName, ext string) bool {
	return imgName != "" && ext != ""
}
