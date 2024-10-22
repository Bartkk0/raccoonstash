package internal

import (
	"log"
	"os"
	"path"
	"raccoonstash/internal/repository"
)

var FilesDir string
var PastesDir string

func InitializeStorage() {
	FilesDir = path.Join(Config.DataDir, "files/")
	PastesDir = path.Join(Config.DataDir, "pastes/")

	// Required directories
	checkAndCreateDirectory(Config.DataDir)
	checkAndCreateDirectory(FilesDir)
	checkAndCreateDirectory(PastesDir)
}

func FindFile(file repository.File) (*os.File, error) {
	return os.Open(path.Join(FilesDir, file.Hash+file.Extension))
}

func CreateFile(hash string, extension string) (*os.File, error) {
	return os.OpenFile(path.Join(FilesDir, hash+extension), os.O_CREATE|os.O_RDWR, 0666)
}

func FindPaste(paste repository.Paste) (*os.File, error) {
	return os.Open(path.Join(PastesDir, paste.Hash+".txt"))
}

func CreatePaste(hash string) (*os.File, error) {
	targetFilename := hash + ".txt"
	return os.OpenFile(path.Join(Config.DataDir, "pastes/", targetFilename), os.O_CREATE|os.O_WRONLY, 0666)
}

func checkAndCreateDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println("Data directory", path, "does not exist, creating...")
		err = os.Mkdir(path, 0770)
		if err != nil {
			log.Fatalln("Could not create directory", path, "!\n", err.Error())
		}
	}
}
