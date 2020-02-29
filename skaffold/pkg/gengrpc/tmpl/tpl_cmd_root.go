package tmpl

const CmdRootTpl = `
package cmd

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "",
    Long:  "",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "INFO")
}
`
