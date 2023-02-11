package storage

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const defaultPerm = 0774

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) FileStorage {
	return FileStorage{basePath: basePath}
}

func (s FileStorage) Save(p *Page) error { // saving info
	fPath := filepath.Join(s.basePath, p.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return fmt.Errorf("can't create directory : %s", err.Error())
	}

	fName, err := fileName(p)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(p); err != nil {
		return err
	}

	return nil
}

func (s FileStorage) PickRandom(username string) (*Page, error) {
	path := filepath.Join(s.basePath, username)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("can't read directory with path %s, error : %s", path, err)
	}

	if len(files) == 0 {
		return nil, ErrNoSavedFiles
	}

	rand.Seed(time.Now().Unix())
	n := rand.Intn(len(files)) // random id

	f := files[n]

	return s.decodePage(filepath.Join(path, f.Name()))
}

func (s FileStorage) Remove(p *Page) error {
	fName, err := fileName(p)
	if err != nil {
		return err
	}

	path := filepath.Join(s.basePath, p.UserName, fName)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("can't remove the file with path : %s, error : %s", path, err.Error())
	}

	return nil
}

func (s FileStorage) IsExists(p *Page) (bool, error) {
	fName, err := fileName(p)
	if err != nil {
		return false, err
	}

	path := filepath.Join(s.basePath, p.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("can't check file existence : %s, error : %s", path, err.Error())
	}

	return true, nil
}

func (s FileStorage) decodePage(fPath string) (*Page, error) { // custom decoder for web page
	f, err := os.Open(fPath)
	if err != nil {
		return nil, fmt.Errorf("can't open file %s, error : %s", fPath, err.Error())
	}
	defer func() { _ = f.Close() }()

	var p Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("can't decode tha page : %s, error : %s", p.URL, err)
	}

	return &p, nil
}
