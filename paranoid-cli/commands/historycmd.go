package commands

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/cpssd/paranoid/pfsd/activitylogger"
	pb "github.com/cpssd/paranoid/proto/activitylogger"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strconv"
)

var cl *cli.Context

func History(c *cli.Context) {
	cl = c
	args := c.Args()
	if len(args) < 1 {
		cli.ShowCommandHelp(c, "history")
		os.Exit(1)
	}

	usr, err := user.Current()
	if err != nil {
		Log.Fatal(err)
	}
	givenString := args[0]
	target := ""
	if fileSystemExists(givenString) {
		target = path.Join(usr.HomeDir, ".pfs", givenString, "meta", "activity_logs")
	} else {
		target = givenString
	}

	Read(target)
}

func fileSystemExists(fsname string) bool {
	usr, err := user.Current()
	if err != nil {
		Log.Fatal(err)
	}

	dirpath := path.Join(usr.HomeDir, ".pfs", fsname)
	_, err = ioutil.ReadDir(dirpath)
	return err == nil
}

func Read(directory string) {
	tempDir := os.TempDir()
	filePath := path.Join(tempDir, "log.pfslog")
	LogsToLogfile(directory, filePath)
	defer os.Remove(filePath)

	less := exec.Command("less", filePath)
	less.Stdout = os.Stdout
	less.Stdin = os.Stdin
	less.Stderr = os.Stderr
	less.Run()
}

func LogsToLogfile(logDir, filePath string) {
	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		Log.Verbose("read dir:", logDir, "err:", err)
		cli.ShowCommandHelp(cl, "history")
		os.Exit(1)
	}

	writeFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	t := len(strconv.Itoa(len(files)))
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		p, err := fileToProto(file, logDir)
		if err != nil {
			log.Fatalln(err)
		}

		writeFile.WriteString(toLine(i+1, t, p))
	}
	writeFile.Close()
}

func fileToProto(file os.FileInfo, directory string) (entry *pb.Entry, err error) {
	filePath := path.Join(directory, file.Name())
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("Failed to read logfile: " + file.Name())
	}

	entry = &pb.Entry{}
	err = proto.Unmarshal(fileData, entry)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal file data")
	}

	return entry, nil
}

func toLine(i, t int, p *pb.Entry) string {
	iStr := strconv.Itoa(i)
	pad := padding(t + 1 - len(iStr))

	typeStr := activitylogger.TypeString(p.Type)
	p.Type = 0
	pad2 := padding(10 - len(typeStr))

	size := len(p.Data)
	if size > 0 {
		p.Data = nil
		return fmt.Sprint(iStr, pad, ": ", typeStr, pad2, p, "Data: ", bytesString(size), "\n")
	}

	return fmt.Sprint(iStr, pad, ": ", typeStr, pad2, p, "\n")
}

func padding(i int) string {
	str := ""
	for j := 0; j < i; j++ {
		str += " "
	}
	return str
}

func bytesString(bytes int) string {
	if bytes < 1000 {
		return fmt.Sprint(bytes, "B")
	} else if bytes < 1000000 {
		return fmt.Sprint(bytes/1000, "KB")
	} else if bytes < 1000000000 {
		return fmt.Sprint(bytes/1000000, "MB")
	} else if bytes < 1000000000000 {
		return fmt.Sprint(bytes/1000000000, "GB")
	} else {
		return fmt.Sprint(bytes/1000000000000, "TB")
	}
}
