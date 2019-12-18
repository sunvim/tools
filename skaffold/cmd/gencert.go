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
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sunvim/tools/skaffold/pkg/constant"
	"github.com/sunvim/tools/skaffold/pkg/gencert"
	"k8s.io/client-go/util/homedir"
)

// gencertCmd represents the gencert command
var gencertCmd = &cobra.Command{
	Use:   "gencert",
	Short: "generate  certificate authority file",
	Long:  ``,
	Run:   gencert.PrepareBasicWorkEnvironments,
}

func init() {
	rootCmd.AddCommand(gencertCmd)
	gencertCmd.Flags().StringP(constant.ServiceName, "n", "web-hook-scaffold", "web hook service name")
	gencertCmd.Flags().StringP(constant.Namespace, "s", "default", "web hook service namespace")
	gencertCmd.Flags().StringP(constant.NsLabel, "l", "webhook-injector",
		"namespace label ex: kubectl label ns default webhook-injector=enabled")
	if home := homedir.HomeDir(); home != "" {
		gencertCmd.Flags().StringP(constant.KubernetesConfig, "k", filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file")
	} else {
		gencertCmd.Flags().StringP(constant.KubernetesConfig, "k", "",
			"absolute path to the kubeconfig file")
	}
}
