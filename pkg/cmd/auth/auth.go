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

package auth

import (
	authConfigureCmd "github.com/databendcloud/bendsql/pkg/cmd/auth/configure"
	authLoginCmd "github.com/databendcloud/bendsql/pkg/cmd/auth/login"
	authTokenCmd "github.com/databendcloud/bendsql/pkg/cmd/auth/token"
	"github.com/databendcloud/bendsql/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth <command>",
		Short: "Authenticate bendsql",
		Annotations: map[string]string{
			"IsCore": "true",
		},
	}

	cmd.AddCommand(authLoginCmd.NewCmdLogin(f, nil))
	cmd.AddCommand(authTokenCmd.NewCmdAuthToken(f, nil))
	cmd.AddCommand(authConfigureCmd.NewCmdConfigure(f))
	return cmd
}
