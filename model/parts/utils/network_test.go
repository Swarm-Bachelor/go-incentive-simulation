package utils_test

import (
	net "go-incentive-simulation/model/parts/utils"
	"log"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetwork(t *testing.T) {
	n := net.Network{}
	n.NewNetwork(8, 16)
	assert.Equal(t, 8, n.Bits)
	assert.Equal(t, 16, n.Bin)
	assert.NotNil(t, n.Nodes)
}

func TestAddNode(t *testing.T) {
	n := net.Network{}
	n.NewNetwork(8, 16)

	n.Nodes[1] = net.Node{
		Net: &n,
		Id: 1111,
		Adj: []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		StorageSet: make(map[int]net.Node),
		CacheSet: make(map[int]net.Node),
		CanPay: true,
	}

	assert.Equal(t, 1, len(n.Nodes))
	assert.Equal(t, 1111, n.Nodes[1].Id)
	assert.Equal(t, 16, len(n.Nodes[1].Adj))
	assert.NotNil(t, n.Nodes[1].StorageSet)
	assert.NotNil(t, n.Nodes[1].CacheSet)
	assert.Equal(t, true, n.Nodes[1].CanPay)
}

func TestAddNodesToNetwork(t *testing.T) {
	n1 := net.Network{}
	n1.NewNetwork(8, 16)

	assert.Equal(t, 0, len(n1.Nodes))

	node1 := net.Node{
		Id: 1111,
		Adj: []int{10, 11, 12, 13, 14, 15, 16, 17},
		StorageSet: make(map[int]net.Node),
		CacheSet: make(map[int]net.Node),
		CanPay: true,
	}

	node2 := net.Node{
		Id: 2222,
		Adj: []int{2, 3, 4, 5, 6, 7, 8,},
		StorageSet: make(map[int]net.Node),
		CacheSet: make(map[int]net.Node),
		CanPay: false,
	}
	n1.AddNode(&node1)
	n1.AddNode(&node2)

	assert.Equal(t, 2, len(n1.Nodes))
	log.Printf("Node %+v\n", n1.Nodes)

	n2 := net.Network{}
	n2.NewNetwork(8, 16)

	assert.Equal(t, 0, len(n2.Nodes))

	node3 := net.Node{
		Id: 1111,
		Adj: []int{13, 14, 15, 16, 17},
		StorageSet: make(map[int]net.Node),
		CacheSet: make(map[int]net.Node),
		CanPay: true,
	}

	node4 := net.Node{
		Id: 2222,
		Adj: []int{7, 8, 9, 10, 11, 12, 13},
		StorageSet: make(map[int]net.Node),
		CacheSet: make(map[int]net.Node),
		CanPay: false,
	}
	n2.AddNode(&node3)
	n2.AddNode(&node4)

	assert.Equal(t, 2, len(n2.Nodes))
	log.Printf("Node: %+v\n", n2.Nodes)

}