package glob

import (
	"fmt"
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfsd/globals"
	"path/filepath"
	"strings"
)

var Log *logger.ParanoidLogger

func ignoreExists() bool {
	_, err, _ := commands.StatCommand(globals.ParanoidDir, ".pfsignore")
	if err != nil {
		Log.Info("No Ignore File Found:", err)
		return false
	}
	return true
}

func isIgnoredFile(filePath string) bool {
	return filepath.Base(filePath) != ".pfsignore"
}

func dirIgnored(currentPattern string, globs []string) bool {
	dirFound := false
	for _, globPattern := range globs {
		if strings.HasPrefix(currentPattern, globPattern) && currentPattern != globPattern {
			dirFound = true // if a directory is ignored Negation is removed
			break
		}
	}
	return dirFound
}

func ShouldIgnore(filePath string, changeInode bool) (bool, error) {
	shouldIgnore := false

	var prevIgnore bool
	var code returncodes.Code

	if isIgnoredFile(filePath) {
		prevIgnore, code = commands.PreviouslyIgnored(globals.ParanoidDir, filePath)
		if ignoreExists() {
			negationSet := false
			_, err, returnData := commands.ReadCommand(globals.ParanoidDir, ".pfsignore", -1, -1)
			if err != nil {
				Log.Error("Cannot Read .pfsignore File, Defaulting to sending over the network")
				return false, fmt.Errorf("Cannot Read .pfsignore File")
			}
			globs := strings.Split(string(returnData), "\n")
			for _, pattern := range globs {
				if pattern != "" {
					dirIgnored := dirIgnored(pattern, globs)
					if string(pattern[0]) == Negation {
						negationSet = !dirIgnored
						pattern = strings.Trim(pattern, Negation)
					}

					globResponse := Glob(pattern, filePath)
					shouldIgnore = shouldIgnore || globResponse || dirIgnored
				}
			}
			if negationSet && shouldIgnore {
				shouldIgnore = !shouldIgnore
			}
		}
	}
	if shouldIgnore {
		Log.Info("File:", filePath, "has been ignored")
		if code != returncodes.OK {
			commands.UpdateInodeIgnore(globals.ParanoidDir, filePath, shouldIgnore)
		}
	} else if prevIgnore && code == returncodes.OK && !changeInode {
		Log.Error(filePath, "was previously ignored and will not sync")
		shouldIgnore = true
	} else if changeInode && prevIgnore {
		commands.UpdateInodeIgnore(globals.ParanoidDir, filePath, false) //setting ignore to false
	}
	return shouldIgnore, nil
}
