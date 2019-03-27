package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints version of the service builder",
	Long: `for example:

$ servicebuilder version 
	`,
	Run: func(c *cobra.Command, args []string) {
		fmt.Printf("%s\n Version:  %s\n Git Commit:  %s\n Go Version:  %s\n OS/Arch:  %s/%s\n Built:  %s\n",
			rootCmd.Use, version, gitCommit,
			runtime.Version(), runtime.GOOS, runtime.GOARCH, compiledAt)
	},
}

var (
	version    = "unknown"
	gitCommit  = "unknown"
	compiledAt = "unknown"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
