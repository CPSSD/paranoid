package glob

import (
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

func ShouldIgnore(filePath string, changeInode bool) (bool, returncodes.Code) {
	shouldIgnore := false

	var prevIgnore bool
	var code returncodes.Code

	if isIgnoredFile(filePath) {
		prevIgnore, code = commands.PreviouslyIgnored(globals.ParanoidDir, filePath)
		if ignoreExists() {
			negationSet := false
			code, err, returnData := commands.ReadCommand(globals.ParanoidDir, ".pfsignore", -1, -1)
			if err != nil || code != returncodes.OK {
				Log.Error("Cannot Read .pfsignore File, Defaulting to sending over the network")
				return false, code
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
			code, err := commands.UpdateInodeIgnore(globals.ParanoidDir, filePath, shouldIgnore)
			if err != nil && code != returncodes.OK {
				Log.Error("Cannot Update iNode for", filePath)
			}
		}
	} else if prevIgnore && code == returncodes.OK && !changeInode {
		Log.Warn(filePath, "was previously ignored and will not sync")
		shouldIgnore = true
	} else if changeInode && prevIgnore {
		code, err := commands.UpdateInodeIgnore(globals.ParanoidDir, filePath, false) //setting ignore to false
		if err != nil && code != returncodes.OK {
			Log.Error("Cannot Update iNode for", filePath, err)
		}
	}
	return shouldIgnore, returncodes.OK
}
