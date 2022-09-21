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
	"fmt"
	"github.com/databendcloud/bendsql/pkg/iostreams"
	"github.com/databendcloud/bendsql/pkg/text"
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

var HelpTopics = map[string]map[string]string{
	"environment": {
		"short": "Environment variables that can be used with bendsql",
		"long": heredoc.Doc(`
			BENDSQL_CONFIG_DIR: the directory where bendsql will store configuration files. Default:
			"$HOME/.config/bendsql".
		`),
	},
	"reference": {
		"short": "A comprehensive reference of all bendsql commands",
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
			"markdown:basename": "bendsql_help_" + topic,
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
	fmt.Fprintf(w, "Usage: bendsql help %s", command.Use)
	return nil
}
