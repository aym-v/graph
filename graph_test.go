package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnexistingVertices(t *testing.T) {
	g := New()

	got := g.Connect("A", "B")

	want := &UnknownVertexError{
		VertexID:      "A",
		ValidVertexes: g.getVertexIDs(),
	}

	assert.Equal(t, want, got)
}

func TestLoop(t *testing.T) {
	g := New()

	g.AddVertex("A", nil)

	got := g.Connect("A", "A")

	want := &LoopError{
		VertexID: "A"}

	assert.Equal(t, want, got)
}

func TestAddInOut(t *testing.T) {
	g := New()

	in := make(chan []byte)
	out := make(chan []byte)

	g.AddVertex("A", nil)

	g.addInput(in)
	g.addOutput(out)

	assert.Equal(t, []chan []byte{in}, g.vertices["A"].In)
	assert.Equal(t, []chan []byte{out}, g.vertices["A"].Out)
}

func TestPortError(t *testing.T) {
	g := New()

	in := make(chan []byte)
	out := make(chan []byte)

	g.AddVertex("A", nil)
	g.AddVertex("B", nil)

	// Create a circuit
	g.Connect("A", "B")
	g.Connect("B", "A")

	gotIn := g.addInput(in)
	gotOut := g.addOutput(out)

	assert.Equal(t, &InputPortError{}, gotIn)
	assert.Equal(t, &OutputPortError{}, gotOut)
}

func TestLeafError(t *testing.T) {
	g := New()

	g.AddVertex("A", nil)
	g.AddVertex("B", nil)
	g.AddVertex("C", nil)

	g.Connect("A", "B")
	g.Connect("B", "C")

	err := g.validate()

	assert.Equal(t, &LeafVertexError{VertexID: "A", Port: "in"}, err)
}

func TestRun(t *testing.T) {
	g := New()

	task := func(b ...[]byte) []byte {
		add := byte(0x00)
		return append(b[0], add)
	}

	g.AddVertex("A", task)
	g.AddVertex("B", task)
	g.AddVertex("C", task)
	g.AddVertex("D", task)
	g.AddVertex("E", task)

	g.Connect("A", "B")
	g.Connect("B", "C")
	g.Connect("C", "D")
	g.Connect("D", "E")

	got, _ := g.Run([]byte{0x00})

	assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, <-got)
}

func TestAllValues(t *testing.T) {
	g := New()

	task := func(b ...[]byte) []byte {
		add := byte(0x00)
		return append(b[0], add)
	}
	g.AddVertex("A", task)
	g.AddVertex("B", task)
	g.AddVertex("C", task)

	g.Connect("A", "B")
	g.Connect("A", "C")

	got, _ := g.Run([]byte{0x00})

	for b := range got {
		assert.Equal(t, []byte{0x00, 0x00, 0x00}, b)
	}
}
