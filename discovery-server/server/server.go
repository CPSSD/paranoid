package server

import (
	"fmt"
	"github.com/cpssd/paranoid/logger"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type FileCache struct {
	Uuid           string
	AccessAmmount  int32
	AccessLimit    int32
	FileData       []byte
	FilePath       string
	isServing      bool
	ExpirationTime time.Time
}

type FileserverServer struct{}

var FileMap map[string]*FileCache
var Log *logger.ParanoidLogger
var Port string

func getFileFromHash(hash string) ([]byte, string, error) {
	value, ok := FileMap[hash]
	if !ok {
		return []byte(""), "", fmt.Errorf("No Valid File Found")
	}
	if time.Now().After(value.ExpirationTime) || value.AccessAmmount >= value.AccessLimit {
		Log.Info("Expired Filed attempted to be accessed")
		if !FileMap[hash].isServing {
			delete(FileMap, hash)
		}
		FileMap[hash].isServing = true
		return []byte(""), "", fmt.Errorf("File Expired")
	}
	value.AccessAmmount++
	return value.FileData, filepath.Base(value.FilePath), nil
}

func ServeFiles(serverPort string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			log.Println(r.URL)
			w.Write([]byte("Welcome to the Paranoid File Server. Please enter your file Hash in the URL"))
		} else {
			file, name, err := getFileFromHash(r.URL.Path[1:])
			if err != nil {
				w.Write([]byte("File Not Found"))
			} else {
				w.Header().Set("Content-Disposition", "attachment; filename="+name)
				w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
				Log.Info("sending File", name, "to user")
				w.Write(file)
				FileMap[r.URL.Path[1:]].isServing = false
			}
		}
	})
	Port = ":" + serverPort
	http.ListenAndServe(Port, nil)
}
