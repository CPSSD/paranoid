package returncodes

//Current pfsm supported return codes
const (
	OK          = iota
	ENOENT      //No such file or directory.
	EACCES      //Can not access file
	EEXIST      //File already exists
	ENOTEMPTY   //Directory not empty
	EISDIR      //Is Directory
	EIO         //Input/Output error
	ENOTDIR     //Isn't Directory
	EBUSY       //System is busy
	EUNEXPECTED //Unforseen error
)
