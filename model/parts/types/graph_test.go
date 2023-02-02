package types

import (
	"testing"
)

func testContainsNode(t *testing.T) {

}

func testAddNode(t *testing.T) {

}

func testGetNode(t *testing.T) {

}

func TestAddEdge(t *testing.T) {
	path := "../../../data/nodes_data_8_10000.txt"
	network := Network{}
	_, _, nodes := network.Load(path)
	var testNodes []*Node
	counter := 0

	for _, node := range nodes {
		if counter == 10 {
			break
		}
		counter++
		testNodes = append(testNodes, node)
	}
	graph := &Graph{Network: &network, Nodes: testNodes, Edges: map[int][]Edge{}, NodesMap: network.Nodes}

	edge1 := Edge{testNodes[0].Id, testNodes[1].Id, EdgeAttrs{10, 20, 16}}
	err := graph.AddEdge(edge1)
	if err != nil {
		t.Error("addEdge function returned an error message: ", err)
	}

	//err = graph.AddEdge(edge1)
	//if err == nil {
	//	t.Error("addedge should have returned an error for trying to add the same edge twice")
	//}

}

func TestGetEdgeData(t *testing.T) {
	// fromNodes := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// toNodes := []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	// graph := Graph{}
}
