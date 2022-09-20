package stage

import (
	stageListCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/stage/stage_ls"
	stageUploadCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/stage/stage_upload"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdStage(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stage <command>",
		Short: "Operate stage",
		Annotations: map[string]string{
			"IsCore": "true",
		},
	}

	cmd.AddCommand(stageUploadCmd.NewCmdStageUpload(f))
	cmd.AddCommand(stageListCmd.NewCmdStageList(f))

	return cmd
}
