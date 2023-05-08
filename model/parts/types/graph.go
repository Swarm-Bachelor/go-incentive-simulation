package types

import (
	"fmt"
	"go-incentive-simulation/config"
	"sort"
	"sync"
)

// Graph structure, node Ids in array and edges in map
type Graph struct {
	*Network
	CurState  State
	Nodes     []*Node
	NodeIds   []NodeId
	Edges     map[NodeId]map[NodeId]*Edge
	RespNodes map[ChunkId][4]NodeId
	Mutex     sync.Mutex
}

// Edge that connects to NodesMap with attributes about the connection
type Edge struct {
	FromNodeId NodeId
	ToNodeId   NodeId
	Attrs      EdgeAttrs
	Mutex      *sync.Mutex
}

// EdgeAttrs Edge attributes structure,
// "a2b" show how much this node asked from other node,
// "lastEpoch" is the epoch where it was last forgiven.
// "threshold" is for the adjustable threshold limit.
type EdgeAttrs struct {
	A2B       int
	LastEpoch int
	Threshold int
}

//func (g *Graph) FindResponsibleNodes(chunkId int) [4]int {
//	return g.RespNodes[chunkId]
//}

func BinarySearchClosest(arr []NodeId, target int, n int) []NodeId {
	left, right := 0, len(arr)-1
	mid := 0
	for left <= right {
		mid = (left + right) / 2
		curNodeId := arr[mid].ToInt()
		if curNodeId == target {
			//if curNodeId > target-(n/2) && curNodeId < target+(n/2) {
			break
		} else if curNodeId < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	left = mid - n
	if left < 0 {
		left = 0
	}
	right = mid + n
	if right > len(arr) {
		right = len(arr)
	}
	return arr[left:right]
}

func (g *Graph) FindResponsibleNodes(chunkId ChunkId) [4]NodeId {
	chunkIdInt := chunkId.ToInt()
	if config.IsPrecomputeRespNodes() {
		return g.RespNodes[chunkId]

	} else {

		if _, ok := g.RespNodes[chunkId]; ok {
			return g.RespNodes[chunkId]

		} else {
			numNodesSearch := config.GetBits()
			closestNodes := BinarySearchClosest(g.NodeIds, chunkIdInt, numNodesSearch)

			distances := make([]int, len(closestNodes))
			result := [4]NodeId{}

			for i, nodeId := range closestNodes {
				distances[i] = nodeId.ToInt() ^ chunkIdInt
			}

			sort.Slice(distances, func(i, j int) bool { return distances[i] < distances[j] })

			for i := 0; i < 4; i++ {
				result[i] = NodeId(distances[i] ^ chunkIdInt) // this results in the nodeId again
			}
			g.RespNodes[chunkId] = result

			return result
		}
	}
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

func (g *Graph) GetNodeAdj(nodeId NodeId) [][]NodeId {
	return g.GetNode(nodeId).AdjIds
}

// AddEdge will add an edge from a node to a node
func (g *Graph) AddEdge(fromNodeId NodeId, toNodeId NodeId, attrs EdgeAttrs) error {
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
		mutex := &sync.Mutex{}
		if g.EdgeExists(toNodeId, fromNodeId) {
			mutex = g.GetEdge(toNodeId, fromNodeId).Mutex
		}
		newEdge := &Edge{FromNodeId: fromNodeId, ToNodeId: toNodeId, Attrs: attrs, Mutex: mutex}
		g.Edges[fromNodeId][toNodeId] = newEdge
		return nil
	}
}

func (g *Graph) LockEdge(nodeA NodeId, nodeB NodeId) {
	edge := g.GetEdge(nodeA, nodeB)
	edge.Mutex.Lock()
}

func (g *Graph) UnlockEdge(nodeA NodeId, nodeB NodeId) {
	edge := g.GetEdge(nodeA, nodeB)
	edge.Mutex.Unlock()
}

func (g *Graph) GetEdge(fromNodeId NodeId, toNodeId NodeId) *Edge {
	if g.EdgeExists(fromNodeId, toNodeId) {
		return g.Edges[fromNodeId][toNodeId]
	}
	return &Edge{}
}

func (g *Graph) GetEdgeData(fromNodeId NodeId, toNodeId NodeId) EdgeAttrs {
	if g.EdgeExists(fromNodeId, toNodeId) {
		return g.GetEdge(fromNodeId, toNodeId).Attrs
	}
	return EdgeAttrs{}
}

func (g *Graph) EdgeExists(fromNodeId NodeId, toNodeId NodeId) bool {
	if _, ok := g.Edges[fromNodeId][toNodeId]; ok {
		return true
	}
	return false
}

func (g *Graph) SetEdgeData(fromNodeId NodeId, toNodeId NodeId, edgeAttrs EdgeAttrs) bool {
	if g.EdgeExists(fromNodeId, toNodeId) {
		g.Edges[fromNodeId][toNodeId].Attrs = edgeAttrs
		return true
	}
	return false
}

// GetNode getNode will return a node point if exists or return nil
func (g *Graph) GetNode(nodeId NodeId) *Node {
	//g.Mutex.Lock()
	//defer g.Mutex.Unlock()
	node, ok := g.NodesMap[nodeId]
	if ok {
		return node
	}
	return nil
}

func ContainsNode(Nodes []*Node, node *Node) bool {
	for _, curNode := range Nodes {
		if curNode.Id == node.Id {
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
