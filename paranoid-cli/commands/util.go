package commands

import (
	"errors"
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/raft"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"math/rand"
	"os"
	"os/user"
	"path"
	"time"
)

var Log *logger.ParanoidLogger

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func getRandomName() string {
	prefix := []string{
		"raging",
		"violent",
		"calm",
		"peaceful",
		"strange",
		"hungry",
	}
	postfix := []string{
		"dolphin",
		"snake",
		"elephant",
		"fox",
		"dog",
		"cat",
		"rabbit",
	}

	rand.Seed(time.Now().Unix())
	return prefix[rand.Int()%len(prefix)] + "_" + postfix[rand.Int()%len(postfix)]
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
