package utils

import (
	"encoding/json"
	"log"
	"os"
)

// Network

type Network struct {
	Bits  int // address length
	Bin   int // buckets 
	Nodes map[int]Node // nodes
}

// create a new network
func (net *Network) NewNetwork(bits int, bin int) {
	net.Bits = bits
	net.Bin = bin
	net.Nodes = make(map[int]Node)
}

// Nodes

type Node struct {
	Net *Network // network
	Id  int    // id
	Adj []int // adjacent nodes
	StorageSet map[int]Node // storage set
	CacheSet map[int]Node // cache set
	CanPay bool // can pay
}



// add a node to the network
func (net *Network) AddNode(node *Node) bool {
	if node.Net != nil {
		return false
	}
	node.Net = net
	net.Nodes[node.Id] = *node
	return true
}



func LoadNodes(path string){
	file, _ := os.Open(path)

	defer file.Close()

	decoder := json.NewDecoder(file)



    log.Println(&decoder)

}

// func LoadNodes(path string) bool {
// 	file, _ := os.Open(path)

// 	defer file.Close()

// 	decoder := json.NewDecoder(file)

// 	var test Test
// 	decoder.Decode(&test)

// 	fmt.Println(test.Bin)
// 	fmt.Println(test.Bits)

// 	for i := 0; i < len(test.Nodes); i++ {
// 		fmt.Println(test.Nodes[i])
// 	}

// 	return true
// }
