package bfs

// Producer is the interface implemented by the nodes of the graph.
//
// ID() must returns a string that uniquely identify the node.
// Value() should returns the value of the node.
type Producer interface {
	ID() string
	Value() interface{}
}

// Visitor is the interface that wraps the method to walks the graph.
//
// Visit() takes a Producer, visit it, and returns a slice containing
// all the generators reachable in one step.
type Visitor interface {
	Visit(node Producer) ([]Producer, error)
}

// VisitorFunc is an adapter to make Bfs accept functions
type VisitorFunc func(node Producer) ([]Producer, error)

// Visit call the function itself
func (f VisitorFunc) Visit(node Producer) ([]Producer, error) {
	return f(node)
}

// Bfs implements the Breadth First Search algorithm
func Bfs(node Producer, visitor Visitor) ([]Producer, error) {
	visited := make(map[string]bool)
	var result []Producer

	for queue := []Producer{node}; len(queue) > 0; {
		next := queue[0]
		queue = queue[1:]

		if _, found := visited[next.ID()]; found {
			continue
		}

		neighbors, err := visitor.Visit(next)
		if err != nil {
			return nil, err
		}
		visited[next.ID()] = true

		result = append(result, next)

		queue = append(queue, neighbors...)
	}

	return result, nil
}
