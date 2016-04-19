package glob

import (
	"strings"
)

const recursiveStar string = "**"
const star string = "*"
const Negation string = "!"

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
	isGlobbed := false
	if strings.Contains(pattern, star) {
		isGlobbed = starGlob(pattern, file) || dualStarGlob(pattern, file)
	}

	return isGlobbed || pattern == file
}
