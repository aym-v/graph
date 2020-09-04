package graph

import (
	"fmt"
	"reflect"
)

// UnknownVertexError is triggered when an edge has one of its endpoints pointing to nowhere.
type UnknownVertexError struct {
	VertexID      string
	ValidVertexes []string
}

func (e *UnknownVertexError) Error() string {
	return fmt.Sprintf("graph: unknown vertex '%s', valid vertexes are: '%s'", e.VertexID, e.ValidVertexes)
}

// LoopError is triggered by graph loops.
type LoopError struct {
	VertexID string
}

func (e *LoopError) Error() string {
	return fmt.Sprintf("graph: vertex '%s' cannot point to itself", e.VertexID)
}

// InputPortError is returned when there is no input port.
type InputPortError struct{}

func (e *InputPortError) Error() string {
	return fmt.Sprintf("graph: no input ports were found")
}

// OutputPortError is returned when there is no output port.
type OutputPortError struct{}

func (e *OutputPortError) Error() string {
	return fmt.Sprintf("graph: no output ports were found")
}

// LeafVertexError is returned when a vertex has an empty port.
type LeafVertexError struct {
	VertexID string
	Port     string
}

func (e *LeafVertexError) Error() string {
	return fmt.Sprintf("graph: endpoint '%s' on vertex '%s' should be assigned to a channel", e.Port, e.VertexID)
}

// IncompatibleVerticesError is returned when two vertices have incompatible types.
type IncompatibleVerticesError struct {
	InVertex  string
	OutVertex string
}

func (e *IncompatibleVerticesError) Error() string {
	return fmt.Sprintf("graph: vertex '%s' is not compatible with vertex '%s'", e.InVertex, e.OutVertex)
}

// InvalidInputTypeError is returned when the graph input is is not applicable to any of the graph's vertex.
type InvalidInputTypeError struct {
	in  string
	out string
}

func (e *InvalidInputTypeError) Error() string {
	return fmt.Sprintf("graph: the graph input type is not applicable to any of the vertices. '%s' != '%s'", e.in, e.out)
}

// WrongTaskTypeError is returned when the input task is not a function.
type WrongTaskTypeError struct {
	Kind reflect.Kind
}

func (e *WrongTaskTypeError) Error() string {
	return fmt.Sprintf("graph: task should be a function, got '%s'", reflect.Kind(e.Kind))
}

// InvalidArgumentNumber is returned when there is no arguments.
type InvalidArgumentNumber struct{}

func (e *InvalidArgumentNumber) Error() string {
	return fmt.Sprintf("graph: task should have at least one argument.")
}
