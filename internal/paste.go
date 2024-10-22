package internal

import (
	"context"
	"database/sql"
	"os"
	"raccoonstash/internal/repository"
	"time"
)

type Paste struct {
	data       repository.Paste
	Filename   string
	Language   sql.NullString
	UploadTime time.Time
	Size       int64

	file *os.File
}

func FindPasteByHash(context context.Context, hash string) (Paste, error) {
	var paste Paste

	pasteData, err := Queries.GetPasteByHash(context, hash)
	if err != nil {
		return Paste{}, nil
	}

	paste.data = pasteData
	paste.Filename = paste.data.Filename
	paste.UploadTime = paste.data.UploadedAt
	paste.Language = paste.data.Language
	paste.Size = paste.data.Size

	return paste, nil
}

func (p *Paste) open() error {
	var err error
	p.file, err = FindPaste(p.data)
	return err
}

func (p *Paste) CheckPassword(password string) bool {
	if !p.data.Password.Valid {
		return true
	}

	return p.data.Password.String == password
}

func (p *Paste) Read(out []byte) (n int, err error) {
	if p.file == nil {
		err := p.open()
		if err != nil {
			return 0, err
		}
	}
	return p.file.Read(out)
}

func (p *Paste) Seek(offset int64, whence int) (int64, error) {
	if p.file == nil {
		err := p.open()
		if err != nil {
			return 0, err
		}
	}
	return p.file.Seek(offset, whence)
}
