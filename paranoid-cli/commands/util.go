package commands

import (
	"errors"
	"fmt"
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

type fileSystemAttributes struct {
	Encrypted    bool `json:"encrypted"`
	KeyGenerated bool `json:"keygenerated"`
	NetworkOff   bool `json:"networkoff"`
}

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

	dirpath := path.Join(usr.HomeDir, ".pfs", "filesystems", fsname)
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

func getFsMeta(usr *user.User, pfsName string) (string, string, string, string) {
	pfsDir := path.Join(usr.HomeDir, ".pfs", "filesystems", pfsName)
	if _, err := os.Stat(pfsDir); err != nil {
		fmt.Printf("%s does not exist. Please call 'paranoid-cli init' before running this command.", pfsDir)
		Log.Fatal("PFS directory does not exist.")
	}

	uuid, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "uuid"))

	if err != nil {
		fmt.Println("Error Reading supplied file:")
		Log.Fatal("Cant Reading uuid file")
	}

	ip, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "ip"))
	port, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "port"))
	pool, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "pool"))
	if err != nil {
		fmt.Println("Could not find Ip address of the file server")
		Log.Fatal("Unable to read Ip and Port of discovery server", err)
	}
	return string(ip), string(port), string(uuid), string(pool)
}
