package root

import (
	"fmt"
	"github.com/datafuselabs/bendcloud-cli/pkg/iostreams"
	"github.com/datafuselabs/bendcloud-cli/pkg/text"
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var HelpTopics = map[string]map[string]string{
	"environment": {
		"short": "Environment variables that can be used with bendctl",
		"long": heredoc.Doc(`
			BENDCTL_CONFIG_DIR: the directory where bendctl will store configuration files. Default:
			"$HOME/.config/bendctl".
		`),
	},
	"reference": {
		"short": "A comprehensive reference of all bendctl commands",
	},
}

func NewHelpTopic(ios *iostreams.IOStreams, topic string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     topic,
		Short:   HelpTopics[topic]["short"],
		Long:    HelpTopics[topic]["long"],
		Example: HelpTopics[topic]["example"],
		Hidden:  true,
		Annotations: map[string]string{
			"markdown:generate": "true",
			"markdown:basename": "bendctl_help_" + topic,
		},
	}

	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		helpTopicHelpFunc(ios.Out, c, args)
	})
	cmd.SetUsageFunc(func(c *cobra.Command) error {
		return helpTopicUsageFunc(ios.ErrOut, c)
	})

	return cmd
}

func helpTopicHelpFunc(w io.Writer, command *cobra.Command, args []string) {
	fmt.Fprint(w, command.Long)
	if command.Example != "" {
		fmt.Fprintf(w, "\n\nEXAMPLES\n")
		fmt.Fprint(w, text.Indent(command.Example, "  "))
	}
}

func helpTopicUsageFunc(w io.Writer, command *cobra.Command) error {
	fmt.Fprintf(w, "Usage: bendctl help %s", command.Use)
	return nil
}
