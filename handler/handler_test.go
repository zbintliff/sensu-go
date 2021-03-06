package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncAddHandler(t *testing.T) {
	handler := NewMessageHandler()
	handler.AddHandler("MessageType", func(payload []byte) error {
		assert.EqualValues(t, []byte{0}, payload)
		return nil
	})

	assert.NoError(t, handler.Handle("MessageType", []byte{0}))
}

func BenchmarkHandleMessage(b *testing.B) {
	handler := NewMessageHandler()
	handler.AddHandler("MessageType", func(payload []byte) error {
		return nil
	})

	for n := 0; n < b.N; n++ {
		_ = handler.Handle("MessageType", []byte{0})
	}
}
