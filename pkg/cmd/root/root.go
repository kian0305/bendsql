// Copyright 2022 Datafuse Labs.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package root

import (
	"github.com/MakeNowJust/heredoc"
	authCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/auth"
	completionCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/completion"
	queryCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/query"
	stageCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/stage"
	versionCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/version"
	warehouseCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/warehouse"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands

func NewCmdRoot(f *cmdutil.Factory, version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bendctl <command> <subcommand> [flags]",
		Short: "Dababend Cloud CLI",
		Long:  `Work seamlessly with Databend Cloud from the command line.`,

		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ bendctl auth login
			$ bendctl warehouse status
			$ bendctl warehouse create
			$ bendctl ls stage
		`),
		Annotations: map[string]string{
			"versionInfo": versionCmd.Format(version, buildDate),
		},
	}

	cmd.SetErr(f.IOStreams.ErrOut) // just let it default to os.Stderr instead

	cmd.Flags().Bool("version", false, "Show bendctl version")
	cmd.PersistentFlags().Bool("help", false, "Show help for command")
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		rootHelpFunc(f, c, args)
	})
	cmd.SetUsageFunc(func(c *cobra.Command) error {
		return rootUsageFunc(f.IOStreams.ErrOut, c)
	})
	cmd.SetFlagErrorFunc(rootFlagErrorFunc)

	// Child commands
	cmd.AddCommand(versionCmd.NewCmdVersion(f, version, buildDate))
	cmd.AddCommand(completionCmd.NewCmdCompletion(f.IOStreams))
	cmd.AddCommand(authCmd.NewCmdAuth(f))
	cmd.AddCommand(warehouseCmd.NewWarehouseCmd(f))
	cmd.AddCommand(stageCmd.NewCmdStage(f))
	cmd.AddCommand(queryCmd.NewCmdQuerySQL(f))
	return cmd
}
