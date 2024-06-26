package main

import (
	"crypto/sha1" // nolint:gosec // ok for non-secure comparisons
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/function61/gokit/fileexists"
	"github.com/spf13/cobra"
)

type State struct {
	oldDirHashes            map[string]string
	hashesMissingFromNewDir map[string]string
}

func NewState() *State {
	return &State{
		oldDirHashes: map[string]string{},
		// - oldDirHashes is cloned to here after old dir is processed.
		// - when new dir has a match, an entry is removed from here.
		// - in the end, the hashes from old dir not present in new dir are what is left here.
		hashesMissingFromNewDir: map[string]string{},
	}
}

func (s *State) initializeMissingMap() {
	for hash, path := range s.oldDirHashes {
		s.hashesMissingFromNewDir[hash] = path
	}
}

func (s *State) markNotMissing(hash string) {
	delete(s.hashesMissingFromNewDir, hash)
}

func stopOnErrors(fn func(path string, info os.FileInfo) error) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// do not process directories (but recurse to them)
		if info.IsDir() {
			return nil
		}

		return fn(path, info)
	}
}

func compare(olddir string, newdir string) error {
	state := NewState()

	r := &report{}

	r.info(fmt.Sprintf("olddir<%s> newdir<%s>", olddir, newdir))

	// quick sanity check to see that new dir also exists. if we checked it after walking,
	// olddir, we might have lots of wasted work
	if newdirExists, err := fileexists.Exists(newdir); !newdirExists || err != nil {
		return errors.New("newdir does not exist, or error accessing it")
	}

	r.info("Scanning olddir")

	if err := filepath.Walk(olddir, stopOnErrors(func(path string, info os.FileInfo) error {
		hash, err := sha1File(path)
		if err != nil {
			return err
		}

		state.oldDirHashes[hash] = path

		return nil
	})); err != nil {
		return err
	}

	r.info("starting initializeMissingMap")

	state.initializeMissingMap()

	r.info("Scanning newdir")

	if err := filepath.Walk(newdir, stopOnErrors(func(newPath string, info os.FileInfo) error {
		hash, err := sha1File(newPath)
		if err != nil {
			// TODO: configurable dontStopOnError
			r.readError(newPath, err)
			return nil
		}

		oldPath, hashExists := state.oldDirHashes[hash]

		if hashExists {
			r.sameFile(oldPath, newPath)
			state.markNotMissing(hash)
		} else {
			r.addedFile(newPath)
		}

		return nil
	})); err != nil {
		return err
	}

	r.info("Listing missing files")

	for _, oldPath := range state.hashesMissingFromNewDir {
		r.missingFile(oldPath)
	}

	r.info("Done")

	return nil
}

func compareEntrypoint() *cobra.Command {
	return &cobra.Command{
		Use:   "compare [olddir] [newdir]",
		Short: "Lists added and removed files in newdir",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := compare(args[0], args[1]); err != nil {
				panic(err)
			}
		},
	}
}

func sha1File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	// nolint:gosec // ok for non-secure comparisons
	hash := sha1.New()

	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	hashHex := fmt.Sprintf("%x", hash.Sum(nil))

	return hashHex, nil
}
