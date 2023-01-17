package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (net *Network) LoadNodes(path string) {
	fileContent, _ := os.Open(path)
	defer fileContent.Close()
	
	byteResult, _ := ioutil.ReadAll(fileContent)

	var res map[string]interface{}
	json.Unmarshal([]byte(byteResult), &res)
 
	fmt.Println(res["users"])
}




// func (node *Node) add(other *Node) bool {
// 	net := &node.Net
// 	if (net == nil) || &other.Net != net || node == other {
// 		return false
// 	}
// 	bits := net.Bits - (node.Id ^ other.Id)
// }