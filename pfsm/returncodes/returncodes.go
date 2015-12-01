package returncodes

import (
	"strconv"
)

//Current pfsm supported return codes
const (
	OK        = iota
	ENOENT    //No such file or directory.
	EACCES    //Can not access file
	EEXIST    //File already exists
	ENOTEMPTY //Directory not empty
	EISDIR    //Is Directory
	ENOTDIR   //Isn't Directory
)

//Gets the integer return code for a given Enum of the code represented as a 2 byte string.
func GetReturnCode(code int) string {
	strcode := strconv.Itoa(code)
	if len(strcode) < 2 {
		return "0" + strcode
	}
	return strcode
}
