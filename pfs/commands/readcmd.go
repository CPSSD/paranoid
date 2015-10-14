package commands

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

func ReadCommand(args []string) {
	if len(args) < 2 {
		log.Fatal("Not enough arguments!")
	}
	directory := args[0]
	fileNameBytes, err := ioutil.ReadFile(path.Join(directory, "names", args[1]))
	fileName := string(fileNameBytes)
	file, err := os.Open(path.Join(directory, "contents", fileName))
	checkErr("read", err)
	if len(args) == 2 {
		bytesRead := make([]byte, 1024)
		for {
			n, err := file.Read(bytesRead)
			checkErr("read", err)
			io.WriteString(os.Stdout, string(bytesRead))
			if n < 1024 {
				break
			}
		}
	} else if len(args) > 2 {
		bytesRead := make([]byte, 1024)
		maxRead := 100000000
		if len(args) > 3 {
			maxRead, err = strconv.Atoi(args[3])
			checkErr("read", err)
		}
		off, err := strconv.Atoi(args[2])
		checkErr("read", err)
		offset := int64(off)
		for {
			n, err := file.ReadAt(bytesRead, offset)
			if n > maxRead {
				bytesRead = bytesRead[0:maxRead]
				io.WriteString(os.Stdout, string(bytesRead))
				break
			}
			maxRead = maxRead - n
			if err == io.EOF {
				io.WriteString(os.Stdout, string(bytesRead))
				break
			}
			checkErr("read", err)
			io.WriteString(os.Stdout, string(bytesRead))
		}
	}
}
