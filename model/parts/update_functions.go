package policy

import (
	"fmt"
	ct "go-incentive-simulation/model/constants"
	. "go-incentive-simulation/model/parts/utils"
)

type State struct {
	network                 *Graph
	originators             []int
	originatorsIndex        int
	nodesId                 []int
	routeLists              []Route
	pendingDict             map[int]int
	rerouteDict             map[int][]int
	cacheDict               map[int]int
	originatorIndex         int
	successfulFound         int
	failedRequestsThreshold int
	failedRequestsAccess    int
	timeStep                int
}

type Policy struct {
	found           bool
	route           Route
	thresholdFailed []int
	originatorIndex int
	accessFailed    bool
	paymentList     []Payment
}

func UpdateSuccessfulFound(prevState State, policyInput Policy) State {
	oldSuccessCounter := prevState.successfulFound
	newSuccessCounter := oldSuccessCounter
	if policyInput.found {
		newSuccessCounter++
	}
	return prevState
}

func UpdateFailedRequestsThreshold(prevState State, policyInput Policy) State {
	oldFailedCounter := prevState.failedRequestsThreshold
	newFailedCounter := oldFailedCounter
	found := policyInput.found
	// thresholdFailed := policyInput.thresholdFailed
	accessFailed := policyInput.accessFailed
	if !found && !accessFailed {
		newFailedCounter++
	}
	return prevState
}

func UpdateFailedRequestsAccess(prevState State, policyInput Policy) State {
	oldFailedAccessCounter := prevState.failedRequestsAccess
	accessFailed := policyInput.accessFailed
	if accessFailed {
		oldFailedAccessCounter++
	}
	return prevState
}

func UpdateOriginatorIndex(prevState State, policyInput Policy) State {
	oldOriginatorIndex := prevState.originatorIndex
	newOriginatorIndex := oldOriginatorIndex + 1
	if newOriginatorIndex >= ct.Constants.GetOriginators() {
		newOriginatorIndex = 0
	}
	return prevState
}

// TODO: function convert and dump to file

func UpdateRouteListAndFlush(prevState State, policyInput Policy) State {
	prevState.routeLists = append(prevState.routeLists, policyInput.route)
	currTimestep := prevState.timeStep + 1
	if currTimestep%6250 == 0 {
		// TODO: call convert_and_dump
		prevState.routeLists = []Route{}
		return prevState
	}
	return prevState
}

// TODO: Implement this function
func UpdateCacheDictionary(prevState State, policyInput Policy) State {
	return prevState
}

func UpdateRerouteDictionary(prevState State, policyInput Policy) State {
	rerouteDict := prevState.rerouteDict
	if ct.Constants.IsRetryWithAnotherPeer() {
		route := policyInput.route
		originator := route[0]
		if !contains(route, -1) && !contains(route, -2) {
			if _, ok := rerouteDict[originator]; ok {
				val := rerouteDict[originator]
				if val[len(val)-1] == route[len(route)-1] {
					//remove rerouteDict[originator]
					delete(rerouteDict, originator)
				}
			}
		} else {
			if len(route) > 3 {
				if _, ok := rerouteDict[originator]; ok {
					val := rerouteDict[originator]
					if !contains(val, route[1]) {
						val = append([]int{route[1]}, val...)
						rerouteDict[originator] = val
					}
				} else {
					rerouteDict[originator] = []int{route[1], route[len(route)-1]}
				}
			}
		}
		if _, ok := rerouteDict[originator]; ok {
			if len(rerouteDict[originator]) > ct.Constants.GetBinSize() {
				delete(rerouteDict, originator)
			}
		}
	}
	return prevState
}

func UpdatePendingDictionary(prevState State, policyInput Policy) State {
	pendingDict := prevState.pendingDict
	if ct.Constants.IsWaitingEnabled() {
		route := policyInput.route
		originator := route[0]
		if !contains(route, -1) && !contains(route, -2) {
			if _, ok := pendingDict[originator]; ok {
				if pendingDict[originator] == route[len(route)-1] {
					delete(pendingDict, originator)
				}
			}

		} else {
			pendingDict[originator] = route[len(route)-1]
		}
	}
	return prevState
}

func UpdateNetwork(prevState State, policyInput Policy) State {
	network := prevState.network
	currTinmeStep := prevState.timeStep + 1
	route := policyInput.route
	paymentsList := policyInput.paymentList

	if ct.Constants.GetPaymentEnabled() {
		for _, payment := range paymentsList {
			var p Payment
			if payment != p {
				if payment.firstNodeId != -1 {
					edgeData1 := network.GetEdgeData(payment.firstNodeId, payment.payNextId)
					edgeData2 := network.GetEdgeData(payment.payNextId, payment.firstNodeId)
					price := peerPriceChunk(payment.payNextId, payment.chunkId)
					val := edgeData1.a2b - edgeData2.a2b + price
					if ct.Constants.IsPayOnlyForCurrentRequest() {
						val = price
					}
					if val < 0 {
						continue
					} else {
						if !ct.Constants.IsPayOnlyForCurrentRequest() {
							edgeData1.a2b = 0
							edgeData2.a2b = 0
						}
					}
					fmt.Println("Payment from ", payment.firstNodeId, " to ", payment.payNextId, " for chunk ", payment.chunkId, " with price ", val)
				} else {
					edgeData1 := network.GetEdgeData(payment.firstNodeId, payment.payNextId)	
					edgeData2 := network.GetEdgeData(payment.payNextId, payment.firstNodeId)
					price := peerPriceChunk(payment.payNextId, payment.chunkId)
					val := edgeData1.a2b - edgeData2.a2b + price
					if ct.Constants.IsPayOnlyForCurrentRequest() {
						val = price
					}
					if val < 0 {
						continue
					} else {
						if !ct.Constants.IsPayOnlyForCurrentRequest() {
							edgeData1.a2b = 0
							edgeData2.a2b = 0
						}
					}
					fmt.Println("-1", "Payment from ", payment.firstNodeId, " to ", payment.payNextId, " for chunk ", payment.chunkId, " with price ", val) //Means that the first one is the originator
				}			
			}
		}
	}
	
	return prevState
}


func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
