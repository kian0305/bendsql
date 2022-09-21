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

package build

import (
	"fmt"
	"github.com/cli/safeexec"
	"io"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

// Version is dynamically set by the toolchain or overridden by the Makefile.
var Version = version()

// Date is dynamically set at build time in the Makefile.
var Date = date() // YYYY-MM-DD

func init() {
	if Version == "DEV" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}

func version() string {
	if versionEnv := os.Getenv("BENDSQL_VERSION"); versionEnv != "" {
		return versionEnv
	}
	if desc, err := cmdOutput("git", "describe", "--tags"); err == nil {
		fmt.Println(desc)
		return desc
	}
	rev, _ := cmdOutput("git", "rev-parse", "--short", "HEAD")
	fmt.Println(rev)
	return rev
}

func date() string {
	t := time.Now()
	if sourceDate := os.Getenv("SOURCE_DATE_EPOCH"); sourceDate != "" {
		if sec, err := strconv.ParseInt(sourceDate, 10, 64); err == nil {
			t = time.Unix(sec, 0)
		}
	}
	return t.Format("2006-01-02")
}

func cmdOutput(args ...string) (string, error) {
	exe, err := safeexec.LookPath(args[0])
	if err != nil {
		return "", err
	}
	cmd := exec.Command(exe, args[1:]...)
	cmd.Stderr = io.Discard
	out, err := cmd.Output()
	return strings.TrimSuffix(string(out), "\n"), err
}
