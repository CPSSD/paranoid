package commands

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"syscall"
	"time"
)

type statInfo struct {
	Length int64     `json:"length"`
	Ctime  time.Time `json:"ctime"`
	Mtime  time.Time `json:"mtime"`
	Atime  time.Time `json:"atime"`
}

func StatCommand(args []string) {
	if len(args) < 2 {
		log.Fatal("Not enough arguments!")
	}
	directory := args[0]
	fileNameBytes, err := ioutil.ReadFile(path.Join(directory, "names", args[1]))
	fileName := string(fileNameBytes)
	file, err := os.Open(path.Join(directory, "contents", fileName))
	checkErr("stat", err)
	fi, err := file.Stat()
	stat := fi.Sys().(*syscall.Stat_t)
	atime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	statData := &statInfo{
		Length: fi.Size(),
		Mtime:  fi.ModTime(),
		Ctime:  ctime,
		Atime:  atime}
	jsonData, err := json.Marshal(statData)
	checkErr("stat", err)
	io.WriteString(os.Stdout, string(jsonData))
}
