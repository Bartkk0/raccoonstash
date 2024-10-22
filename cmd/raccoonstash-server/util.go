package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"maps"
	"math"
	"net/url"
	"time"
)

var templates = map[string]*template.Template{
	"index.html":    nil,
	"paste.html":    nil,
	"token.html":    nil,
	"password.html": nil,
}

func formatFileSize(bytes int64) string {
	if bytes < 1000 {
		return fmt.Sprintf("%d B", bytes)
	}
	unit := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	i := int(math.Floor(math.Log(float64(bytes)) / math.Log(1000)))
	val := float64(bytes) / math.Pow(1000, float64(i))
	return fmt.Sprintf("%.1f %s", val, unit[i])
}

var templateFunctions = template.FuncMap{
	"formatFileSize": formatFileSize,
}

func loadTemplates() {
	for key := range maps.Keys(templates) {
		templates[key] = template.Must(template.New(key).Funcs(templateFunctions).ParseFS(templatesFs, "_template.html", key))
	}
}

func renderTemplate(writer io.Writer, name string, data any) error {
	if templates[name] == nil {
		log.Panicln("Template not found", name)
	}
	return templates[name].Execute(writer, data)
}

func formGetNullTimeDuration(form url.Values, key string) (sql.NullTime, error) {
	if !form.Has(key) {
		return sql.NullTime{}, nil
	}

	if form.Get(key) == "" {
		return sql.NullTime{}, nil
	}

	duration, err := time.ParseDuration(form.Get(key))
	if err != nil {
		return sql.NullTime{}, err
	}

	return sql.NullTime{
		Valid: true,
		Time:  time.Now().Add(duration),
	}, nil
}

func formGetNullString(form url.Values, key string) sql.NullString {
	if !form.Has(key) {
		return sql.NullString{}
	}

	return sql.NullString{
		Valid:  true,
		String: form.Get(key),
	}
}
