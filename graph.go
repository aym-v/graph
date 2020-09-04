package graph

// Graph represents a graph composed of vertices.
type Graph struct {
	vertices map[string]*Vertex
}

// New creates a new graph without vertices.
func New() *Graph {
	return &Graph{
		vertices: map[string]*Vertex{},
	}
}

// AddVertex adds a vertex and a corresponding workload to the graph
func (g *Graph) AddVertex(id string, task func(...[]byte) []byte) {
	g.vertices[id] = newVertex(task)
}

// Connect establish a communication channel between the given vertex ids.
func (g *Graph) Connect(src, target string) error {
	ch := make(chan []byte)

	if src == target {
		return &LoopError{VertexID: src}
	}

	if v, ok := g.vertices[src]; ok {
		v.addOutputs(ch)
	} else {
		return &UnknownVertexError{
			VertexID:      src,
			ValidVertexes: g.getVertexIDs(),
		}
	}
	if v, ok := g.vertices[target]; ok {
		v.addInputs(ch)
	} else {
		return &UnknownVertexError{
			VertexID:      target,
			ValidVertexes: g.getVertexIDs(),
		}
	}
	return nil
}

// Run runs the graph with the given input and returns the output.
func (g *Graph) Run(in []byte) (chan []byte, error) {
	ic := make(chan []byte)
	oc := make(chan []byte)

	g.addInput(ic)
	g.addOutput(oc)

	err := g.validate()
	if err != nil {
		return nil, err
	}

	for _, v := range g.vertices {
		go v.run()
	}

	ic <- in

	return oc, nil
}

// getVertexIDs returns a list of all the available vertices.
func (g *Graph) getVertexIDs() []string {
	var ids []string
	for k := range g.vertices {
		ids = append(ids, k)
	}
	return ids
}

// addInput connects head input ports with the given input source.
func (g *Graph) addInput(ch chan []byte) error {
	for _, v := range g.vertices {
		if v.In == nil {
			v.addInputs(ch)
			return nil
		}
	}
	return &InputPortError{}
}

// addOutput connects leaf output ports with the given output source.
func (g *Graph) addOutput(ch chan []byte) error {
	found := false
	for _, v := range g.vertices {
		if v.Out == nil {
			v.addOutputs(ch)
			found = true
		}
	}

	if found {
		return nil
	}

	return &OutputPortError{}
}

// Validate checks that every port is connected to a channel.
func (g *Graph) validate() error {
	for id, v := range g.vertices {
		if v.In == nil {
			return &LeafVertexError{
				VertexID: id,
				Port:     "in",
			}
		}
		if v.Out == nil {
			return &LeafVertexError{
				VertexID: id,
				Port:     "out",
			}
		}
	}
	return nil
}
