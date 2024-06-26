package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func dirDiffEntrypoint() *cobra.Command {
	return &cobra.Command{
		Use:   "dirdiff [olddir] [newdir]",
		Short: "Diff dirs only based on filenames",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := dirDiff(args[0], args[1]); err != nil {
				panic(err)
			}
		},
	}
}

func dirDiff(oldDir string, newDir string) error {
	oldEntries, err := os.ReadDir(oldDir)
	if err != nil {
		return err
	}

	newEntries, err := os.ReadDir(newDir)
	if err != nil {
		return err
	}

	oldNames := names(oldEntries)
	newNames := names(newEntries)

	for _, oldEntry := range oldEntries {
		if !exists(oldEntry.Name(), newNames) {
			fmt.Printf("missing from new: %s\n", filepath.Join(oldDir, oldEntry.Name()))
		} else if oldEntry.IsDir() { // recurse into subdirs
			if err := dirDiff(
				filepath.Join(oldDir, oldEntry.Name()),
				filepath.Join(newDir, oldEntry.Name())); err != nil {
				return err
			}
		}
	}

	for _, newEntry := range newEntries {
		if !exists(newEntry.Name(), oldNames) {
			fmt.Printf("missing from old: %s\n", filepath.Join(newDir, newEntry.Name()))
		}
	}

	return nil
}

func names(infos []fs.DirEntry) []string {
	ret := []string{}

	for _, info := range infos {
		ret = append(ret, info.Name())
	}

	return ret
}

func exists(item string, list []string) bool {
	for _, other := range list {
		if item == other {
			return true
		}
	}

	return false
}
