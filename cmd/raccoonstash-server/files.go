package main

import "C"
import (
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"raccoonstash/internal"
	"raccoonstash/internal/repository"
	"strings"
)

func getFileHandler(writer http.ResponseWriter, request *http.Request) {
	hash := request.PathValue("hash")
	hash = strings.Split(hash, ".")[0]
	password := request.FormValue("password")

	file, err := internal.FindFileByHash(request.Context(), hash)
	if err != nil {
		http.NotFound(writer, request)
		return
	}

	if !file.CheckPassword(password) {
		writer.WriteHeader(http.StatusUnauthorized)
		err := renderTemplate(writer, "password.html", nil)
		if err != nil {
			log.Println(err)
		}
		return
	}

	writer.Header().Add("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", url.QueryEscape(file.Filename)))
	http.ServeContent(writer, request, file.Filename, file.UploadTime, &file)
}

func postFileHandler(writer http.ResponseWriter, request *http.Request) {
	if !internal.VerifyToken(request.Context(), request.FormValue("token")) {
		writer.WriteHeader(http.StatusUnauthorized)
		err := renderTemplate(writer, "token.html", nil)
		if err != nil {
			log.Println(err)
		}
		return
	}

	file, fileHeader, err := request.FormFile("file")
	if err != nil {
		log.Println("Failed to read file from form!")
		http.Error(writer, "Invalid form!", http.StatusBadRequest)
		return
	}

	password := formGetNullString(request.Form, "password")

	expiresTime, err := formGetNullTimeDuration(request.Form, "expires")
	if err != nil {
		http.Error(writer, "Unable to parse expires field", http.StatusBadRequest)
		return
	}

	hash := crypto.SHA256.New()
	io.Copy(hash, file)
	file.Seek(0, io.SeekStart)
	hashStr := hex.EncodeToString(hash.Sum(nil))

	fileExtension := path.Ext(fileHeader.Filename)
	filename := fmt.Sprint(hashStr, fileExtension)

	targetFile, _ := internal.CreateFile(hashStr, fileExtension)
	defer targetFile.Close()

	tx, _ := internal.DB.BeginTx(request.Context(), nil)
	queries := internal.Queries.WithTx(tx)

	_, err = queries.InsertFile(request.Context(), repository.InsertFileParams{
		Hash:      hashStr,
		Filename:  fileHeader.Filename,
		Extension: fileExtension,
		Size:      fileHeader.Size,
		Password:  password,
		ExpiresAt: expiresTime,
	})

	_, err = io.Copy(targetFile, file)
	if err != nil {
		log.Println("Failed to copy uploaded file!")
		return
	}
	log.Println("Uploaded file", fileHeader.Filename, "saved as", filename)
	http.Redirect(writer, request, getFileUrl(hashStr, fileExtension), http.StatusTemporaryRedirect)

	err = tx.Commit()
	if err != nil {
		log.Println("Error while commiting uploaded file!")
		log.Println(err)
		return
	}
}

func getFileUrl(hash string, extension string) string {
	return "/file/" + hash + extension
}
