package main

import (
	"testing"

	"github.com/function61/gokit/assert"
)

func TestEncodeAndDecode(t *testing.T) {
	encoded := RleSliceEncode([]string{"foo", "bar", "lintu sanoo kva|ak"})

	assert.EqualString(t, encoded, "AwADABIA foo|bar|lintu sanoo kva|ak")

	decoded, err := RleSliceDecode(encoded)
	assert.True(t, err == nil)

	assert.True(t, len(decoded) == 3)
	assert.EqualString(t, decoded[0], "foo")
	assert.EqualString(t, decoded[1], "bar")
	assert.EqualString(t, decoded[2], "lintu sanoo kva|ak")
}
