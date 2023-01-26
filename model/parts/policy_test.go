package policy

import (
	"fmt"
	. "go-incentive-simulation/model/parts/types"
	"gotest.tools/assert"
	"testing"
)

func TestResponisbleNodes(t *testing.T) {
	nodesId := []int{64132, 49693, 45280, 42779, 41852, 43812, 47987, 43377, 41471}
	chunkAdd := 11
	values := findResponsibleNodes(nodesId, chunkAdd)

	assert.Equal(t, len(values), 4)
}

func TestSendRequest(t *testing.T) {
	path := "../../data/nodes_data_8_10000.txt"
	network := Network{}
	_, _, nodes := network.Load(path)

	var testNodes []*Node

	graph := Graph{}

	for _, v := range nodes {
		testNodes = append(testNodes, v)
		addNode := graph.AddNode(v)
		if addNode != nil {
			panic(addNode)
		}
	}
	edge1 := Edge{testNodes[0].Id, testNodes[1].Id, EdgeAttrs{A2b: 10, Last: 20}}
	addEdg := graph.AddEdge(&edge1)
	if addEdg != nil {
		panic(addEdg)
	}

	state := State{
		Graph:                   &graph,
		Originators:             []int{35506, 61875, 20655, 55142, 33831, 33831, 33831, 33831, 33831},
		NodesId:                 []int{64132, 49693, 45280, 42779, 41852, 43812, 47987, 43377, 41471},
		RouteLists:              []Route{},
		PendingMap:              PendingMap{},
		RerouteMap:              RerouteMap{},
		CacheListMap:            CacheListMap{},
		OriginatorIndex:         1,
		SuccessfulFound:         0,
		FailedRequestsThreshold: 0,
		FailedRequestsAccess:    0,
		TimeStep:                0,
	}
	found, route, thresholdFailed, accessFailed, paymentsList := SendRequest(&state)
	fmt.Println(found, route, thresholdFailed, accessFailed, paymentsList)
}
