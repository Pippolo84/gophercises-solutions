package graph

import (
	"testing"
)

type node struct {
	id        string
	neighbors []node
}

func (n node) ID() string {
	return n.id
}

func (n node) Value() interface{} {
	return n.ID()
}

func nodeVisitor(n Producer) ([]Producer, error) {
	var next []Producer

	for _, p := range n.(node).neighbors {
		next = append(next, p)
	}
	return next, nil
}

func TestBfsWithZeroDepth(t *testing.T) {
	n := node{
		id:        "0",
		neighbors: []node{},
	}
	res, err := Bfs(n, VisitorFunc(nodeVisitor), 0)

	if err != nil {
		t.Errorf("Expecting no errors, but got one from Bfs: %v", err)
	}

	if len(res) != 0 {
		t.Errorf("Expecting empty slice, got %v", res)
	}
}

func TestBfsWithDepthOne(t *testing.T) {
	n2 := node{
		id:        "1",
		neighbors: []node{},
	}
	n1 := node{
		id:        "0",
		neighbors: []node{n2},
	}
	res, err := Bfs(n1, VisitorFunc(nodeVisitor), 1)

	if err != nil {
		t.Errorf("Expecting no errors, but got one from Bfs: %v", err)
	}

	if len(res) != 1 {
		t.Errorf("Expecting slice of lenght one, got %v", res)
	}

	if res[0].ID() != n1.id {
		t.Errorf("Expected node with ID %v, got ID %v", n1.id, res[0].ID())
	}
}

func TestBfsWithDepthGreaterThanOne(t *testing.T) {
	n4 := node{
		id:        "3",
		neighbors: []node{},
	}
	n3 := node{
		id:        "2",
		neighbors: []node{n4},
	}
	n2 := node{
		id:        "1",
		neighbors: []node{},
	}
	n1 := node{
		id:        "0",
		neighbors: []node{n2, n3},
	}
	res, err := Bfs(n1, VisitorFunc(nodeVisitor), 4)

	if err != nil {
		t.Errorf("Expecting no errors, but got one from Bfs: %v", err)
	}

	if len(res) != 4 {
		t.Errorf("Expecting slice of lenght four, got %v", res)
	}

	if res[0].ID() != n1.id || res[1].ID() != n2.ID() || res[2].ID() != n3.ID() || res[3].ID() != n4.ID() {
		t.Errorf("Expected nodes with IDs %v %v %v %v, got IDs %v %v %v %v",
			n1.id, n2.id, n3.id, n4.id, res[0].ID(), res[1].ID(), res[2].ID(), res[3].ID())
	}
}

func TestBfsWithCycle(t *testing.T) {
	n2 := node{
		id:        "1",
		neighbors: []node{},
	}
	n1 := node{
		id:        "0",
		neighbors: []node{n2},
	}

	n2.neighbors = append(n2.neighbors, n1)

	res, err := Bfs(n1, VisitorFunc(nodeVisitor), ^uint(0))

	if err != nil {
		t.Errorf("Expecting no errors, but got one from Bfs: %v", err)
	}

	if len(res) != 2 {
		t.Errorf("Expecting slice of lenght two, got %v", res)
	}

	if res[0].ID() != n1.id || res[1].ID() != n2.ID() {
		t.Errorf("Expected nodes with IDs %v %v, got IDs %v %v", n1.id, n2.id, res[0].ID(), res[1].ID())
	}
}
