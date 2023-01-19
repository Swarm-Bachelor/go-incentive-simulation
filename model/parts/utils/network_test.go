package utils_test

import (
	net "go-incentive-simulation/model/parts/utils"
	"log"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetwork(t *testing.T) {
	// Creates a network
	n := net.Network{}
	n.NewNetwork(8, 16)

	// Asserts the network
	assert.Equal(t, 8, n.Bits)
	assert.Equal(t, 16, n.Bin)
	assert.NotNil(t, n.Nodes)
}

func TestAddNode(t *testing.T) {
	// Creates a network
	n := net.Network{}
	n.NewNetwork(8, 16)

	// Creates a node
	n.Nodes[1] = net.Node{
		Net:        &n,
		Id:         1111,
		Adj:        []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		StorageSet: make(map[int]net.Node),
		CacheSet:   make(map[int]net.Node),
		CanPay:     true,
	}

	// Asserts the network
	assert.Equal(t, 1, len(n.Nodes))
	assert.Equal(t, 1111, n.Nodes[1].Id)
	assert.Equal(t, 16, len(n.Nodes[1].Adj))
	assert.NotNil(t, n.Nodes[1].StorageSet)
	assert.NotNil(t, n.Nodes[1].CacheSet)
	assert.Equal(t, true, n.Nodes[1].CanPay)
}

func TestAddNodesToNetwork(t *testing.T) {
	// Creates network 1
	n1 := net.Network{}
	n1.NewNetwork(8, 16)

	assert.Equal(t, 0, len(n1.Nodes))

	node1 := net.Node{
		Id:         1111,
		Adj:        []int{10, 11, 12, 13, 14, 15, 16, 17},
		StorageSet: make(map[int]net.Node),
		CacheSet:   make(map[int]net.Node),
		CanPay:     true,
	}

	node2 := net.Node{
		Id:         2222,
		Adj:        []int{2, 3, 4, 5, 6, 7, 8},
		StorageSet: make(map[int]net.Node),
		CacheSet:   make(map[int]net.Node),
		CanPay:     false,
	}
	// Adds tow nodes to network 1
	n1.AddNode(&node1)
	n1.AddNode(&node2)

	// Checks if the nodes were added to the network
	assert.Equal(t, 2, len(n1.Nodes))
	// Logs the nodes
	log.Printf("Node %+v\n", n1.Nodes)

	// Creates network 2
	n2 := net.Network{}
	n2.NewNetwork(8, 16)

	assert.Equal(t, 0, len(n2.Nodes))

	node3 := net.Node{
		Id:         1111,
		Adj:        []int{13, 14, 15, 16, 17},
		StorageSet: make(map[int]net.Node),
		CacheSet:   make(map[int]net.Node),
		CanPay:     true,
	}

	node4 := net.Node{
		Id:         2222,
		Adj:        []int{7, 8, 9, 10, 11, 12, 13},
		StorageSet: make(map[int]net.Node),
		CacheSet:   make(map[int]net.Node),
		CanPay:     false,
	}
	// Adds tow nodes to network 2
	n2.AddNode(&node3)
	n2.AddNode(&node4)

	// Checks if the nodes were added to the network
	assert.Equal(t, 2, len(n2.Nodes))
	// Logs the nodes
	log.Printf("Node: %+v\n", n2.Nodes)
}

func TestLoad(t *testing.T) {

	net.LoadNodes("input_test.json")

}
