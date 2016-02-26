package ignore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var IgnoreFile string
var PfsDir string

func IgnoreExists() bool {
	_, err := os.Stat(IgnoreFile)
	return !os.IsNotExist(err)
}

func contains(globs []string, file string) bool {
	for _, globElement := range globs {
		if globElement == file {
			return true
		} else if strings.Contains(file, globElement) {
			return true
		}
	}
	return false
}

func glob(updateFile string) bool {
	shouldIgnore := false
	file, err := ioutil.ReadFile(IgnoreFile)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	values := strings.Split(string(file), "\n")
	for _, element := range values {
		if element != "" {
			glob, _ := filepath.Glob(PfsDir + "/" + element)
			if glob != nil {
				fmt.Println(glob, updateFile)
				shouldIgnore = shouldIgnore || contains(glob, updateFile)
			}
		}
	}
	return shouldIgnore
}

func PfsIgnore(file string) bool {
	ignore := false
	file = path.Join(PfsDir, file)
	if IgnoreExists() {
		ignore = glob(file)
	}
	return ignore
}
