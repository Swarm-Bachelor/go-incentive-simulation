package types

import (
	"fmt"
	"sync"
)

// Graph structure, node Ids in array and edges in map
type Graph struct {
	Network      *Network
	Nodes        []*Node
	NodeIds      []int
	Edges        map[int]map[int]Edge
	NodesMap     map[int]*Node
	RespNodes    [][4]int
	GetNodeMutex sync.Mutex
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

func (g *Graph) FindResponsibleNodes(chunkId int) [4]int {
	return g.RespNodes[chunkId]
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
func (g *Graph) AddEdge(fromNodeId int, toNodeId int, attrs EdgeAttrs) error {
	toNode := g.GetNode(toNodeId)
	fromNode := g.GetNode(fromNodeId)
	if toNode == nil || fromNode == nil {
		return fmt.Errorf("not a valid edge from %d ---> %d", fromNode.Id, toNode.Id)
	} else if g.EdgeExists(fromNodeId, toNodeId) {
		//if ContainsEdge(g.Edges[edge.FromNodeId], edge) {
		return fmt.Errorf("edge from node %d ---> %d already exists", fromNodeId, toNodeId)
	} else {
		//newEdges := append(g.Edges[edge.FromNodeId], edge)
		//g.Edges[edge.FromNodeId] = newEdges
		newEdge := Edge{FromNodeId: fromNodeId, ToNodeId: toNodeId, Attrs: attrs}
		g.Edges[fromNodeId][toNodeId] = newEdge
		return nil
	}
}

func (g *Graph) GetEdge(fromNodeId int, toNodeId int) Edge {
	if g.EdgeExists(fromNodeId, toNodeId) {
		return g.Edges[fromNodeId][toNodeId]
	}
	return Edge{}
}

func (g *Graph) GetEdgeData(fromNodeId int, toNodeId int) EdgeAttrs {
	if g.EdgeExists(fromNodeId, toNodeId) {
		return g.GetEdge(fromNodeId, toNodeId).Attrs
	}
	return EdgeAttrs{}
}

func (g *Graph) EdgeExists(fromNodeId int, toNodeId int) bool {
	if _, ok := g.Edges[fromNodeId][toNodeId]; ok {
		return true
	}
	return false
}

func (g *Graph) SetEdgeData(fromNodeId int, toNodeId int, edgeAttrs EdgeAttrs) bool {
	if g.EdgeExists(fromNodeId, toNodeId) {
		newEdge := Edge{FromNodeId: fromNodeId, ToNodeId: toNodeId, Attrs: edgeAttrs}
		g.Edges[fromNodeId][toNodeId] = newEdge
		return true
	}
	return false
}

// GetNode getNode will return a node point if exists or return nil
func (g *Graph) GetNode(nodeId int) *Node {
	g.GetNodeMutex.Lock()
	defer g.GetNodeMutex.Unlock()
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
