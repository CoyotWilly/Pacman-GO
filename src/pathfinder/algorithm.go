package pathfinder

import "container/heap"

// astar is an A* pathfinding implementation.

// Pattern is an interface which allows A* searching on arbitrary objects which
// can represent a weighted graph.
type Pattern interface {
	// PathNeighbors returns the direct neighboring nodes of this node which
	// can be pathed to.
	PathNeighbors() []Pattern
	// PathNeighborCost calculates the exact movement cost to neighbor nodes.
	PathNeighborCost(to Pattern) float64
	// PathEstimatedCost is a heuristic method for estimating movement costs
	// between non-adjacent nodes.
	PathEstimatedCost(to Pattern) float64
}

// node is a wrapper to store A* data for a Pattern node.
type node struct {
	pattern Pattern
	cost    float64
	rank    float64
	parent  *node
	open    bool
	closed  bool
	index   int
}

// nodeMap is a collection of nodes keyed by Pattern nodes for quick reference.
type nodeMap map[Pattern]*node

// get gets the Pattern object wrapped in a node, instantiating if required.
func (nm nodeMap) get(p Pattern) *node {
	n, ok := nm[p]
	if !ok {
		n = &node{
			pattern: p,
		}
		nm[p] = n
	}
	return n
}

// Path calculates a short path and the distance between the two Pattern nodes.
//
// If no path is found, found will be false.
func Path(from, to Pattern) (path []Pattern, distance float64, found bool) {
	nm := nodeMap{}
	nq := &priorityQueue{}
	heap.Init(nq)
	fromNode := nm.get(from)
	fromNode.open = true
	heap.Push(nq, fromNode)
	for {
		if nq.Len() == 0 {
			// There's no path, return found false.
			return
		}
		current := heap.Pop(nq).(*node)
		current.open = false
		current.closed = true

		if current == nm.get(to) {
			// Found a path to the goal.
			var p []Pattern
			curr := current
			for curr != nil {
				p = append(p, curr.pattern)
				curr = curr.parent
			}
			return p, current.cost, true
		}

		for _, neighbor := range current.pattern.PathNeighbors() {
			cost := current.cost + current.pattern.PathNeighborCost(neighbor)
			neighborNode := nm.get(neighbor)
			if cost < neighborNode.cost {
				if neighborNode.open {
					heap.Remove(nq, neighborNode.index)
				}
				neighborNode.open = false
				neighborNode.closed = false
			}
			if !neighborNode.open && !neighborNode.closed {
				neighborNode.cost = cost
				neighborNode.open = true
				neighborNode.rank = cost + neighbor.PathEstimatedCost(to)
				neighborNode.parent = current
				heap.Push(nq, neighborNode)
			}
		}
	}
}
