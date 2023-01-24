package utils

import (
	"fmt"
	ct "go-incentive-simulation/model/constants"
	"math/rand"
	"time"
)

func CreateGraphNetwork(filename string) (*Graph, error) {
	fmt.Println("Creating graph network...")
	graph := &Graph{
		edges: make(map[int][]*Edge),
	}
	net := new(Network)
	_, _, nodes := net.load(filename)
	for _, node := range nodes {
		err := graph.AddNode(node)
		if err != nil {
			return nil, err
		}
	}
	for _, node := range graph.Nodes() {
		nodeAdj := node.adj
		for _, adjItems := range nodeAdj {
			for _, item := range adjItems {
				// "a2b" show how much this node asked from other node,
				// "last" is for the last forgiveness time
				attrs := EdgeAttrs{a2b: 0, last: 0}
				edge := Edge{fromNodeId: node.id, toNodeId: item.id, attrs: attrs}
				err := graph.AddEdge(&edge)
				if err != nil {
					return nil, err
				}
				// graph.SetEdgeAttributes()
			}
		}
	}
	fmt.Println("Graph network is created.")
	return graph, nil
}

func isThresholdFailed(firstNodeId int, secondNodeId int, chunkId int, g *Graph) bool {
	if ct.Constants.GetThresholdEnabled() {
		edgeDataFirst := g.GetEdgeData(firstNodeId, secondNodeId)
		p2pFirst := edgeDataFirst.a2b

		edgeDataSecond := g.GetEdgeData(secondNodeId, firstNodeId)
		p2pSecond := edgeDataSecond.a2b

		price := p2pFirst - p2pSecond + peerPriceChunk(secondNodeId, chunkId)
		fmt.Printf("price: %d", price)
		return price > ct.Constants.GetThreshold()
	}
	return false
}

type Payment struct {
	firstNodeId  int
	payNextId    int
	chunkId      int
	isOriginator bool
}

type Threshold [2]*Node

type RerouteMap map[int][]*Node

func getNext(firstNode *Node, chunkId int, graph *Graph, mainOriginator *Node, prevNodePaid bool, rerouteMap RerouteMap) (int, *Node, []Threshold, bool, bool, Payment, bool) {
	var nextNode *Node = nil
	var payNext *Node = nil
	var thresholdList []Threshold
	var thresholdFailed bool
	var accessFailed bool
	var payment Payment
	resultInt := 1
	lastDistance := firstNode.id ^ chunkId
	fmt.Printf("last distance is : %d, chunk is: %d, first is: %d", lastDistance, chunkId, firstNode.id)
	fmt.Printf("which bucket: %d", 16-BitLength(chunkId^firstNode.id))

	currDist := lastDistance
	payDist := lastDistance
	for _, adj := range firstNode.adj {
		fmt.Println("adj: ", adj)
		for _, node := range adj {
			dist := node.id ^ chunkId
			if BitLength(dist) >= BitLength(lastDistance) {
				continue
			}

			if !isThresholdFailed(firstNode.id, node.id, chunkId, graph) {
				thresholdFailed = false

				// Could probably clean this one up, but keeping it close to original for now
				if dist < currDist {
					if ct.Constants.IsRetryWithAnotherPeer() {
						_, ok := rerouteMap[mainOriginator.id]
						if ok {
							allExceptLast := len(rerouteMap[mainOriginator.id]) - 1
							if containsNode(rerouteMap[mainOriginator.id][:allExceptLast], node) {
								continue
							} else {
								currDist = dist
								nextNode = node
							}
						} else {
							currDist = dist
							nextNode = node
						}
					} else {
						currDist = dist
						nextNode = node
					}
				}
			} else {
				thresholdFailed = true
				if ct.Constants.GetPaymentEnabled() {
					if dist < payDist {
						payDist = dist
						payNext = node
					}
				}
				listItem := Threshold{firstNode, node}
				thresholdList = append(thresholdList, listItem)
			}
		}
	}
	if nextNode != nil {
		thresholdFailed = false
		accessFailed = false
	} else {
		if !thresholdFailed {
			accessFailed = true
			resultInt = -2
			nextNode = nil
			// nextNode = -2 // accessFailed, TYPE MISMATCH ??
		} else {
			resultInt = -1
			nextNode = nil
			// nextNode = -1 // thresholdFailed, TYPE MISMATCH ??
		}
		if ct.Constants.GetPaymentEnabled() {
			if payNext != nil {
				accessFailed = false
				if ct.Constants.IsOnlyOriginatorPays() {
					if firstNode == mainOriginator {
						payment.isOriginator = true
						payment.firstNodeId = firstNode.id
						payment.payNextId = payNext.id
						payment.chunkId = chunkId
						nextNode = payNext
					} else {
						thresholdFailed = true
						resultInt = -1
						nextNode = nil
						// nextNode = -1 TYPE MISMATCH ??
					}
				} else if ct.Constants.IsPayIfOrigPays() {
					if prevNodePaid {
						nextNode = payNext
						thresholdFailed = false
						if firstNode == mainOriginator {
							payment.isOriginator = true
						} else {
							payment.isOriginator = false
						}
						payment.firstNodeId = firstNode.id
						payment.payNextId = payNext.id
						payment.chunkId = chunkId
					} else {
						if firstNode == mainOriginator {
							payment.isOriginator = true
							payment.firstNodeId = firstNode.id
							payment.payNextId = payNext.id
							payment.chunkId = chunkId
							nextNode = payNext
						} else {
							thresholdFailed = true
							resultInt = -1
							nextNode = nil
							// nextNode = -1 // TYPE MISMATCH ??
							payNext = nil
						}
					}
				} else {
					nextNode = payNext
					thresholdFailed = false
					if firstNode == mainOriginator {
						payment.isOriginator = true
					} else {
						payment.isOriginator = false
					}
					payment.firstNodeId = firstNode.id
					payment.payNextId = payNext.id
					payment.chunkId = chunkId
				}
			}
		}
	}
	// TODO: Usikker på dette
	if ct.Constants.GetPaymentEnabled() {
	out:
		for i, item := range thresholdList {
			for _, node := range item {
				if node == payNext {
					if ct.Constants.IsPayIfOrigPays() {
						if firstNode == mainOriginator {
							thresholdList = append(thresholdList[:i], thresholdList[i+1:]...)
						}
					} else {
						thresholdList = append(thresholdList[:i], thresholdList[i+1:]...)
					}
					break out
				}
			}
		}
	}
	if payment == (Payment{}) {
		prevNodePaid = true
	} else {
		prevNodePaid = false
	}
	fmt.Printf("next node is: %d", nextNode.id)
	return resultInt, nextNode, thresholdList, thresholdFailed, accessFailed, payment, prevNodePaid
}

type Request struct {
	originator *Node
	chunkId    int
}

type CacheListMap map[*Node][]map[int]int

type Route []int

// ConsumeTask cacheDict is map of nodes containing an array of maps with key as a chunkAddr and a popularity counter
func ConsumeTask(request *Request, graph *Graph, respNodes []*Node, rerouteMap RerouteMap, cacheListMap CacheListMap) (bool, Route, [][]Threshold, bool, []Payment) {
	var thresholdFailedList [][]Threshold
	var paymentList []Payment
	originator := request.originator
	chunkId := request.chunkId
	mainOriginator := originator
	found := false
	foundByCaching := false
	route := Route{originator.id}
	var resultInt int
	var nextNode *Node
	var thresholdList []Threshold
	// thresholdFailed := false
	var accessFailed bool
	var payment Payment
	var prevNodePaid bool
	if ct.Constants.IsPayIfOrigPays() {
		prevNodePaid = true
	}
	if containsNode(respNodes, originator) {
		// originator has the chunk
		found = true
	} else {
	out:
		for _, node := range respNodes {
			// fmt.Printf("orig: %d, chunk_id: %d", originator.id, chunkId)
			if node != originator {
				// nextNode, thresholdList, thresholdFailed, accessFailed, payment, prevNodePaid = getNext(originator, chunkId, graph, mainOriginator, prevNodePaid, rerouteMap)
				resultInt, nextNode, thresholdList, _, accessFailed, payment, prevNodePaid = getNext(originator, chunkId, graph, mainOriginator, prevNodePaid, rerouteMap)
				if payment != (Payment{}) {
					paymentList = append(paymentList, payment)
				}
				if len(thresholdList) > 0 {
					thresholdFailedList = append(thresholdFailedList, thresholdList)
				}
				route = append(route, nextNode.id)
				// if not isinstance(next_node, int), originale versjonen
				if !(resultInt <= -1) && nextNode != nil {
					if containsNode(respNodes, nextNode) {
						fmt.Println("is not in cache")
						found = true
						break out
					}
					cacheList, ok := cacheListMap[nextNode]
					if ok {
						for _, cacheMap := range cacheList {
							_, ok2 := cacheMap[chunkId]
							if ok2 {
								fmt.Println("is in cache")
								found = true
								foundByCaching = true
								break out
							}
						}
					}
					originator = nextNode
				} else {
					break out
				}
			}
		}
	}
	route = append(route, chunkId)

	if ct.Constants.IsForwarderPayForceOriginatorToPay() {
		// if !Contains(route, -2) { // Gir ikke mening lengre
		if resultInt != -2 {
			// NOT accessFailed
			if len(paymentList) > 0 {
				firstPayment := paymentList[0]
				if !firstPayment.isOriginator {
					// TODO: Dobbelsjekk at logikken under her matcher originalen
					for i := range route[:len(route)-1] {
						p := Payment{route[i], route[i+1], route[len(route)-1], false}
						for j, tmp := range paymentList {
							if p.payNextId == tmp.payNextId && p.firstNodeId == tmp.firstNodeId && p.chunkId == tmp.chunkId {
								break
							}
							if j == len(paymentList) {
								// payment is now definitely not in paymentList
								if i == 0 {
									p.isOriginator = true
								}
								if i != len(route)-2 {
									paymentList = append(paymentList[:i+1], paymentList[i:]...)
									paymentList[i] = p
								} else {
									continue
								}
							}
						}
					}
				} else {
					// TODO: Dobbelsjekk at logikken under her matcher originalen
					for i := range route[1 : len(route)-1] {
						p := Payment{route[i], route[i+1], route[len(route)-1], false}
						for j, tmp := range paymentList {
							if p.payNextId == tmp.payNextId && payment.firstNodeId == tmp.firstNodeId && p.chunkId == tmp.chunkId {
								break
							}
							if j == len(paymentList) {
								// payment is now definitely not in paymentList
								if i != len(route)-2 {
									paymentList = append(paymentList[:i+1], paymentList[i:]...)
									paymentList[i] = p
								} else {
									continue
								}
							}
						}
					}
				}
			}
		} else {
			paymentList = []Payment{}
		}
	}
	if foundByCaching {
		// route = append(route, "C") // TYPE MISMATCH
		route = append(route, -3) // TODO: midlertidig fix?
	}
	return found, route, thresholdFailedList, accessFailed, paymentList
}

// TODO: Not used in original
//func getBin(src int, dest int, index int) int {
//	distance := src ^ dest
//	result := index
//	for distance > 0 {
//		distance >>= 1
//		result -= 1
//	}
//	return result
//}

// TODO: Not used in original
//func whichPowerTwo(rangeAddress int) int {
//	return BitLength(rangeAddress) - 1
//}

func getProximityChunk(firstNodeId int, chunkId int) int {
	retVal := ct.Constants.GetBits() - BitLength(firstNodeId^chunkId)
	if retVal <= ct.Constants.GetMaxProximityOrder() {
		return retVal
	} else {
		return ct.Constants.GetMaxProximityOrder()
	}
}

func peerPriceChunk(firstNodeId int, chunkId int) int {
	return (ct.Constants.GetMaxProximityOrder() - getProximityChunk(firstNodeId, chunkId) + 1) * ct.Constants.GetPrice()
}

func choice(nodes []int, k int) []int {
	res := make([]int, 0, k)

	rand.Seed(time.Now().UnixMicro())

	for i := 0; i < k; i++ {
		res = append(res, nodes[rand.Intn(len(nodes))])
	}
	return res
}

// TODO: Not used in original
//func MakeFiles() []int {
//	fmt.Println("Making files...")
//	var filesList []int
//
//	for i := 0; i <= ct.Constants.GetOriginators(); i++ {
//		// TODO: fix this, GetChunks should be a list?
//		// chunksList := choice(ct.Constants.GetChunks(), ct.Constants.GetRangeAddress())
//		// filesList = append(chunksList)
//		fmt.Println(i)
//	}
//	// Gets all constants
//	consts := ct.Constants
//
//	for i := 0; i <= consts.GetOriginators(); i++ {
//		chunksList := rand.Perm(consts.GetChunks())
//		filesList = append(chunksList)
//	}
//	fmt.Println("Files made!")
//	return filesList
//}

func (net *Network) CreateDowloadersList() []int {
	fmt.Println("Creating downloaders list...")

	nodesValue := make([]int, 0, len(net.nodes))
	for i := range net.nodes {
		nodesValue = append(nodesValue, net.nodes[i].id)
	}
	downloadersList := choice(nodesValue, ct.Constants.GetOriginators())

	fmt.Println("Downloaders list create...!")
	return downloadersList
}

// TODO: Not used in original
//func (net *Network) PushSync(fileName string, files []string) {
//	fmt.Println("Pushing sync...")
//	if net == nil {
//		fmt.Println("Network is nil!")
//		return
//	}
//	nodes := net.nodes
//	for i := range nodes {
//		fmt.Println(nodes[i].id)
//	}
//
//	fmt.Println("Pushing sync finished...")
//}
