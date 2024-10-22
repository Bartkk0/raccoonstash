package internal

import (
	"context"
	"errors"
	"os"
	"raccoonstash/internal/repository"
	"time"
)

type File struct {
	data       repository.File
	Filename   string
	UploadTime time.Time

	file *os.File
}

func FindFileByHash(context context.Context, hash string) (File, error) {
	var file File

	fileData, err := Queries.GetFileByHash(context, hash)
	if err != nil {
		return File{}, errors.New("file not found")
	}

	file.data = fileData
	file.Filename = file.data.Filename
	file.UploadTime = file.data.UploadedAt

	return file, nil
}

func (f *File) open() error {
	var err error
	f.file, err = FindFile(f.data)
	return err
}

func (f *File) CheckPassword(password string) bool {
	if !f.data.Password.Valid {
		return true
	}

	return f.data.Password.String == password
}

func (f *File) Read(p []byte) (n int, err error) {
	if f.file == nil {
		err := f.open()
		if err != nil {
			return 0, err
		}
	}
	return f.file.Read(p)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.file == nil {
		err := f.open()
		if err != nil {
			return 0, err
		}
	}
	return f.file.Seek(offset, whence)
}
