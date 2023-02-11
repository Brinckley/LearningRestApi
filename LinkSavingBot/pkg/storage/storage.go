package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"time"
)

type Storage interface {
	Save(page *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(page *Page) error
	IsExists(p *Page) (bool, error)
}

var ErrNoSavedFiles = errors.New("no saved files")

type Page struct {
	URL      string
	UserName string
	Created  time.Time
}

func (p Page) Hash() (string, error) {
	sh := sha1.New()

	if _, err := io.WriteString(sh, p.URL); err != nil {
		return "", fmt.Errorf("can't hash the link : %s", err.Error())
	}

	if _, err := io.WriteString(sh, p.UserName); err != nil {
		return "", fmt.Errorf("can't hash the username : %s", err.Error())
	}

	return fmt.Sprintf("%x", sh.Sum(nil)), nil
}

func fileName(p *Page) (string, error) {
	return p.Hash()
}
