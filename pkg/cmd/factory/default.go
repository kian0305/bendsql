package factory

import (
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/datafuselabs/bendcloud-cli/pkg/iostreams"
)

func New(appVersion string) *cmdutil.Factory {
	f := &cmdutil.Factory{
		ExecutableName: "bendctl",
	}

	f.IOStreams = ioStreams(f) // Depends on Config

	return f
}

func ioStreams(f *cmdutil.Factory) *iostreams.IOStreams {
	io := iostreams.System()

	return io
}
