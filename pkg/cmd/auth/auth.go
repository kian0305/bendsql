package auth

import (
	authConfigureCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/auth/configure"
	authLoginCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/auth/login"
	authTokenCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/auth/token"
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
	cmd.AddCommand(authTokenCmd.NewCmdAuthToken(f, nil))
	cmd.AddCommand(authConfigureCmd.NewCmdConfigure(f))
	return cmd
}
