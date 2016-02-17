package commands

import (
	"github.com/cpssd/paranoid/logger"
	"math/rand"
	"os"
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
