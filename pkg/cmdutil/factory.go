package cmdutil

import (
	"github.com/datafuselabs/bendcloud-cli/pkg/iostreams"
	"os"
	"path/filepath"
	"strings"
)

type Factory struct {
	IOStreams *iostreams.IOStreams
	// TODO: add interface for bendctl config
	//Config     func() (config.Config, error)
	ExecutableName string
}

// Executable is the path to the currently invoked binary
func (f *Factory) Executable() string {
	if !strings.ContainsRune(f.ExecutableName, os.PathSeparator) {
		f.ExecutableName = executable(f.ExecutableName)
	}
	return f.ExecutableName
}

// Finds the location of the executable for the current process as it's found in PATH, respecting symlinks.
// If the process couldn't determine its location, return fallbackName. If the executable wasn't found in
// PATH, return the absolute location to the program.
func executable(fallbackName string) string {
	exe, err := os.Executable()
	if err != nil {
		return fallbackName
	}

	base := filepath.Base(exe)
	path := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(path) {
		p, err := filepath.Abs(filepath.Join(dir, base))
		if err != nil {
			continue
		}
		f, err := os.Lstat(p)
		if err != nil {
			continue
		}

		if p == exe {
			return p
		} else if f.Mode()&os.ModeSymlink != 0 {
			if t, err := os.Readlink(p); err == nil && t == exe {
				return p
			}
		}
	}

	return exe
}
