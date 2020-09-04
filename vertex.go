package graph

// Vertex represents a node in the graph.
type Vertex struct {
	In   []chan []byte
	Out  []chan []byte
	task func(...[]byte) []byte
}

func newVertex(task func(...[]byte) []byte) *Vertex {
	return &Vertex{
		In:   nil,
		Out:  nil,
		task: task,
	}
}

// run executes the task associated with the given vertex.
func (v *Vertex) run() {
	var inputs [][]byte

	// wait for each input channel to send a value
	for _, ch := range v.In {
		inputs = append(inputs, <-ch)
	}

	res := v.task(inputs...)

	// deliver response to all destinations
	for _, out := range v.Out {
		out <- res
		close(out)
	}
}

// addInput adds inputs channels to the vertex
func (v *Vertex) addInputs(ch ...chan []byte) {
	v.In = ch
}

// addOutput adds outputs channels to the vertex
func (v *Vertex) addOutputs(ch ...chan []byte) {
	v.Out = ch
}
