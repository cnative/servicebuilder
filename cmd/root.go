package cmd

import (
	"fmt"
	"os"

	"github.com/cnative/servicebuilder/internal/term"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	debug  = false
	silent = false
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "servicebuilder",
	Short: "Service builder assists in building a cloud native service.",
	Long: `Service builder enables developers to get a micro service up and 
running by providing a standard and consistent service development. 

For example:

$ servicebuilder new --name contacts --repo $GOPATH/src/github.com/kustomers

This will generate the all required scafolding to get started with the service at

$GOPATH/src/github.com/kustomers/contacts

This application is a tool to generate the needed files
to quickly create a cloud native micro service.`,
	PersistentPreRunE: preRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(semVar, commit, time string) {
	version = semVar
	gitCommit = commit
	compiledAt = time

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "print debug information")
	rootCmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "silent. no verbose output")
}

// preRun
func preRun(c *cobra.Command, args []string) error {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&term.TextFormatter{})
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if silent {
		log.SetLevel(log.FatalLevel)
	}

	return nil
}
