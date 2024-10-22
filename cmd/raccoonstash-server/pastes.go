package main

import (
	"bufio"
	"bytes"
	"crypto"
	"encoding/hex"
	"fmt"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"raccoonstash/internal"
	"raccoonstash/internal/repository"
)

func getPasteHandler(writer http.ResponseWriter, request *http.Request) {
	hash := request.PathValue("hash")
	password := request.FormValue("password")

	paste, err := internal.FindPasteByHash(request.Context(), hash)
	if err != nil {
		http.NotFound(writer, request)
		return
	}
	if !paste.CheckPassword(password) {
		writer.WriteHeader(http.StatusUnauthorized)
		err := renderTemplate(writer, "password.html", nil)
		if err != nil {
			log.Println(err)
		}
		return
	}

	data, _ := io.ReadAll(&paste)
	dataStr := string(data)

	var lexer chroma.Lexer
	if !paste.Language.Valid {
		lexer = lexers.Fallback
	} else if paste.Language.String == "auto" {
		lexer = lexers.Analyse(dataStr)
	} else {
		lexer = lexers.Get(paste.Language.String)
	}

	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get("dracula")

	formatter := html.New(html.TabWidth(4), html.WithLineNumbers(true), html.WithLinkableLineNumbers(true, "line-"), html.WrapLongLines(true), html.WithClasses(true))
	tokens, _ := lexer.Tokenise(nil, dataStr)

	var buffer bytes.Buffer
	var bufferWriter = bufio.NewWriter(&buffer)
	err = formatter.Format(bufferWriter, style, tokens)
	if err != nil {
		log.Println(err)
	}
	bufferWriter.Flush()

	var cssBuffer bytes.Buffer
	var cssBufferWriter = bufio.NewWriter(&cssBuffer)
	formatter.WriteCSS(cssBufferWriter, style)
	cssBufferWriter.Flush()

	err = renderTemplate(writer, "paste.html", struct {
		FormattedText template.HTML
		Paste         internal.Paste
		Css           template.CSS
		RawUrl        template.URL
	}{
		FormattedText: template.HTML(buffer.String()),
		Paste:         paste,
		Css:           template.CSS(cssBuffer.String()),
		RawUrl:        template.URL(getRawPasteUrl(hash)),
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func getPasteRawHandler(writer http.ResponseWriter, request *http.Request) {
	hash := request.PathValue("hash")
	password := request.FormValue("password")

	paste, err := internal.FindPasteByHash(request.Context(), hash)
	if err != nil {
		http.NotFound(writer, request)
		return
	}
	if !paste.CheckPassword(password) {
		writer.WriteHeader(http.StatusUnauthorized)
		err := renderTemplate(writer, "password.html", nil)
		if err != nil {
			log.Println(err)
		}
		return
	}

	writer.Header().Add("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", url.QueryEscape(paste.Filename)))
	http.ServeContent(writer, request, paste.Filename, paste.UploadTime, &paste)
}

func postPasteHandler(writer http.ResponseWriter, request *http.Request) {
	if !internal.VerifyToken(request.Context(), request.FormValue("token")) {
		writer.WriteHeader(http.StatusUnauthorized)
		err := renderTemplate(writer, "token.html", nil)
		if err != nil {
			log.Println(err)
		}
		return
	}

	text := request.PostFormValue("text")
	filename := request.FormValue("filename")

	language := formGetNullString(request.Form, "language")
	password := formGetNullString(request.Form, "password")

	if language.Valid && language.String == "" {
		language.Valid = false
	}

	expiresTime, err := formGetNullTimeDuration(request.Form, "expires")
	if err != nil {
		http.Error(writer, "Unable to parse expires field", http.StatusBadRequest)
		return
	}

	if filename == "" {
		filename = "paste.txt"
	}

	hash := crypto.MD5.New()
	hash.Write([]byte(text))
	hashStr := hex.EncodeToString(hash.Sum(nil))

	targetFile, _ := internal.CreatePaste(hashStr)
	defer targetFile.Close()

	tx, _ := internal.DB.BeginTx(request.Context(), nil)
	repo := internal.Queries.WithTx(tx)

	err = repo.InsertPaste(request.Context(), repository.InsertPasteParams{
		Hash:      hashStr,
		Filename:  filename,
		Size:      int64(len(text)),
		Language:  language,
		Password:  password,
		ExpiresAt: expiresTime,
	})

	if err != nil {
		log.Panicln(err)
	}

	_, err = targetFile.WriteString(text)
	if err != nil {
		log.Println("Failed to copy uploaded paste!")
		return
	}
	log.Println("Uploaded paste", filename, "saved as", targetFile.Name())
	http.Redirect(writer, request, getPasteUrl(hashStr), http.StatusTemporaryRedirect)

	err = tx.Commit()
	if err != nil {
		log.Println("Error while commiting paste transaction")
		log.Println(err)
	}
}

func getPasteUrl(hash string) string {
	return "/paste/" + hash
}

func getRawPasteUrl(hash string) string {
	return "/paste/" + hash + "/raw"
}
