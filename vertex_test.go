package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunTask(t *testing.T) {
	s := func(b ...[]byte) []byte {
		add := byte(0x00)
		return append(b[0], add)
	}

	got := newVertex(s)

	in := make(chan []byte)
	out := make(chan []byte)

	got.addInputs(in)
	got.addOutputs(out)

	go got.run()

	in <- []byte{0x00}
	result := <-out

	assert.Equal(t, result, []byte{0x00, 0x00})
}

func TestFanOut(t *testing.T) {
	s := func(b ...[]byte) []byte {
		add := byte(0x00)
		return append(b[0], add)
	}

	got := newVertex(s)

	in := make(chan []byte)
	out := make(chan []byte)
	out2 := make(chan []byte)

	got.addInputs(in)
	got.addOutputs(out, out2)

	go got.run()

	in <- []byte{0x00}
	result := <-out
	result2 := <-out2

	assert.Equal(t, result, []byte{0x00, 0x00})
	assert.Equal(t, result2, []byte{0x00, 0x00})
}

func TestFanIn(t *testing.T) {
	reduce := func(b ...[]byte) []byte {
		var acc [][]byte

		for _, el := range b {
			acc = append(acc, el)
		}

		return acc[1]
	}

	got := newVertex(reduce)

	in := make(chan []byte)
	in2 := make(chan []byte)
	out := make(chan []byte)

	got.addInputs(in, in2)
	got.addOutputs(out)

	go got.run()

	in <- []byte{0x00}
	in2 <- []byte{0x01}
	result := <-out

	assert.Equal(t, result, []byte{0x01})

}
