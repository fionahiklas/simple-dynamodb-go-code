package version

// These values are set by "-X" options to the ldflags for Go, see the Makefile
var (
	commitHash  = "null"
	codeVersion = "null"
	appName     = "null"
)

func CommitHash() string  { return commitHash }
func CodeVersion() string { return codeVersion }
func AppName() string     { return appName }
