package main

import (
	"fmt"
	. "go-incentive-simulation/model/parts/policy"
	. "go-incentive-simulation/model/parts/types"
	. "go-incentive-simulation/model/parts/update"
	. "go-incentive-simulation/model/state"
	"sync"
	"time"
)

func MakePolicyOutput(state State, index int) Policy {
	//fmt.Println("start of make initial policy")

	found, route, thresholdFailed, accessFailed, paymentsList := SendRequest(&state)
	policy := Policy{
		Found:                found,
		Route:                route,
		ThresholdFailedLists: thresholdFailed,
		AccessFailed:         accessFailed,
		PaymentList:          paymentsList,
	}
	return policy
}

func main() {
	start := time.Now()
	state := MakeInitialState("./data/nodes_data_16_10000.txt")

	const iterations = 1000000
	const numGoroutines = 1

	numLoops := iterations / numGoroutines
	stateArray := make([]State, numLoops)

	for i := 0; i < numLoops; i++ {
		//fmt.Println("Start of lop ", time.Since(start))
		var policyOutputs [numGoroutines]Policy
		var wg sync.WaitGroup
		for j := 0; j < numGoroutines; j++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				policyOutputs[index] = MakePolicyOutput(state, index)
			}(j)
		}
		wg.Wait()
		//fmt.Println("end of lop ", time.Since(start))
		for j := 0; j < numGoroutines; j++ {
			state = UpdatePendingMap(state, policyOutputs[j])
			state = UpdateRerouteMap(state, policyOutputs[j])
			state = UpdateCacheMap(state, policyOutputs[j])
			state = UpdateOriginatorIndex(state, policyOutputs[j])
			state = UpdateSuccessfulFound(state, policyOutputs[j])
			state = UpdateFailedRequestsThreshold(state, policyOutputs[j])
			state = UpdateFailedRequestsAccess(state, policyOutputs[j])
			state = UpdateRouteListAndFlush(state, policyOutputs[j])
			state = UpdateNetwork(state, policyOutputs[j])
		}

		curState := State{
			Graph:                   state.Graph,
			Originators:             state.Originators,
			NodesId:                 state.NodesId,
			RouteLists:              state.RouteLists,
			PendingMap:              state.PendingMap,
			RerouteMap:              state.RerouteMap,
			CacheStruct:             state.CacheStruct,
			OriginatorIndex:         state.OriginatorIndex,
			SuccessfulFound:         state.SuccessfulFound,
			FailedRequestsThreshold: state.FailedRequestsThreshold,
			FailedRequestsAccess:    state.FailedRequestsAccess,
			TimeStep:                state.TimeStep,
			PriceList:               state.PriceList}

		stateArray[i] = curState
		//fmt.Println("payment list: ", policyOutputs[0].PaymentList)
	}
	fmt.Println("end of main: ")
	elapsed := time.Since(start)
	fmt.Println("Time taken:", elapsed)
	fmt.Println("Number of iterations: ", iterations)
	fmt.Println("Number of Goroutines: ", numGoroutines)
	// allReq, thresholdFails, requestsToBucketZero, rejectedBucketZero, rejectedFirstHop := ReadRoutes("routes.json")
	// fmt.Println("allReq: ", allReq)
	// fmt.Println("thresholdFails: ", thresholdFails)
	// fmt.Println("requestsToBucketZero: ", requestsToBucketZero)
	// fmt.Println("rejectedBucketZero: ", rejectedBucketZero)
	// fmt.Println("rejectedFirstHop: ", rejectedFirstHop)
	PrintState(state)
}

func PrintState(state State) {
	fmt.Println("SuccessfulFound: ", state.SuccessfulFound)
	fmt.Println("FailedRequestsThreshold: ", state.FailedRequestsThreshold)
	fmt.Println("FailedRequestsAccess: ", state.FailedRequestsAccess)
	fmt.Println("CacheHits:", state.CacheStruct.CacheHits)
	fmt.Println("TimeStep: ", state.TimeStep)
	fmt.Println("OriginatorIndex: ", state.OriginatorIndex)
	//get average of all the prices
	avgPrice := 0
	for _, price := range state.PriceList {
		avgPrice += price
	}
	avgPrice = avgPrice / len(state.PriceList)
	fmt.Println("Average Price: ", avgPrice)
	fmt.Println("Length of price list: ", len(state.PriceList))

	//fmt.Println("PendingMap: ", state.PendingMap)
	//fmt.Println("RerouteMap: ", state.RerouteMap)
	//fmt.Println("RouteLists: ", state.RouteLists)
	//fmt.Println("CacheMap: ", state.CacheStruct.CacheMap)
}
