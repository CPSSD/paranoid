package ignore

import (
	"github.com/cpssd/paranoid/logger"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"
)

type SolvedPaths struct {
	value  bool
	solved time.Time
}

var IgnoreFile string
var Log *logger.ParanoidLogger
var IgnoreData []string
var FileLastUpdated time.Time
var FoundGlobs map[string]SolvedPaths

const recursiveStar string = "**"
const star string = "*"
const negation string = "!"

func IgnoreExists() bool {
	_, err := os.Stat(IgnoreFile)
	return !os.IsNotExist(err)
}

func dualStarGlob(pattern, file string) bool {
	if pattern == recursiveStar {
		return true
	}

	patternSplit := strings.Split(pattern, recursiveStar)

	// No ** split occured
	if len(patternSplit) == 1 {
		return pattern == file
	}
	//checking files start the same
	if !strings.HasPrefix(file, patternSplit[0]) {
		return false
	}

	//search the middle Patterns
	for i := 1; i < len(patternSplit)-1; i++ {
		if !strings.Contains(file, patternSplit[i]) {
			return false
		}
		index := strings.Index(file, patternSplit[i]) + len(patternSplit[i])
		file = file[index:]
	}

	return strings.HasSuffix(pattern, recursiveStar) || strings.HasSuffix(file, patternSplit[len(patternSplit)-1])

}

func starGlob(pattern, file string) bool {

	patternSplit := strings.Split(pattern, star)
	//no star found
	if len(patternSplit) == 1 {
		return pattern == file
	}
	//checking if pattern contains * and ensuring that the file contains the pattern
	shouldGlob := strings.HasPrefix(file, patternSplit[0])
	if !shouldGlob {
		return false
	}
	//search the middle Patterns
	for i := 1; i < len(patternSplit)-1; i++ {
		shouldGlob = shouldGlob || !strings.Contains(file, patternSplit[i])
		index := strings.Index(file, patternSplit[i]) + len(patternSplit[i])
		file = file[index:]
	}
	if patternSplit[len(patternSplit)-1] == "" {
		return true
	} else {
		return shouldGlob && strings.HasSuffix(file, patternSplit[len(patternSplit)-1])
	}
}

func checkDir(pattern, file string) bool {
	patternArr := strings.Split(pattern, "/")
	fileArr := strings.Split(file, "/")
	dirsMatch := false

	if len(patternArr) == 1 {
		return patternArr[0] == fileArr[0]
	}

	for i := 0; i < len(patternArr); i++ {
		if patternArr[i] == "" || fileArr[i] == "" {
			break
		}
		dirsMatch = dirsMatch && (patternArr[i] == fileArr[i])
	}

	return dirsMatch
}

func Glob(pattern, file string) bool {
	//Checking if the start of the File a Negation is set
	negationSet := false
	if string(pattern[0]) == "!" {
		negationSet = true
		//strip the not as its not of use anymore
		pattern = strings.Trim(pattern, "!")
	}
	//Removing trailing slashes
	pattern = strings.TrimSuffix(pattern, "/")

	if pattern == star || pattern == recursiveStar {
		return true
	}

	if pattern == "" {
		return false
	}

	isGlobbed := starGlob(pattern, file) || dualStarGlob(pattern, file) || pattern == file
	if negationSet {
		if isGlobbed {
			return false
		}
	}
	return isGlobbed || checkDir(pattern, file)
}

func buildCache(data string) {
	newData := strings.Split(data, "\n")
	if reflect.DeepEqual(newData, IgnoreData) {
		return
	} else {
		IgnoreData = nil
		for _, pattern := range newData {
			IgnoreData = append(IgnoreData, pattern)
		}
	}
}

func PfsIgnore(filePath string) bool {
	if IgnoreExists() {
		interval := 10 * time.Second
		mapPattern, ok := FoundGlobs[filePath]
		if ok && mapPattern.solved.Before(FileLastUpdated.Add(interval)) {
			return FoundGlobs[filePath].value
		}
		//No Valid Data in Solved Cache //checking file Cache
		if time.Now().After(FileLastUpdated.Add(interval)) {
			//If We cant read the file we may as well use the old version.
			data, err := ioutil.ReadFile(IgnoreFile)
			if err != nil {
				Log.Error("Can not read .pfsignore file:", err)
			} else {
				buildCache(string(data))
				FileLastUpdated = time.Now()
			}
		}
		shouldGlob := false
		negationSet := false
		for _, pattern := range IgnoreData {
			if pattern != "" {
				if string(pattern[0]) == "!" {
					negationSet = true
				}
				globResponse := Glob(pattern, filePath)
				shouldGlob = shouldGlob || globResponse
			}
		}
		FoundGlobs[filePath] = SolvedPaths{value: shouldGlob, solved: time.Now()}
		if negationSet && shouldGlob {
			return !shouldGlob
		}
	}
	return false
}
