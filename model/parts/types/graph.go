package types

import (
	"fmt"
	"sync"
)

// Graph structure, node Ids in array and edges in map
type Graph struct {
	Network  *Network
	Nodes    []*Node
	NodeIds  []int
	Edges    map[int][]Edge
	NodesMap map[int]*Node
	Mutex    sync.Mutex
}

// Edge that connects to Nodes with attributes about the connection
type Edge struct {
	FromNodeId int
	ToNodeId   int
	Attrs      EdgeAttrs
}

// EdgeAttrs Edge attributes structure,
// "a2b" show how much this node asked from other node,
// "last" is for the last forgiveness time,
// "threshold" is for the adjustable threshold limit.
type EdgeAttrs struct {
	A2B       int
	Last      int
	Threshold int
}

// AddNode will add a Node to a graph
func (g *Graph) AddNode(node *Node) error {
	if ContainsNode(g.Nodes, node) {
		err := fmt.Errorf("node %d already exists", node.Id)
		return err
	} else {
		g.Nodes = append(g.Nodes, node)
		return nil
	}
}

// AddEdge will add an edge from a node to a node
func (g *Graph) AddEdge(edge Edge) error {
	toNode := g.GetNode(edge.ToNodeId)
	fromNode := g.GetNode(edge.FromNodeId)
	if toNode == nil || fromNode == nil {
		return fmt.Errorf("not a valid edge from %d ---> %d", fromNode.Id, toNode.Id)
	} else if ContainsEdge(g.Edges[fromNode.Id], edge) {
		//if ContainsEdge(g.Edges[edge.FromNodeId], edge) {
		return fmt.Errorf("edge from node %d ---> %d already exists", edge.FromNodeId, edge.ToNodeId)
	} else {
		newEdges := append(g.Edges[edge.FromNodeId], edge)
		g.Edges[edge.FromNodeId] = newEdges
		return nil
	}
}

func (g *Graph) GetEdgeData(fromNodeId int, toNodeId int) EdgeAttrs {
	edgeData := EdgeAttrs{}
	for _, edge := range g.Edges[fromNodeId] {
		if edge.ToNodeId == toNodeId {
			edgeData = edge.Attrs
			break
		}
	}
	return edgeData
}

// GetNode getNode will return a node point if exists or return nil
func (g *Graph) GetNode(nodeId int) *Node {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()
	node, ok := g.NodesMap[nodeId]
	if ok {
		return node
	}
	return nil
}

func ContainsNode(v []*Node, node *Node) bool {
	for _, v := range v {
		if v.Id == node.Id {
			return true
		}
	}
	return false
}

func ContainsEdge(v []Edge, edge Edge) bool {
	for _, v := range v {
		if v.ToNodeId == edge.ToNodeId {
			return true
		}
	}
	return false
}

func (g *Graph) Print() {
	for _, v := range g.Nodes {
		fmt.Printf("%d : ", v.Id)
		for _, i := range v.AdjIds {
			for _, v := range i {
				fmt.Printf("%d ", v)
			}
		}
		fmt.Println()
	}
}
