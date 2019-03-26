package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints version of the service builder",
	Long: `For example:

$ servicebuilder version 
	`,
	Run: func(c *cobra.Command, args []string) {
		fmt.Printf("%s\n Version:  %s\n Git Commit:  %s\n Go Version:  %s\n OS/Arch:  %s/%s\n Built:  %s\n",
			rootCmd.Use, version, gitCommit,
			runtime.Version(), runtime.GOOS, runtime.GOARCH, CompiledAt().String())
	},
}

var (
	gitCommit = "unknown"
	version   = "dev"
	compiled  = ""
)

func init() {
	rootCmd.AddCommand(versionCmd)
	if compiled == "" {
		compiled = strconv.FormatInt(time.Now().Unix(), 10)
	}
}

// CompiledAt converts the Unix time Compiled to a time.Time using UTC timezone.
func CompiledAt() time.Time {
	i, err := strconv.ParseInt(compiled, 10, 64)
	if err != nil {
		panic(err)
	}

	return time.Unix(i, 0).UTC()
}
