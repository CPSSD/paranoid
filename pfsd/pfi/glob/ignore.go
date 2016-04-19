package glob

import (
	"encoding/json"
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfsd/globals"
	"io/ioutil"
	"path"
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

func isIgnoreFile(filePath string) bool {
	return strings.HasSuffix(strings.TrimSuffix(filePath, "/"), ".pfsignore")
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

func previouslyIgnored(filePath string) (bool, *commands.Inode) {
	namePath := commands.GetParanoidPath(globals.ParanoidDir, filePath)
	inodeName, _, err := commands.GetFileInode(namePath)
	if err != nil {
		Log.Error("Error Reading Inode:", err)
		return false, nil
	}
	inodeBytes, err := ioutil.ReadFile(path.Join(path.Join(globals.ParanoidDir, "inodes", string(inodeName))))
	if err != nil {
		Log.Error("Cannot Read Inode:", string(inodeName))
		return false, nil
	}
	inodeData := &commands.Inode{}
	err = json.Unmarshal(inodeBytes, &inodeData)
	if err != nil {
		Log.Error("Cannot Parse iNode Json,", err)
		return false, nil
	}
	return inodeData.Ignored, inodeData
}

func updateInode(inodeData commands.Inode, filePath string, newIgnoreVal bool) {
	if inodeData.Ignored != newIgnoreVal {
		Log.Info("Updating iNode for:", filePath)
		inodeData.Ignored = newIgnoreVal
		commands.UpdateInode(globals.ParanoidDir, filePath, inodeData)
	}
}

func ShouldIgnore(filePath string, changeInode bool) bool {
	shouldGlob := false

	var prevIgnore bool
	var nodeData *commands.Inode

	if !isIgnoreFile(filePath) {

		prevIgnore, nodeData = previouslyIgnored(filePath)
		if ignoreExists() {
			negationSet := false
			_, err, returnData := commands.ReadCommand(globals.ParanoidDir, ".pfsignore", -1, -1)
			if err != nil {
				Log.Error("Cannot Read .pfsignore File, Defaulting to sending over the network")
				return false
			}
			globs := strings.Split(string(returnData), "\n")
			for _, pattern := range globs {
				if pattern != "" {
					if string(pattern[0]) == Negation {
						negationSet = !dirIgnored(pattern, globs)
						pattern = strings.Trim(pattern, Negation)
					}

					globResponse := Glob(pattern, filePath)
					shouldGlob = shouldGlob || globResponse
				}
			}
			if negationSet && shouldGlob {
				shouldGlob = !shouldGlob
			}
		}
	}
	if shouldGlob {
		Log.Info("File:", filePath, "has been ignored")
		if nodeData != nil {
			updateInode(*nodeData, filePath, shouldGlob)
		}
	} else if prevIgnore && nodeData != nil && !changeInode {
		Log.Error(filePath, "has been previously Ignored, and will not sync")
		shouldGlob = true
	} else if changeInode && prevIgnore {
		updateInode(*nodeData, filePath, false) //setting ignore to false
	}
	return shouldGlob
}
