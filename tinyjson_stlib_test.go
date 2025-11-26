//go:build !wasm

package tinyjson_test

import (
	"testing"

	"github.com/cdvelop/tinyjson"
)

func TestStdlib(t *testing.T) {
	j := tinyjson.New()

	t.Run("Encode", func(t *testing.T) { EncodeShared(t, j) })
	t.Run("Decode", func(t *testing.T) { DecodeShared(t, j) })
}
