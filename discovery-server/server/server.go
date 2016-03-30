package server

import (
	"fmt"
	"github.com/cpssd/paranoid/logger"
	"log"
	"net/http"
	"time"
)

type FileCache struct {
	AccessAmmount  int32
	AccessLimit    int32
	FileData       []byte
	FileName       string
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
		return []byte(""), "", fmt.Errorf("File Expired")
	}
	return value.FileData, value.FileName, nil
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
			}
		}
	})
	Port = ":" + serverPort
	http.ListenAndServe(Port, nil)
}
