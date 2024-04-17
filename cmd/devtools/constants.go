package devtools

const (
	DataDirName         = "data"
	DefaultInstanceName = "default"
	BinariesDirName     = "bin"
	LocalConfigDirName  = "config-local"
	RemoteConfigDirName = "config-remote"
)

// Flag variables for use in the different run commands
var isRunLocal bool
var isRunRemote bool
var exportLogs bool
var monoRepoPath string
