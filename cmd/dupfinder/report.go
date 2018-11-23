package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
)

type recordType string

const (
	recordTypeAdded     recordType = "+"
	recordTypeRemoved              = "-"
	recordTypeSame                 = "="
	recordTypeInfo                 = "INFO"
	recordTypeReadError            = "!"
)

type report struct{}

// info messages go both to stderr and stdout
func (r *report) info(output string) {
	r.write(recordTypeInfo, output)

	log.Println(output)
}

func (r *report) sameFile(oldPath string, newPath string) {
	r.write(recordTypeSame, RleSliceEncode([]string{oldPath, newPath}))
}

func (r *report) missingFile(path string) {
	r.write(recordTypeRemoved, path)
}

func (r *report) addedFile(path string) {
	r.write(recordTypeAdded, path)
}

func (r *report) readError(path string, err error) {
	descr := fmt.Sprintf("%s %s", path, err.Error())

	r.write(recordTypeReadError, descr)

	log.Println("readError: " + descr)
}

func (r *report) write(typ recordType, line string) {
	fmt.Println(string(typ) + " " + line)
}

var recordParseRe = regexp.MustCompile("^([^ ]+) (.+)")

type recordCallback func(recordType, string) error

func parseReport(reader io.Reader, cb recordCallback) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		match := recordParseRe.FindStringSubmatch(scanner.Text())
		if match == nil {
			panic("match nil")
		}

		// allow callback to interrupt log processing
		if err := cb(recordType(match[1]), match[2]); err != nil {
			return err
		}
	}

	return scanner.Err()
}
