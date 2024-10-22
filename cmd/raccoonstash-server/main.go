package main

import (
	"embed"
	"flag"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"log"
	"net/http"
	"os"
	"raccoonstash/internal"
	"raccoonstash/internal/repository"
)

var serverConfig struct {
	listenAddress string
	dev           bool
}

//go:embed static/*
var staticEmbedFs embed.FS

var staticFs http.FileSystem

//go:embed templates/*
var templatesEmbedFs embed.FS

var templatesFs fs.FS

func main() {
	internal.AddGlobalFlags()
	flag.BoolVar(&serverConfig.dev, "dev", false, "Enable dev mode (use files from source tree instead of EmbedFS)")
	flag.StringVar(&serverConfig.listenAddress, "listen", ":8080", "Address to listen on")
	flag.Parse()

	internal.InitializeDatabase()
	internal.InitializeStorage()

	setupFilesystems()
	loadTemplates()
	if serverConfig.dev {
		setupRefreshing()
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		stats, _ := internal.Queries.GetStats(request.Context())

		err := renderTemplate(writer, "index.html", struct {
			Languages []string
			Stats     repository.GetStatsRow
		}{
			Languages: lexers.Names(false),
			Stats:     stats,
		})
		if err != nil {
			log.Println(err)
			return
		}
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticFs)))
	http.HandleFunc("/file/upload", postFileHandler)
	http.HandleFunc("/file/{hash}", getFileHandler)
	http.HandleFunc("/paste/upload", postPasteHandler)
	http.HandleFunc("/paste/{hash}", getPasteHandler)
	http.HandleFunc("/paste/{hash}/raw", getPasteRawHandler)

	log.Println("Listening on", serverConfig.listenAddress)
	log.Fatal(http.ListenAndServe(serverConfig.listenAddress, nil))
}

func setupRefreshing() {
	log.Println("Staring template watcher...")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panicln(err)
	}

	//defer watcher.Close()
	for template := range templates {
		watcher.Add("cmd/raccoonstash-server/templates/" + template)
	}
	watcher.Add("cmd/raccoonstash-server/templates/_template.html")

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					println("Reloading templates...")
					loadTemplates()
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				println(err)
			}
		}
	}()
}

func setupFilesystems() {
	if serverConfig.dev {
		staticFs = http.Dir("cmd/server/static/")
		templatesFs = os.DirFS("cmd/server/templates/")
	} else {
		sub, _ := fs.Sub(staticEmbedFs, "static")
		staticFs = http.FS(sub)
		sub, _ = fs.Sub(templatesEmbedFs, "templates")
		templatesFs = sub
	}
}
