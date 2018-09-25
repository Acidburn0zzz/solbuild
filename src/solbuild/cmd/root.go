//
// Copyright © 2016-2018 Ikey Doherty <ikey@solus-project.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cmd

import (
	"builder"
	"github.com/spf13/cobra"
	"os"
)

// Shared between most of the subcommands
var profile string

// CLIDebug determines whether to enable debug level logs or not.
var CLIDebug bool

// RootCmd is the main entry point into solbuild
var RootCmd = &cobra.Command{
	Use:   "solbuild",
	Short: "solbuild is the Solus package builder",
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "", "Build profile to use")
	RootCmd.PersistentFlags().BoolVarP(&CLIDebug, "debug", "d", false, "Enable debug messages")
	RootCmd.PersistentFlags().BoolVarP(&builder.DisableColors, "no-color", "n", false, "Disable color output")
}

// FindLikelyArg will look in the current directory to see if common path names exist,
// for when it is acceptable to omit a filename.
func FindLikelyArg() string {
	lookPaths := []string{
		"package.yml",
		"pspec.xml",
	}
	for _, p := range lookPaths {
		if st, err := os.Stat(p); err == nil {
			if st != nil {
				return p
			}
		}
	}
	return ""
}
