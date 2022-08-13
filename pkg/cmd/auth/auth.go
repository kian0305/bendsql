package auth

import (
	authLoginCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/auth/login"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate bendctl",
		Annotations: map[string]string{
			"IsCore": "true",
		},
	}

	cmd.AddCommand(authLoginCmd.NewCmdLogin(f, nil))

	return cmd
}
