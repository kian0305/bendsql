package auth

import (
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/datafuselabs/bendcloud-cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"net/http"
)

type LoginOptions struct {
	IO          *iostreams.IOStreams
	HttpClient  func() (*http.Client, error)
	Interactive bool
	Token       string
}

func NewCmdLogin(f *cmdutil.Factory, runF func(*LoginOptions) error) *cobra.Command {
	cmd := &cobra.Command{}
	return cmd
}
