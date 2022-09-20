package factory

import (
	"github.com/datafuselabs/bendcloud-cli/api"
	"github.com/datafuselabs/bendcloud-cli/internal/config"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/datafuselabs/bendcloud-cli/pkg/iostreams"
)

func New(appVersion string) *cmdutil.Factory {
	f := &cmdutil.Factory{
		ExecutableName: "bendctl",
		ApiClient:      httpClientFunc(),
		Config:         configFunc(),
	}

	f.IOStreams = ioStreams(f)

	return f
}

func ioStreams(f *cmdutil.Factory) *iostreams.IOStreams {
	io := iostreams.System()

	return io
}

func configFunc() func() (config.Configer, error) {
	var cachedConfig config.Configer
	var configError error
	return func() (config.Configer, error) {
		if cachedConfig != nil || configError != nil {
			return cachedConfig, configError
		}
		cachedConfig, configError = config.NewConfig()
		return cachedConfig, configError
	}
}

func httpClientFunc() func() (*api.APIClient, error) {
	return func() (*api.APIClient, error) {
		return api.NewApiClient(), nil
	}
}
