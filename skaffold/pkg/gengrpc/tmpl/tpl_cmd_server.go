package tmpl

const CmdServerTpl = `
package cmd

import (
	"{{ . }}/srv"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "main service",
	Long:  "",
	Run:   srv.MainRun,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringP("httpPort", "h", ":6800", "http service port")
	serverCmd.Flags().StringP("grpcPort", "p", ":6801", "grpc service port")
}
`
