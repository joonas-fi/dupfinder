package main

import (
	"fmt"
	"github.com/function61/gokit/fileexists"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func removeUnchangedFromOld(logReader io.Reader, dry bool) error {
	return removeUnchangedFromX(logReader, 0, dry)
}

func removeUnchangedFromNew(logReader io.Reader, dry bool) error {
	return removeUnchangedFromX(logReader, 1, dry)
}

func removeUnchangedFromX(logReader io.Reader, oldOrNew int, dry bool) error {
	processItem := func(typ recordType, payload string) error {
		switch typ {
		case recordTypeSame:
			paths, err := RleSliceDecode(payload)
			if err != nil {
				panic(err)
			}

			// 0 for old, 1 for new
			path := paths[oldOrNew]

			if dry {
				fmt.Println("would remove " + path)
			} else {
				exists, err := fileexists.Exists(path)
				if err != nil {
					return err
				}

				if exists {
					fmt.Println("removing " + path)

					if err := os.Remove(path); err != nil {
						return fmt.Errorf("Remove: %v", err)
					}
				} else {
					fmt.Println("FILE DOES NOT EXIST " + path)
				}
			}
		case recordTypeAdded, recordTypeRemoved, recordTypeInfo, recordTypeReadError:
			// noop
		default:
			panic("unknown record type")
		}

		return nil
	}

	return parseReport(logReader, processItem)
}

func removeUnchangedFromOldEntrypoint() *cobra.Command {
	really := false

	cmd := &cobra.Command{
		Use:   "removeunchangedfilesfromold",
		Short: "Remove files from old directory that are the same in new directory",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := removeUnchangedFromOld(os.Stdin, !really); err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().BoolVarP(&really, "really", "", really, "Really remove the files. Without this a dry run is performed.")

	return cmd
}

func removeUnchangedFromNewEntrypoint() *cobra.Command {
	really := false

	cmd := &cobra.Command{
		Use:   "removeunchangedfilesfromnew",
		Short: "Remove files from new directory that are the same in old directory",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := removeUnchangedFromNew(os.Stdin, !really); err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().BoolVarP(&really, "really", "", really, "Really remove the files. Without this a dry run is performed.")

	return cmd
}
