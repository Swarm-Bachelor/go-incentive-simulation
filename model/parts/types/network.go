package types

import (
	"encoding/json"
	. "go-incentive-simulation/model/general"
	"os"
	"sort"
)

type Network struct {
	Bits  int
	Bin   int
	Nodes map[int]*Node
}

type Node struct {
	Network    *Network
	Id         int
	AdjIds     [][]int
	storageSet []int
	cacheSet   []int
	canPay     bool
}

type Test struct {
	Bits  int `json:"bits"`
	Bin   int `json:"bin"`
	Nodes []struct {
		ID  int   `json:"id"`
		Adj []int `json:"adj"`
	} `json:"Nodes"`
}

func (network *Network) Load(path string) (int, int, map[int]*Node) {
	file, _ := os.Open(path)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	decoder := json.NewDecoder(file)

	var test Test
	err := decoder.Decode(&test)
	if err != nil {
		return 0, 0, nil
	}

	network.Bits = test.Bits
	network.Bin = test.Bin
	network.Nodes = make(map[int]*Node)

	for _, node := range test.Nodes {
		node1 := network.node(node.ID)
		sort.Ints(node.Adj)
		for _, adj := range node.Adj {
			node2 := network.node(adj)
			node1.add(node2)
		}
	}

	return network.Bits, network.Bin, network.Nodes
}

func (network *Network) node(value int) *Node {
	if value < 0 || value >= (1<<network.Bits) {
		panic("address out of range")
	}
	res := Node{
		Network:    network,
		Id:         value,
		AdjIds:     make([][]int, network.Bits),
		storageSet: []int{0},
		cacheSet:   []int{0},
		canPay:     true,
	}
	if _, ok := network.Nodes[value]; !ok {
		network.Nodes[value] = &res
		return &res
	}
	return network.Nodes[value]

}

func (network *Network) generate(count int) {
	// TODO: implement
}

func (node *Node) add(other *Node) bool {
	if node.Network == nil || node.Network != other.Network || node == other {
		return false
	}
	if node.AdjIds == nil {
		node.AdjIds = make([][]int, node.Network.Bits)
	}
	if other.AdjIds == nil {
		other.AdjIds = make([][]int, other.Network.Bits)
	}
	bit := node.Network.Bits - BitLength(node.Id^other.Id)
	if bit < 0 || bit >= node.Network.Bits {
		return false
	}
	isDup := Contains(node.AdjIds[bit], other.Id) || Contains(other.AdjIds[bit], node.Id)
	if len(node.AdjIds[bit]) < node.Network.Bin && len(other.AdjIds[bit]) < node.Network.Bin && !isDup {
		node.AdjIds[bit] = append(node.AdjIds[bit], other.Id)
		other.AdjIds[bit] = append(other.AdjIds[bit], node.Id)
		return true
	}
	return false
}

func (node *Node) IsNil() bool {
	return node.Id == 0
}
