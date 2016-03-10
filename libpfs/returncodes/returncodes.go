package returncodes

// Code is the basic type for all returncodes
type Code int

//Current pfsm supported return codes
const (
	OK          Code = iota
	ENOENT           //No such file or directory.
	EACCES           //Can not access file
	EEXIST           //File already exists
	ENOTEMPTY        //Directory not empty
	EISDIR           //Is Directory
	EIO              //Input/Output error
	ENOTDIR          //Isn't Directory
	EBUSY            //System is busy
	EUNEXPECTED      //Unforseen error
)

// String returns the code as a string
func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case ENOENT:
		return "ENOENT"
	case EACCES:
		return "EACCES"
	case EEXIST:
		return "EEXIST"
	case ENOTEMPTY:
		return "ENOTEMPTY"
	case EISDIR:
		return "EISDIR"
	case EIO:
		return "EIO"
	case ENOTDIR:
		return "ENOTDIR"
	case EBUSY:
		return "EBUSY"
	case EUNEXPECTED:
		return "EUNEXPECTED"
	default:
		return ""
	}
}
