package glob

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfsd/globals"
	"strings"
)

type SolvedPaths struct {
	value bool
	path  string
}

const recursiveStar string = "**"
const star string = "*"
const negation string = "!"

var Log *logger.ParanoidLogger

func ignoreExists() bool {
	_, err, _ := commands.StatCommand(globals.ParanoidDir, ".pfsignore")
	if err != nil {
		Log.Info("No Ignore File Found:", err)
		return false
	}
	return true
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

	for i := 0; i < len(patternArr)-1; i++ {
		if patternArr[i] == "" || fileArr[i] == "" {
			break
		}
		dirsMatch = dirsMatch && (patternArr[i] == fileArr[i])
	}

	return dirsMatch
}

func Glob(pattern, file string) bool {
	if string(pattern[0]) == "!" {
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

	return isGlobbed || checkDir(pattern, file)
}

func ShouldIgnore(filePath string) bool {
	shouldGlob := false
	isIgnore := strings.HasSuffix(strings.TrimSuffix(filePath, "/"), ".pfsignore")
	if !isIgnore {
		if ignoreExists() {
			_, err, returnData := commands.ReadCommand(globals.ParanoidDir, ".pfsignore", -1, -1)
			if err != nil {
				return false
			}
			negationSet := false
			for _, pattern := range strings.Split(string(returnData), "\n") {
				if pattern != "" {
					if string(pattern[0]) == "!" {
						negationSet = true
					}
					globResponse := Glob(pattern, filePath)
					shouldGlob = shouldGlob || globResponse
					Log.Info(pattern, shouldGlob)
				}
			}
			if negationSet && shouldGlob {
				shouldGlob = !shouldGlob
			}
		}
	}
	if shouldGlob {
		Log.Info("File", filePath, "is being Ignored")
	}
	return shouldGlob
}
