package activitylogger

// protobuf types
const (
	TChmod uint32 = iota
	TCreat
	TLink
	TMkdir
	TRename
	TRmdir
	TSymLink
	TTruncate
	TUnlink
	TUtimes
	TWrite
)

// TypeString returns the string version of the type
func TypeString(t uint32) string {
	switch t {
	case TChmod:
		return "Chmod"
	case TCreat:
		return "Creat"
	case TLink:
		return "Link"
	case TMkdir:
		return "Mkdir"
	case TRename:
		return "Rename"
	case TRmdir:
		return "Rmdir"
	case TSymLink:
		return "Symlink"
	case TTruncate:
		return "Truncate"
	case TUnlink:
		return "Unlink"
	case TUtimes:
		return "Utimes"
	case TWrite:
		return "Write"
	default:
		return "Unknown"
	}
}
