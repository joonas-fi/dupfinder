package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"regexp"
	"strings"
)

// FIXME: this is too complicated. should probably just use "," as an item separator and
//        do c-style escape with "\," if that literal char needs to be used

const (
	separator = "|"
)

func RleSliceEncode(items []string) string {
	buf := &bytes.Buffer{}

	if len(items) == 0 {
		panic("need > 0 items")
	}

	lens := []uint16{}

	for _, item := range items {
		lens = append(lens, uint16(len(item)))
	}

	if err := binary.Write(buf, binary.LittleEndian, lens); err != nil {
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(buf.Bytes()) + " " + strings.Join(items, separator)
}

var stuffRe = regexp.MustCompile("^([^ ]+) (.+)")

func RleSliceDecode(serialized string) ([]string, error) {
	match := stuffRe.FindStringSubmatch(serialized)
	if match == nil {
		return nil, errors.New("no match")
	}

	lensBytes, err := base64.RawURLEncoding.DecodeString(match[1])
	if err != nil {
		return nil, err
	}

	lens := make([]uint16, len(lensBytes)/2)
	if err := binary.Read(bytes.NewBuffer(lensBytes), binary.LittleEndian, lens); err != nil {
		return nil, err
	}

	nextStart := 0

	out := make([]string, len(lens))

	for idx, len_ := range lens {
		out[idx] = match[2][nextStart : nextStart+int(len_)]

		nextStart += int(len_) + len(separator)
	}

	return out, nil
}
