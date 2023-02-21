package update

import (
	"encoding/json"
	. "go-incentive-simulation/model/constants"
	. "go-incentive-simulation/model/general"
	. "go-incentive-simulation/model/parts/types"
	. "go-incentive-simulation/model/parts/utils"
	"io/ioutil"
	"math"
	"sync"
	"sync/atomic"
)

func UpdateSuccessfulFound(prevState *State, policyInput Policy) {
	if policyInput.Found {
		atomic.AddInt32(&prevState.SuccessfulFound, 1)
	}
}

func UpdateFailedRequestsThreshold(prevState *State, policyInput Policy) {
	found := policyInput.Found
	// thresholdFailedList := policyInput.thresholdFailedList
	accessFailed := policyInput.AccessFailed
	if !found && !accessFailed {
		atomic.AddInt32(&prevState.FailedRequestsThreshold, 1)

	}
}

func UpdateFailedRequestsAccess(prevState *State, policyInput Policy) {
	accessFailed := policyInput.AccessFailed
	if accessFailed {
		atomic.AddInt32(&prevState.FailedRequestsAccess, 1)
	}
}

func UpdateOriginatorIndex(prevState *State, policyInput Policy, origMutex *sync.Mutex) {
	origMutex.Lock()
	defer origMutex.Unlock()
	//if prevState.TimeStep%100 == 0 {
	if int(atomic.LoadInt32(&prevState.OriginatorIndex)+1) >= Constants.GetOriginators() {
		atomic.StoreInt32(&prevState.OriginatorIndex, 0)
	} else {
		atomic.AddInt32(&prevState.OriginatorIndex, 1)
	}
	//}
	//prevState.OriginatorIndex = rand.Intn(Constants.GetOriginators() - 1)
}

func convertAndDumpToFile(routes []Route, currTimestep int) error {
	type RouteData struct {
		Timestep int     `json:"timestep"`
		Routes   []Route `json:"routes"`
	}
	data := RouteData{currTimestep, routes}
	file, _ := json.Marshal(data)
	err := ioutil.WriteFile("routes.json", file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func UpdateRouteListAndFlush(prevState *State, policyInput Policy) {
	prevState.RouteLists = append(prevState.RouteLists, policyInput.Route)
	currTimestep := int(atomic.LoadInt32(&prevState.TimeStep) + 1)
	if currTimestep%6250 == 0 {
		convertAndDumpToFile(prevState.RouteLists, currTimestep)
		prevState.RouteLists = []Route{}
	}
}

func UpdateCacheMap(prevState *State, policyInput Policy) {
	cacheStruct := prevState.CacheStruct
	chunkId := 0
	//val := make(map[int]int)

	if Constants.IsCacheEnabled() {
		route := policyInput.Route
		if Contains(route, -3) {
			// -3 means found by caching
			cacheStruct.CacheHits++
			chunkId = route[len(route)-2]
		} else {
			chunkId = route[len(route)-1]
		}
		if !Contains(route, -1) && !Contains(route, -2) {
			if Contains(route, -3) {
				for i := 0; i < len(route)-3; i++ {
					nodeId := route[i]
					//cacheStruct.AddToCache(nodeId, chunkId)
					node := prevState.Graph.GetNode(nodeId)
					node.Mutex.Lock()
					cacheMap := node.CacheMap
					if cacheMap != nil {
						if _, ok := cacheMap[chunkId]; ok {
							cacheMap[chunkId]++
						} else {
							cacheMap[chunkId] = 1
						}
					} else {
						node.CacheMap = map[int]int{node.Id: 1}
					}
					node.Mutex.Unlock()
				}
			} else {
				for i := 0; i < len(route)-2; i++ {
					nodeId := route[i]
					//cacheStruct.AddToCache(nodeId, chunkId)
					node := prevState.Graph.GetNode(nodeId)
					node.Mutex.Lock()
					cacheMap := node.CacheMap
					if cacheMap != nil {
						if _, ok := cacheMap[chunkId]; ok {
							cacheMap[chunkId]++
						} else {
							cacheMap[chunkId] = 1
						}
					} else {
						node.CacheMap = map[int]int{node.Id: 1}
					}
					node.Mutex.Unlock()
				}
			}
		}
	}
	prevState.CacheStruct = cacheStruct
}

func UpdateRerouteMap(prevState *State, policyInput Policy) {
	rerouteMap := prevState.RerouteMap
	if Constants.IsRetryWithAnotherPeer() {
		route := policyInput.Route
		originator := route[0]
		if !Contains(route, -1) && !Contains(route, -2) {
			if _, ok := rerouteMap[originator]; ok {
				val := rerouteMap[originator]
				if val[len(val)-1] == route[len(route)-1] {
					//remove rerouteMap[originator]
					delete(rerouteMap, originator)
				}
			}
		} else {
			if len(route) > 3 {
				if _, ok := rerouteMap[originator]; ok {
					val := rerouteMap[originator]
					if !Contains(val, route[1]) {
						val = append([]int{route[1]}, val...)
						rerouteMap[originator] = val
					}
				} else {
					rerouteMap[originator] = []int{route[1], route[len(route)-1]}
				}
			}
		}
		if _, ok := rerouteMap[originator]; ok {
			if len(rerouteMap[originator]) > Constants.GetBinSize() {
				delete(rerouteMap, originator)
			}
		}
	}
	prevState.RerouteMap = rerouteMap
}

func UpdatePendingMap(prevState *State, policyInput Policy) {
	pendingMap := prevState.PendingMap
	if Constants.IsWaitingEnabled() {
		route := policyInput.Route
		originator := route[0]
		if !Contains(route, -1) && !Contains(route, -2) {
			if _, ok := pendingMap[originator]; ok {
				if pendingMap[originator] == route[len(route)-1] {
					delete(pendingMap, originator)
				}
			}

		} else {
			pendingMap[originator] = route[len(route)-1]
		}
	}
	prevState.PendingMap = pendingMap
}

func UpdateNetwork(prevState *State, policyInput Policy, networkMutex *sync.Mutex) {
	network := prevState.Graph
	//atomic.AddInt32(&prevState.TimeStep, 1)
	currTimeStep := int(atomic.AddInt32(&prevState.TimeStep, 1))
	route := policyInput.Route
	paymentsList := policyInput.PaymentList

	//networkMutex.Lock()
	//defer networkMutex.Unlock()
	if Constants.GetPaymentEnabled() {
		for _, payment := range paymentsList {
			if payment != (Payment{}) {
				if !payment.IsOriginator {
					edgeData1 := network.GetEdgeData(payment.FirstNodeId, payment.PayNextId)
					edgeData2 := network.GetEdgeData(payment.PayNextId, payment.FirstNodeId)
					price := PeerPriceChunk(payment.PayNextId, payment.ChunkId)
					val := edgeData1.A2B - edgeData2.A2B + price
					if Constants.IsPayOnlyForCurrentRequest() {
						val = price
					}
					if val < 0 {
						continue
					} else {
						if !Constants.IsPayOnlyForCurrentRequest() {
							//edgeData1.A2B = 0
							//edgeData2.A2B = 0
							newEdgeData1 := edgeData1
							newEdgeData1.A2B = 0
							network.SetEdgeData(payment.FirstNodeId, payment.PayNextId, newEdgeData1)

							newEdgeData2 := edgeData2
							newEdgeData2.A2B = 0
							network.SetEdgeData(payment.PayNextId, payment.FirstNodeId, newEdgeData2)
						}
					}
					// fmt.Println("Payment from ", payment.FirstNodeId, " to ", payment.PayNextId, " for chunk ", payment.ChunkId, " with price ", val)
				} else {
					edgeData1 := network.GetEdgeData(payment.FirstNodeId, payment.PayNextId)
					edgeData2 := network.GetEdgeData(payment.PayNextId, payment.FirstNodeId)
					price := PeerPriceChunk(payment.PayNextId, payment.ChunkId)
					val := edgeData1.A2B - edgeData2.A2B + price
					if Constants.IsPayOnlyForCurrentRequest() {
						val = price
					}
					if val < 0 {
						continue
					} else {
						if !Constants.IsPayOnlyForCurrentRequest() {
							//edgeData1.A2B = 0
							//edgeData2.A2B = 0
							newEdgeData1 := edgeData1
							newEdgeData1.A2B = 0
							network.SetEdgeData(payment.FirstNodeId, payment.PayNextId, newEdgeData1)

							newEdgeData2 := edgeData2
							newEdgeData2.A2B = 0
							network.SetEdgeData(payment.PayNextId, payment.FirstNodeId, newEdgeData2)
						}
					}
					//fmt.Println("-1", "Payment from ", payment.FirstNodeId, " to ", payment.PayNextId, " for chunk ", payment.ChunkId, " with price ", val) //Means that the first one is the originator
				}
			}
		}
	}
	if !Contains(route, -1) && !Contains(route, -2) {
		var routeWithPrice []int
		if Contains(route, -3) {
			chunkId := route[len(route)-2]
			for i := 0; i < len(route)-3; i++ {
				requesterNode := route[i]
				providerNode := route[i+1]
				price := PeerPriceChunk(providerNode, chunkId)
				edgeData := network.GetEdgeData(requesterNode, providerNode)
				//edgeData1.A2B += price
				newEdgeData := edgeData
				newEdgeData.A2B += price
				network.SetEdgeData(requesterNode, providerNode, newEdgeData)

				if Constants.GetMaxPOCheckEnabled() {
					routeWithPrice = append(routeWithPrice, requesterNode)
					routeWithPrice = append(routeWithPrice, price)
					routeWithPrice = append(routeWithPrice, providerNode)
				}
			}
			if Constants.GetMaxPOCheckEnabled() {
				//fmt.Println("Route with price ", routeWithPrice)
			}
		} else {
			chunkId := route[len(route)-1]
			for i := 0; i < len(route)-2; i++ {
				requesterNode := route[i]
				providerNode := route[i+1]
				price := PeerPriceChunk(providerNode, chunkId)
				edgeData := network.GetEdgeData(requesterNode, providerNode)
				//edgeData.A2B += price
				newEdgeData := edgeData
				newEdgeData.A2B += price
				network.SetEdgeData(requesterNode, providerNode, newEdgeData)

				if Constants.GetMaxPOCheckEnabled() {
					routeWithPrice = append(routeWithPrice, requesterNode)
					routeWithPrice = append(routeWithPrice, price)
					routeWithPrice = append(routeWithPrice, providerNode)
				}
			}
			if Constants.GetMaxPOCheckEnabled() {
				//fmt.Println("Route with price ", routeWithPrice)
			}
		}
	}
	if Constants.GetThresholdEnabled() && Constants.IsForgivenessEnabled() {
		thresholdFailedLists := policyInput.ThresholdFailedLists
		if len(thresholdFailedLists) > 0 {
			for _, thresholdFailedL := range thresholdFailedLists {
				if len(thresholdFailedL) > 0 {
					for _, couple := range thresholdFailedL {
						requesterNode := couple[0]
						providerNode := couple[1]
						edgeData := network.GetEdgeData(requesterNode, providerNode)
						passedTime := (currTimeStep - edgeData.Last) / Constants.GetRequestsPerSecond()
						if passedTime > 0 {
							refreshRate := Constants.GetRefreshRate()
							if Constants.IsAdjustableThreshold() {
								refreshRate = int(math.Ceil(float64(edgeData.Threshold / 2)))
							}
							removedDeptAmount := passedTime * refreshRate
							//edgeData.A2B -= removedDeptAmount
							//if edgeData.A2B < 0 {
							//	edgeData.A2B = 0
							//}
							//edgeData.Last = currTimeStep
							newEdgeData := edgeData
							newEdgeData.A2B -= removedDeptAmount
							if newEdgeData.A2B < 0 {
								newEdgeData.A2B = 0
							}
							newEdgeData.Last = currTimeStep
							network.SetEdgeData(requesterNode, providerNode, newEdgeData)
						}
					}
				}
			}
		}
	}
	// Unlocks all the edges between the nodes in the route
	//networkMutex.Lock()
	if Constants.GetEdgeLock() {
		if !Contains(route, -1) && !Contains(route, -2) {
			if Contains(route, -3) {
				for i := 0; i < len(route)-3; i++ {
					prevState.Graph.UnlockEdge(route[i], route[i+1])
				}
			} else {
				for i := 0; i < len(route)-2; i++ {
					prevState.Graph.UnlockEdge(route[i], route[i+1])
				}
			}
		} else {
			for i := 0; i < len(route)-3; i++ {
				prevState.Graph.UnlockEdge(route[i], route[i+1])
			}
		}
	}
	//networkMutex.Unlock()

	prevState.Graph = network
}
