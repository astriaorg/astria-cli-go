package devtools

const (
	DataDirName         = "data"
	DefaultInstanceName = "default"
	BinariesDirName     = "bin"
	LocalConfigDirName  = "config-local"
	RemoteConfigDirName = "config-remote"

	// Astria monorepo paths
	// All paths are relative to the root of the monorepo
	AstriaTargetDebugPath = "target/debug"
)

// Flag variables for use in the different run commands
var isRunLocal bool
var isRunRemote bool
var exportLogs bool
var monoRepoPath string
