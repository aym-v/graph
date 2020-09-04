package graph

import (
	"reflect"
)

// Encoder implements a generic encoder.
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, vPtr interface{}) error
}

// Task represents a generic task.
//
// Use the following signature: func(p *object) *object
type Task interface{}

// EncodedGraph represents an encoded graph.
type EncodedGraph struct {
	types map[string]reflect.Type

	Graph *Graph
	Enc   Encoder
}

// NewEncoded creates a new encoded graph.
func NewEncoded(g *Graph, e Encoder) *EncodedGraph {
	return &EncodedGraph{
		types: make(map[string]reflect.Type),
		Graph: g,
		Enc:   e,
	}
}

// AddVertex adds a vertex and a corresponding workload to the graph
// task should use the following signature:
// func(...*obj) *obj
func (eg *EncodedGraph) AddVertex(id string, task Task) error {
	taskTyp, err := FuncInfo(task)
	if err != nil {
		return err
	}

	argTyp := taskTyp.In(0)

	eg.types[id] = taskTyp

	taskVal := reflect.ValueOf(task)

	fn := func(in ...[]byte) []byte {
		oV := make([]reflect.Value, len(in))
		oPtr := make([]reflect.Value, len(in))

		for i, arg := range in {
			if argTyp.Kind() != reflect.Ptr {
				oPtr[i] = reflect.New(argTyp)
			} else {
				oPtr[i] = reflect.New(argTyp.Elem())
			}

			if err := eg.Enc.Decode(arg, oPtr[i].Interface()); err != nil {
				panic("graph: Error while decoding task arguments")
			}

			if argTyp.Kind() != reflect.Ptr {
				oPtr[i] = reflect.Indirect(oPtr[i])
			}

			oV[i] = oPtr[i]

		}

		res := taskVal.Call(oV)

		// Encode result
		b, err := eg.Enc.Encode(res[0].Interface())
		if err != nil {
			panic("graph: error while decoding task return value")
		}

		return b
	}

	eg.Graph.vertices[id] = newVertex(fn)
	return nil
}

// Run runs the graph with the given input and returns the output.
func (eg *EncodedGraph) Run(in interface{}) (chan []byte, error) {
	// check that head nodes are compatible with the input type
	if err := eg.checkGraphInput(in); err != nil {
		return nil, err
	}

	b, err := eg.Enc.Encode(in)
	if err != nil {
		return nil, err
	}

	return eg.Graph.Run(b)
}

// Connect checks that the target vertices are communicating with the same types.
func (eg *EncodedGraph) Connect(src, target string) error {
	var out reflect.Type
	var in reflect.Type

	if v, ok := eg.types[src]; ok {
		out = v.Out(0)
	}

	if v, ok := eg.types[target]; ok {
		in = v.In(0)
	}

	if out != in {
		return &IncompatibleVerticesError{
			InVertex:  src,
			OutVertex: target,
		}
	}

	eg.Graph.Connect(src, target)

	return nil
}

func (eg *EncodedGraph) checkGraphInput(in interface{}) error {
	for i, v := range eg.Graph.vertices {
		if v.In == nil {
			inTyp := reflect.TypeOf(in)
			outTyp := eg.types[i].In(0)

			if inTyp != outTyp {
				return &InvalidInputTypeError{
					in: inTyp.String(), out: outTyp.String(),
				}
			}
		}
	}
	return nil
}

// FuncInfo dissects the task signature.
func FuncInfo(t Task) (reflect.Type, error) {
	tType := reflect.TypeOf(t)
	if tType.Kind() != reflect.Func {
		return nil, &WrongTaskTypeError{Kind: tType.Kind()}
	}

	if tType.NumIn() == 0 {
		return nil, &InvalidArgumentNumber{}
	}

	return tType, nil
}
