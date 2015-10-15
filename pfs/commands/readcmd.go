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
	verboseLog("read command called")
	if len(args) < 2 {
		log.Fatal("Not enough arguments!")
	}
	directory := args[0]
	verboseLog("read : given directory = " + directory)
	fileNameBytes, err := ioutil.ReadFile(path.Join(directory, "names", args[1]))
	fileName := string(fileNameBytes)
	file, err := os.Open(path.Join(directory, "contents", fileName))
	checkErr("read", err)
	if len(args) == 2 {
		verboseLog("read : reading whole file")
		bytesRead := make([]byte, 1024)
		for {
			n, err := file.Read(bytesRead)
			if err != io.EOF {
				checkErr("read", err)
			}
			io.WriteString(os.Stdout, string(bytesRead))
			if n < 1024 {
				break
			}
		}
	} else if len(args) > 2 {
		bytesRead := make([]byte, 1024)
		maxRead := 100000000
		if len(args) > 3 {
			verboseLog("read : " + args[3] + " bytes starting at " + args[2])
			maxRead, err = strconv.Atoi(args[3])
			checkErr("read", err)
		} else {
			verboseLog("read : from " + args[2] + " to end of file")
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
