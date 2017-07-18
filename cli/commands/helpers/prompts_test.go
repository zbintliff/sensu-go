package helpers

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestDeletePrompt(t *testing.T) {
	assert := assert.New(t)
	out := exWriter{}

	in := exReader{input: "TEST\n"}
	confirmed := ConfirmDelete("test", &in, &out)

	assert.True(confirmed)
}

type exReader struct {
	input string
}

func (r *exReader) Read(buf []byte) (int, error) {
	ptr := unsafe.Pointer(&buf)
	newBuf := []byte(r.input)
	ptr = &newBuf
	return len(r.input), nil
}

type exWriter struct {
	result string
}

func (w *exWriter) Clean() {
	w.result = ""
}

func (w *exWriter) Write(p []byte) (int, error) {
	w.result += string(p)
	return 0, nil
}
