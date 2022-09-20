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

package cmdutil

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

type enumValue struct {
	string  *string
	options []string
}

func formatValuesForUsageDocs(values []string) string {
	return fmt.Sprintf("{%s}", strings.Join(values, "|"))
}

func (e *enumValue) Set(value string) error {
	if !isIncluded(value, e.options) {
		return fmt.Errorf("valid values are %s", formatValuesForUsageDocs(e.options))
	}
	*e.string = value
	return nil
}

func (e *enumValue) String() string {
	return *e.string
}

func (e *enumValue) Type() string {
	return "string"
}
func isIncluded(value string, opts []string) bool {
	for _, opt := range opts {
		if strings.EqualFold(opt, value) {
			return true
		}
	}
	return false
}

// StringEnumFlag defines a new string flag that only allows values listed in options.
func StringEnumFlag(cmd *cobra.Command, p *string, name, shorthand, defaultValue string, options []string, usage string) *pflag.Flag {
	*p = defaultValue
	val := &enumValue{string: p, options: options}
	f := cmd.Flags().VarPF(val, name, shorthand, fmt.Sprintf("%s: %s", usage, formatValuesForUsageDocs(options)))
	_ = cmd.RegisterFlagCompletionFunc(name, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return options, cobra.ShellCompDirectiveNoFileComp
	})
	return f
}
