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

func Glob(pattern, file string) bool {
	//Removing trailing slashes
	pattern = strings.TrimSuffix(pattern, "/")

	if pattern == star || pattern == recursiveStar {
		return true
	}

	if pattern == "" {
		return false
	}

	isGlobbed := starGlob(pattern, file) || dualStarGlob(pattern, file) || pattern == file

	return isGlobbed
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
			globs := strings.Split(string(returnData), "\n")
			for _, pattern := range globs {
				if pattern != "" {
					if string(pattern[0]) == negation {
						for _, globPattern := range globs {
							negationSet = true
							if strings.HasPrefix(pattern, globPattern) && pattern != globPattern {
								negationSet = false // if a directory is ignored negation is removed
								break
							}
						}
						pattern = strings.Trim(pattern, negation)
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
		Log.Info("File", filePath, "has been ignored")
	}
	return shouldGlob
}
