package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev" // replaced dynamically at build time

func main() {
	rootCmd := &cobra.Command{
		Use:     os.Args[0],
		Short:   "dupfinder helps you find and delete duplicate files",
		Version: version,
	}

	rootCmd.AddCommand(compareEntrypoint())
	rootCmd.AddCommand(dirDiffEntrypoint())
	rootCmd.AddCommand(removeEmptyDirsEntrypoint())
	rootCmd.AddCommand(removeUnchangedFromOldEntrypoint())
	rootCmd.AddCommand(removeUnchangedFromNewEntrypoint())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
