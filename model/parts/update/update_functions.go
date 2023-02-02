package update

import (
	"encoding/json"
	"fmt"
	. "go-incentive-simulation/model/constants"
	. "go-incentive-simulation/model/general"
	. "go-incentive-simulation/model/parts/types"
	. "go-incentive-simulation/model/parts/utils"
	"io/ioutil"
	"math"
)

func UpdateSuccessfulFound(prevState State, policyInput Policy) State {
	if policyInput.Found {
		prevState.SuccessfulFound++
	}
	return prevState
}

func UpdateFailedRequestsThreshold(prevState State, policyInput Policy) State {
	found := policyInput.Found
	// thresholdFailedList := policyInput.thresholdFailedList
	accessFailed := policyInput.AccessFailed
	if !found && !accessFailed {
		prevState.FailedRequestsThreshold++
	}
	return prevState
}

func UpdateFailedRequestsAccess(prevState State, policyInput Policy) State {
	accessFailed := policyInput.AccessFailed
	if accessFailed {
		prevState.FailedRequestsAccess++
	}
	return prevState
}

func UpdateOriginatorIndex(prevState State, policyInput Policy) State {
	if prevState.OriginatorIndex+1 >= Constants.GetOriginators() {
		prevState.OriginatorIndex = 0
		return prevState
	}
	prevState.OriginatorIndex++
	return prevState
}

func convertAndDumpToFile1(routes []Route, currTimestep int) error {
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

func UpdateRouteListAndFlush(prevState State, policyInput Policy) State {
	prevState.RouteLists = append(prevState.RouteLists, policyInput.Route)
	prevState.TimeStep++
	if prevState.TimeStep%6250 == 0 {
		convertAndDumpToFile1(prevState.RouteLists, prevState.TimeStep)
		prevState.RouteLists = []Route{}
		return prevState
	}

	return prevState
}

// TODO: Implement this function
// Does not work with current type of cacheMap
func UpdateCacheMap(prevState State, policyInput Policy) State {
	cacheMap := prevState.CacheStruct.CacheMap
	cacheHits := prevState.CacheStruct.CacheHits
	g := prevState.Graph
	chunkAddr := 0
	//val := make(map[int]int)

	if Constants.IsCacheEnabled() {
		route := policyInput.Route
		if Contains(route, -3) {
			chunkAddr = route[len(route)-2]
		} else {
			chunkAddr = route[len(route)-1]
		}
		if !Contains(route, -1) && !Contains(route, -2) {
			if Contains(route, -3) {
				for i := 0; i < len(route)-3; i++ {
					routeId := g.GetNode(route[i])
					cacheHits++
					if _, ok := cacheMap[routeId]; ok {
						val := cacheMap[routeId]

						if _, ok := val[chunkAddr]; ok {
							val[chunkAddr]++
						} else {
							val[chunkAddr] = 1
						}
					} else {
						cacheMap[routeId] = map[int]int{}
						val := cacheMap[routeId]
						val[chunkAddr] = 1
					}
				}
			} else {
				for i := 0; i < len(route)-2; i++ {
					routeId := g.GetNode(route[i])
					cacheHits++
					if _, ok := cacheMap[routeId]; ok {
						val := cacheMap[routeId]

						if _, ok := val[chunkAddr]; ok {
							val[chunkAddr]++
						} else {
							val[chunkAddr] = 1
						}
					} else {
						cacheMap[routeId] = map[int]int{}
						val := cacheMap[routeId]
						val[chunkAddr] = 1
					}
				}
			}
		}
	}
	prevState.CacheStruct.CacheMap = cacheMap
	prevState.CacheStruct.CacheHits = cacheHits
	return prevState
}

func UpdateRerouteMap(prevState State, policyInput Policy) State {
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
	return prevState
}

func UpdatePendingMap(prevState State, policyInput Policy) State {
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
	return prevState
}

func UpdateNetwork(prevState State, policyInput Policy) State {
	network := prevState.Graph
	currTimeStep := prevState.TimeStep + 1
	route := policyInput.Route
	paymentsList := policyInput.PaymentList

	if Constants.GetPaymentEnabled() {
		for _, payment := range paymentsList {
			var p Payment
			if payment != p {
				if payment.FirstNodeId != -1 {
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
							edgeData1.A2B = 0
							edgeData2.A2B = 0
						}
					}
					fmt.Println("Payment from ", payment.FirstNodeId, " to ", payment.PayNextId, " for chunk ", payment.ChunkId, " with price ", val)
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
							edgeData1.A2B = 0
							edgeData2.A2B = 0
						}
					}
					fmt.Println("-1", "Payment from ", payment.FirstNodeId, " to ", payment.PayNextId, " for chunk ", payment.ChunkId, " with price ", val) //Means that the first one is the originator
				}
			}
		}
	}
	if !Contains(route, -1) && !Contains(route, -2) {
		routeWithPrice := []int{}
		if Contains(route, -3) {
			chunkId := route[len(route)-2]
			for i := 0; i < len(route)-3; i++ {
				requesterNode := route[i]
				providerNode := route[i+1]
				price := PeerPriceChunk(providerNode, chunkId)
				edgeData1 := network.GetEdgeData(requesterNode, providerNode)
				edgeData1.A2B += price
				if Constants.GetMaxPOCheckEnabled() {
					routeWithPrice = append(routeWithPrice, requesterNode)
					routeWithPrice = append(routeWithPrice, price)
					routeWithPrice = append(routeWithPrice, providerNode)
				}
			}
			if Constants.GetMaxPOCheckEnabled() {
				fmt.Println("Route with price ", routeWithPrice)
			}
		} else {
			chunkId := route[len(route)-1]
			for i := 0; i < len(route)-2; i++ {
				requesterNode := route[i]
				providerNode := route[i+1]
				price := PeerPriceChunk(providerNode, chunkId)
				edgeData1 := network.GetEdgeData(requesterNode, providerNode)
				edgeData1.A2B += price
				if Constants.GetMaxPOCheckEnabled() {
					routeWithPrice = append(routeWithPrice, requesterNode)
					routeWithPrice = append(routeWithPrice, price)
					routeWithPrice = append(routeWithPrice, providerNode)
				}
			}
			if Constants.GetMaxPOCheckEnabled() {
				fmt.Println("Route with price ", routeWithPrice)
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
						edgeData1 := network.GetEdgeData(requesterNode, providerNode)
						passedTime := (currTimeStep - edgeData1.Last) / Constants.GetRequestsPerSecond()
						if passedTime > 0 {
							refreshRate := Constants.GetRefreshRate()
							if Constants.IsAdjustableThreshold() {
								refreshRate = int(math.Ceil(float64(edgeData1.Threshold / 2)))
							}
							removedDeptAmount := passedTime * refreshRate
							edgeData1.A2B -= removedDeptAmount
							if edgeData1.A2B < 0 {
								edgeData1.A2B = 0
							}
							edgeData1.Last = currTimeStep
						}
					}
				}
			}
		}
	}
	prevState.Graph = network
	return prevState
}
