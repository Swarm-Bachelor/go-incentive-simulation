package utils

import (
	. "go-incentive-simulation/model/constants"
	. "go-incentive-simulation/model/general"
	. "go-incentive-simulation/model/parts/types"
	"sort"
)

func PrecomputeClosestNodes(nodesId []int) [][4]int {
	numPossibleChunks := Constants.GetRangeAddress()
	result := make([][4]int, numPossibleChunks)
	numNodesSearch := Constants.GetBits()

	for chunkId := 0; chunkId < numPossibleChunks; chunkId++ {

		closestNodes := BinarySearchClosest(nodesId, chunkId, numNodesSearch)
		distances := make([]int, len(closestNodes))

		for i, nodeId := range closestNodes {
			distances[i] = nodeId ^ chunkId
		}

		sort.Slice(distances, func(i, j int) bool { return distances[i] < distances[j] })

		for i := 0; i < 4; i++ {
			result[chunkId][i] = distances[i] ^ chunkId // this results in the nodeId again
		}
	}
	return result
}

func SortedKeys(m map[int]*Node) []int {
	keys := make([]int, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func CreateGraphNetwork(net *Network) (*Graph, error) {
	//fmt.Println("Creating graph network...")
	sortedNodeIds := SortedKeys(net.Nodes)
	numNodes := len(net.Nodes)
	Edges := make(map[int]map[int]Edge)
	//for _, nodeId := range sortedNodeIds {
	//	Edges[nodeId] = make(map[int]Edge)
	//}
	respNodes := PrecomputeClosestNodes(sortedNodeIds)
	graph := &Graph{
		Network:   net,
		Nodes:     make([]*Node, 0, numNodes),
		Edges:     Edges,
		NodeIds:   sortedNodeIds,
		NodesMap:  net.Nodes,
		RespNodes: respNodes,
	}

	for _, nodeId := range sortedNodeIds {
		graph.Edges[nodeId] = make(map[int]Edge)

		node := net.Nodes[nodeId]
		err1 := graph.AddNode(node)
		if err1 != nil {
			return nil, err1
		}

		nodeAdj := node.AdjIds
		for _, adjItems := range nodeAdj {
			for _, item := range adjItems {
				threshold := BitLength(nodeId ^ item)
				attrs := EdgeAttrs{A2B: 0, Last: 0, Threshold: threshold}
				err := graph.AddEdge(node.Id, item, attrs)
				if err != nil {
					return nil, err
				}
				// graph.SetEdgeAttributes()
			}
		}
	}

	//fmt.Println("Graph network is created.")
	return graph, nil
}

func isThresholdFailed(firstNodeId int, secondNodeId int, chunkId int, g *Graph) bool {
	if Constants.GetThresholdEnabled() {
		edgeDataFirst := g.GetEdgeData(firstNodeId, secondNodeId)
		p2pFirst := edgeDataFirst.A2B
		edgeDataSecond := g.GetEdgeData(secondNodeId, firstNodeId)
		p2pSecond := edgeDataSecond.A2B

		threshold := Constants.GetThreshold()
		if Constants.IsAdjustableThreshold() {
			threshold = edgeDataFirst.Threshold
		}

		peerPriceChunk := PeerPriceChunk(secondNodeId, chunkId)
		price := p2pFirst - p2pSecond + peerPriceChunk
		// fmt.Printf("price: %d = p2pFirst: %d - p2pSecond: %d + PeerPriceChunk: %d \n", price, p2pFirst, p2pSecond, peerPriceChunk)
		return price > threshold
	}
	return false
}

func getNext(firstNodeId int, chunkId int, graph *Graph, mainOriginatorId int, prevNodePaid bool, rerouteMap RerouteMap) (int, []Threshold, bool, bool, Payment, bool) {
	var nextNodeId int
	var payNextId int
	var thresholdList []Threshold
	var thresholdFailed bool
	var accessFailed bool
	var payment Payment
	lastDistance := firstNodeId ^ chunkId
	//fmt.Printf("\n last distance is : %d, chunk is: %d, first is: %d", lastDistance, chunkId, firstNodeId)
	//fmt.Printf("\n which bucket: %d \n", 16-BitLength(chunkId^firstNodeId))

	currDist := lastDistance
	payDist := lastDistance

	//firstNode := graph.NodesMap[firstNodeId]
	bin := Constants.GetBits() - BitLength(firstNodeId^chunkId)
	firstNodeAdjIds := graph.GetNodeAdj(firstNodeId)

	for _, nodeId := range firstNodeAdjIds[bin] {
		dist := nodeId ^ chunkId
		if BitLength(dist) >= BitLength(lastDistance) {
			continue
		}
		if !isThresholdFailed(firstNodeId, nodeId, chunkId, graph) {
			thresholdFailed = false
			// Could probably clean this one up, but keeping it close to original for now
			if dist < currDist {
				if Constants.IsRetryWithAnotherPeer() {
					_, ok := rerouteMap[mainOriginatorId]
					if ok {
						allExceptLast := len(rerouteMap[mainOriginatorId])
						if Contains(rerouteMap[mainOriginatorId][:allExceptLast], nodeId) {
							continue
						} else {
							currDist = dist
							nextNodeId = nodeId
						}
					} else {
						currDist = dist
						nextNodeId = nodeId
					}
				} else {
					currDist = dist
					nextNodeId = nodeId
				}
			}
		} else {
			thresholdFailed = true
			if Constants.GetPaymentEnabled() {
				if dist < payDist {
					payDist = dist
					payNextId = nodeId
				}
			}
			listItem := Threshold{firstNodeId, nodeId}
			thresholdList = append(thresholdList, listItem)
		}

	}

	if nextNodeId != 0 {
		thresholdFailed = false
		accessFailed = false
	} else {
		if !thresholdFailed {
			accessFailed = true
			nextNodeId = -2 // Access Failed
		} else {
			nextNodeId = -1 // Threshold Failed
		}
		if Constants.GetPaymentEnabled() {
			if payNextId != 0 {
				accessFailed = false
				if Constants.IsOnlyOriginatorPays() {
					if firstNodeId == mainOriginatorId {
						payment.IsOriginator = true
						payment.FirstNodeId = firstNodeId
						payment.PayNextId = payNextId
						payment.ChunkId = chunkId
						nextNodeId = payNextId
					} else {
						thresholdFailed = true
						nextNodeId = -1
					}
				} else if Constants.IsPayIfOrigPays() {
					if prevNodePaid {
						nextNodeId = payNextId
						thresholdFailed = false
						if firstNodeId == mainOriginatorId {
							payment.IsOriginator = true
						} else {
							payment.IsOriginator = false
						}
						payment.FirstNodeId = firstNodeId
						payment.PayNextId = payNextId
						payment.ChunkId = chunkId
					} else {
						if firstNodeId == mainOriginatorId {
							payment.IsOriginator = true
							payment.FirstNodeId = firstNodeId
							payment.PayNextId = payNextId
							payment.ChunkId = chunkId
							nextNodeId = payNextId
						} else {
							thresholdFailed = true
							nextNodeId = -1
							payNextId = 0
						}
					}
				} else {
					nextNodeId = payNextId
					thresholdFailed = false
					if firstNodeId == mainOriginatorId {
						payment.IsOriginator = true
					} else {
						payment.IsOriginator = false
					}
					payment.FirstNodeId = firstNodeId
					payment.PayNextId = payNextId
					payment.ChunkId = chunkId
				}
			}
		}
	}
	// TODO: Usikker på dette
	if Constants.GetPaymentEnabled() {
	out:
		for i, item := range thresholdList {
			for _, nodeId := range item {
				if nodeId == payNextId {
					if Constants.IsPayIfOrigPays() {
						if firstNodeId == mainOriginatorId {
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
	if payment != (Payment{}) {
		prevNodePaid = true
	} else {
		prevNodePaid = false
	}
	// RASMUS: nil reference error
	if nextNodeId != 0 {
		//fmt.Printf("\n next node is: %d", nextNodeId)
	}
	return nextNodeId, thresholdList, thresholdFailed, accessFailed, payment, prevNodePaid
}

// ConsumeTask cacheDict is map of nodes containing an array of maps with key as a chunkAddr and a popularity counter
func ConsumeTask(request *Request, graph *Graph, respNodes [4]int, rerouteMap RerouteMap, cacheMap CacheMap) (bool, Route, [][]Threshold, bool, []Payment) {
	var thresholdFailedList [][]Threshold
	var paymentList []Payment
	originatorId := request.OriginatorId
	chunkId := request.ChunkId
	mainOriginatorId := originatorId
	found := false
	foundByCaching := false
	route := Route{mainOriginatorId}
	//var resultInt int
	var nextNodeId int
	var thresholdList []Threshold
	// thresholdFailed := false
	var accessFailed bool
	var payment Payment
	var prevNodePaid bool

	if Constants.IsPayIfOrigPays() {
		prevNodePaid = true
	}
	if ArrContains(respNodes, mainOriginatorId) {
		// originator has the chunk
		found = true
	} else {
		counter := 0
	out:
		for !ArrContains(respNodes, originatorId) {
			counter++
			//fmt.Printf("\n orig: %d, chunk_id: %d", mainOriginatorId, chunkId)
			// nextNodeId, thresholdList, thresholdFailed, accessFailed, payment, prevNodePaid = getNext(originator, chunkId, graph, mainOriginator, prevNodePaid, rerouteMap)
			nextNodeId, thresholdList, _, accessFailed, payment, prevNodePaid = getNext(originatorId, chunkId, graph, mainOriginatorId, prevNodePaid, rerouteMap)
			//fmt.Println("payment: ", payment)
			if payment != (Payment{}) {
				paymentList = append(paymentList, payment)
			}
			if len(thresholdList) > 0 {
				thresholdFailedList = append(thresholdFailedList, thresholdList)
			}
			// RASMUS: Nil reference error
			if nextNodeId != 0 {
				route = append(route, nextNodeId)
			}
			// if not isinstance(next_node, int), originale versjonen
			if !(nextNodeId <= -1) {
				if ArrContains(respNodes, nextNodeId) {
					//fmt.Println("is not in cache")
					found = true
					break out
				}
				if Constants.IsCacheEnabled() {
					nextNode := graph.NodesMap[nextNodeId]
					chunkMap, ok := cacheMap[nextNode]
					if ok {
						if chunkMap[chunkId] > 1 {
							//fmt.Println("is in cache")
							found = true
							foundByCaching = true
							break out
						}
					}
				}
				// NOTE !
				originatorId = nextNodeId
			} else {
				break out
			}
		}
	}

	route = append(route, chunkId)

	if Constants.IsForwarderPayForceOriginatorToPay() {
		//if nextNodeId != -2 {
		if !Contains(route, -2) {
			// NOT accessFailed
			if len(paymentList) > 0 {
				firstPayment := paymentList[0]
				if !firstPayment.IsOriginator {
					// TODO: Dobbelsjekk at logikken under her matcher originalen
					for i := range route[:len(route)-1] {
						p := Payment{FirstNodeId: route[i], PayNextId: route[i+1], ChunkId: route[len(route)-1]}
						for j, tmp := range paymentList {
							if p.PayNextId == tmp.PayNextId && p.FirstNodeId == tmp.FirstNodeId && p.ChunkId == tmp.ChunkId {
								break
							}
							if j == len(paymentList) {
								// payment is now definitely not in paymentList
								if i == 0 {
									p.IsOriginator = true
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
						p := Payment{FirstNodeId: route[i], PayNextId: route[i+1], ChunkId: route[len(route)-1]}
						for j, tmp := range paymentList {
							if p.PayNextId == tmp.PayNextId && p.FirstNodeId == tmp.FirstNodeId && p.ChunkId == tmp.ChunkId {
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

func getProximityChunk(firstNodeId int, chunkId int) int {
	retVal := Constants.GetBits() - BitLength(firstNodeId^chunkId)
	if retVal <= Constants.GetMaxProximityOrder() {
		return retVal
	} else {
		return Constants.GetMaxProximityOrder()
	}
}

func PeerPriceChunk(firstNodeId int, chunkId int) int {
	val := (Constants.GetMaxProximityOrder() - getProximityChunk(firstNodeId, chunkId) + 1) * Constants.GetPrice()
	return val
}

func CreateDownloadersList(g *Graph) []int {
	//fmt.Println("Creating downloaders list...")

	downloadersList := Choice(g.NodeIds, Constants.GetOriginators())

	//fmt.Println("Downloaders list create...!")
	return downloadersList
}

func CreateNodesList(g *Graph) []int {
	//fmt.Println("Creating nodes list...")
	nodesValue := g.NodeIds
	//fmt.Println("Nodes list create...!")
	return nodesValue
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

// TODO: Not used in original
//func MakeFiles() []int {
//	fmt.Println("Making files...")
//	var filesList []int
//
//	for i := 0; i <= ct.Constants.GetOriginators(); i++ {
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
