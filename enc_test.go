package graph

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valaymerick/braidboard/internal/graph/encoders"
	"github.com/valaymerick/braidboard/internal/graph/encoders/testdata"
)

type A struct {
	foo string
}

type B struct {
	foo int
}

func TestFuncInfo(t *testing.T) {
	tests := []struct {
		in  interface{}
		typ reflect.Kind
	}{
		{func(i, j string) {}, reflect.String},
		{func(i string) {}, reflect.String},
	}

	for _, tt := range tests {
		typ, _ := FuncInfo(tt.in)
		assert.Equal(t, typ.In(0).Kind(), tt.typ)
	}
}

func TestAddVertex_NoArgs(t *testing.T) {
	g := New()
	eg := NewEncoded(g, encoders.NewProtobuf())

	err := eg.AddVertex("A", func() {})

	assert.Equal(t, true, err != nil)
}

func TestConnect(t *testing.T) {
	g := New()
	eg := NewEncoded(g, encoders.NewProtobuf())

	strToStr := func(a string) string { return "a" }
	intToStr := func(a int) string { return "a" }
	atoA := func(a *A) *A { return &A{} }
	btoB := func(p *B) *B { return &B{} }

	eg.AddVertex("A", strToStr)
	eg.AddVertex("B", strToStr)
	eg.AddVertex("C", intToStr)
	eg.AddVertex("D", intToStr)

	eg.AddVertex("E", atoA)
	eg.AddVertex("F", atoA)
	eg.AddVertex("G", btoB)

	tests := []struct {
		src    string
		target string
		err    bool
	}{
		{"A", "B", false},
		{"B", "C", true},
		{"C", "D", true},

		{"E", "F", false},
		{"F", "G", true},
	}

	for _, tt := range tests {
		err := eg.Connect(tt.src, tt.target)
		assert.Equal(t, tt.err, err != nil)
	}
}

func TestCheckGraphInput(t *testing.T) {
	tests := []struct {
		task    interface{}
		graphIn interface{}
		err     bool
	}{
		{func(a int) int { return 0 }, 1, false},
		{func(a string) int { return 0 }, 1, true},

		{func(a *A) int { return 0 }, &A{}, false},
		{func(a A) int { return 0 }, A{}, false},
		{func(a A) int { return 0 }, &A{}, true},
		{func(a A) int { return 0 }, &B{}, true},
	}

	for _, tt := range tests {
		g := New()
		eg := NewEncoded(g, encoders.NewProtobuf())

		eg.AddVertex("A", tt.task)
		err := eg.checkGraphInput(tt.graphIn)

		assert.Equal(t, tt.err, err != nil)
	}
}

func a() interface{} {
	return nil
}

func TestGraphRun(t *testing.T) {
	g := New()
	eg := NewEncoded(g, encoders.NewProtobuf())

	exp := &testdata.Person{Name: "foo"}

	task := func(a *testdata.Person) *testdata.Person { return exp }

	eg.AddVertex("A", task)

	eg.Connect("A", "B")

	p := &testdata.Person{}

	res, _ := eg.Run(p)
	eg.Enc.Decode(<-res, p)

	assert.EqualValues(t, exp.ProtoReflect(), p.ProtoReflect())
}
