/*
Copyright © 2019 Mobius <sv0220@163.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sunvim/tools/skaffold/pkg/constant"
	"github.com/sunvim/tools/skaffold/pkg/gencode/frame"
)

// gencodeCmd represents the gencode command
var gencodeCmd = &cobra.Command{
	Use:   "gencode",
	Short: "generate web-hook framework",
	Long:  ``,
	Run:   frame.GenCode,
}

func init() {
	rootCmd.AddCommand(gencodeCmd)

	gencodeCmd.Flags().StringP(constant.PkgName, "p", "frame", "package name")
	gencodeCmd.Flags().StringP(constant.ServiceName, "n", "web-hook-scaffold", "service name")
	gencodeCmd.Flags().StringP(constant.CompanyName, "c", "scaffold", "company name in docker hub")
}
