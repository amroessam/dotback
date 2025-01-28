package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dotback",
	Short: "DotBack - Backup and restore your dotfiles using GitHub",
	Long: `DotBack is a command-line tool that helps you backup and restore your
dotfiles and application configurations using GitHub as storage.
It provides an easy way to manage multiple machine configurations
and keep your development environment in sync.`,
}

func init() {
	// Commands will be added here in subsequent steps
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
