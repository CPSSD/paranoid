package commands

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	pb "github.com/cpssd/paranoid/proto/raft"
	"github.com/cpssd/paranoid/raft"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strconv"
)

// History is the top level function when paranoid-cli is called
func History(c *cli.Context) {
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
	var target string
	if fileSystemExists(givenString) {
		target = path.Join(usr.HomeDir, ".pfs", givenString, "meta", "activity_logs")
	} else {
		target = givenString
	}
	read(target, c)
}

// fileSystemExists checks if there is a folder in ~/.pfs with the given name
func fileSystemExists(fsname string) bool {
	usr, err := user.Current()
	if err != nil {
		Log.Fatal(err)
	}

	dirpath := path.Join(usr.HomeDir, ".pfs", fsname)
	_, err = ioutil.ReadDir(dirpath)
	return err == nil
}

// read shows the history of a log in the given directory
func read(directory string, c *cli.Context) {
	tempDir := os.TempDir()
	filePath := path.Join(tempDir, "log.pfslog")
	logsToLogfile(directory, filePath, c)
	defer os.Remove(filePath)
	less := exec.Command("less", filePath)
	less.Stdout = os.Stdout
	less.Stdin = os.Stdin
	less.Stderr = os.Stderr
	less.Run()
}

// logsToLogFile converts the binary logs in the logDir paramenter
// to a human readable file in the given filePath paramenter.
func logsToLogfile(logDir, filePath string, c *cli.Context) {
	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		Log.Verbose("read dir:", logDir, "err:", err)
		cli.ShowCommandHelp(c, "history")
		os.Exit(1)
	}
	writeFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	numRunes := len(strconv.Itoa(len(files)))
	numRunesString := strconv.Itoa(numRunes)
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		p, err := fileToProto(file, logDir)
		if err != nil {
			log.Fatalln(err)
		}

		writeFile.WriteString(toLine(i+1, numRunesString, p))
	}
	writeFile.Close()
}

// fileToProto converts a given file with a protobuf to a protobuf object
func fileToProto(file os.FileInfo, directory string) (entry *pb.LogEntry, err error) {
	filePath := path.Join(directory, file.Name())
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("Failed to read logfile: " + file.Name())
	}
	entry = &pb.LogEntry{}
	err = proto.Unmarshal(fileData, entry)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal file data")
	}
	return entry, nil
}

// toLine converts a protobuf object to a human readable string representation
// the logNum parameter is the number of the logfile and the pad parameter is
// the number of runes to be used for writing the logNum and Term in string form.
func toLine(logNum int, pad string, p *pb.LogEntry) string {
	marker := fmt.Sprintf("file: %-"+pad+"d, term: %-"+pad+"d ", logNum, p.Term)

	if p.Entry.Type == 0 { // command
		return fmt.Sprintln(marker, "Command: ", commandString(p.Entry.Command))
	} else if p.Entry.Type == 1 { // config
		return fmt.Sprintln(marker, "ConfigChange: ", p.Entry.Config)
	} else { // demo
		return fmt.Sprintln(marker, "Demo: ", p.Entry.Demo)
	}
}

// commandString returns the string representation of a StateMachineCommand
func commandString(cmd *pb.StateMachineCommand) string {
	typeStr := typeString(cmd.Type)
	cmd.Type = 0
	size := len(cmd.Data)
	if size > 0 {
		cmd.Data = nil
		return fmt.Sprint(fmt.Sprintf("%-10s", typeStr), cmd, "Data: ", bytesString(size), "\n")
	}
	return fmt.Sprint(fmt.Sprintf("%-10s", typeStr), cmd, "\n")
}

// bytesString returns the human readable representation of a data size=
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

// typeString returns the string representation of a log type.
func typeString(ty uint32) string {
	switch ty {
	case raft.TYPE_WRITE:
		return "Write"
	case raft.TYPE_CREAT:
		return "Creat"
	case raft.TYPE_CHMOD:
		return "Chmod"
	case raft.TYPE_TRUNCATE:
		return "Truncate"
	case raft.TYPE_UTIMES:
		return "Utimes"
	case raft.TYPE_RENAME:
		return "Rename"
	case raft.TYPE_LINK:
		return "Link"
	case raft.TYPE_SYMLINK:
		return "Symlink"
	case raft.TYPE_UNLINK:
		return "Unlink"
	case raft.TYPE_MKDIR:
		return "Mkdir"
	case raft.TYPE_RMDIR:
		return "Rmdir"
	default:
		return "Unknown"
	}
}
