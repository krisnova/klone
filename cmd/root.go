// Copyright © 2017 Kris Nova <kris@nivenly.com>
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
//
//  _  ___
// | |/ / | ___  _ __   ___
// | ' /| |/ _ \| '_ \ / _ \
// | . \| | (_) | | | |  __/
// |_|\_\_|\___/|_| |_|\___|
//
// root.go is the cobra command for the primary klone command

package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/kris-nova/klone/pkg/klone"
	"github.com/kris-nova/klone/pkg/local"
	"errors"
)

var RootCmd = &cobra.Command{
	Use:   "klone",
	Short: "Used to clone a git repository with style.",
	Long:  `Klone allows you define custom logic based on repository programming language.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		local.PrintStartBanner()
		if len(args) > 0 {
			err := klone.Klone(args[0])
			if err != nil {
				local.PrintError(err)
			}
		} else {
			local.PrintErrorExitCode(errors.New("Missing argument. Use: 'klone $name'"), 99)
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.klone.yaml)")
}
