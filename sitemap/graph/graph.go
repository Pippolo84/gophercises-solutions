package graph

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
// It visits p and all the nodes reachable from p,
// using v to get the list of neighbors at each step.
// No more than maxDepth steps are done.
func Bfs(p Producer, v Visitor, maxDepth uint) ([]Producer, error) {
	visited := make(map[string]bool)
	var result []Producer

	curQueue := []Producer{p}
	nextQueue := []Producer{}
	var curDepth uint
	for len(curQueue) > 0 && curDepth < maxDepth {
		next := curQueue[0]
		curQueue = curQueue[1:]

		if _, found := visited[next.ID()]; found {
			continue
		}

		neighbors, err := v.Visit(next)
		if err != nil {
			return nil, err
		}
		visited[next.ID()] = true

		result = append(result, next)

		nextQueue = append(nextQueue, neighbors...)

		if len(curQueue) == 0 {
			curQueue, nextQueue = nextQueue, curQueue
			curDepth++
		}
	}

	return result, nil
}
