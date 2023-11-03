package version

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

// Version is the current version of the application.
// Set via linker flags.
var Version string

// GetShort returns the short version string.
func GetShort() string {
	if Version == "" {
		return "(devel)"
	}

	return Version
}

// GetFull returns the full version string.
func GetFull(prog string) string {
	short := GetShort()

	var buildInfos []string

	info, ok := debug.ReadBuildInfo()
	if ok {
		buildInfos = append(buildInfos, info.GoVersion)
		var commit string
		var modified bool
		for _, settings := range info.Settings {
			switch settings.Key {
			case "vcs.revision":
				commit = settings.Value
			case "vcs.modified":
				modified = settings.Value == "true"
			}
		}
		if commit == "" {
			buildInfos = append(buildInfos, "commit=devel")
		} else {
			s := "commit=" + commit
			if modified {
				s += "-modified"
			}
			buildInfos = append(buildInfos, s)
		}
	} else {
		buildInfos = append(buildInfos, "unknown")
	}

	return fmt.Sprintf("%s version %s (%s)", prog, short, strings.Join(buildInfos, " "))
}

type Flag struct{ Prog string }

func (Flag) IsBoolFlag() bool { return true }

func (Flag) Get() interface{} { return nil }

func (Flag) String() string { return "" }

func (f Flag) Set(s string) error {
	var version string

	if s == "full" {
		version = GetFull(f.Prog)
	} else {
		version = GetShort()
	}

	fmt.Println(version)
	os.Exit(0)
	return nil
}
