package main

import (
	"fmt"
	. "go-incentive-simulation/model/constants"
	. "go-incentive-simulation/model/parts/policy"
	. "go-incentive-simulation/model/parts/types"
	. "go-incentive-simulation/model/parts/update"
	. "go-incentive-simulation/model/state"
	"sync"
	"sync/atomic"
	"time"
)

func MakePolicyOutput(state *State, index int) Policy {
	//fmt.Println("start of make initial policy")

	//found, route, thresholdFailed, accessFailed, paymentsList := SendRequest(&state)
	found, route, thresholdFailed, accessFailed, paymentsList := SendRequest(state, index)

	policy := Policy{
		Found:                found,
		Route:                route,
		ThresholdFailedLists: thresholdFailed,
		AccessFailed:         accessFailed,
		PaymentList:          paymentsList,
	}
	return policy
}

func UpdateWorker(stateChan chan *State, policyChan chan Policy, globalState *State, stateArray []State, origMutex *sync.Mutex, networkMutex *sync.Mutex) {

	for {
		policyOutput := <-policyChan

		UpdatePendingMap(globalState, policyOutput)
		UpdateRerouteMap(globalState, policyOutput)
		UpdateCacheMap(globalState, policyOutput)
		UpdateOriginatorIndex(globalState, policyOutput, origMutex)
		UpdateSuccessfulFound(globalState, policyOutput)
		UpdateFailedRequestsThreshold(globalState, policyOutput)
		UpdateFailedRequestsAccess(globalState, policyOutput)
		//UpdateRouteListAndFlush(globalState, policyOutput)
		UpdateNetwork(globalState, policyOutput, networkMutex)

		newState := State{
			Graph:                   globalState.Graph,
			Originators:             globalState.Originators,
			NodesId:                 globalState.NodesId,
			RouteLists:              globalState.RouteLists,
			PendingMap:              globalState.PendingMap,
			RerouteMap:              globalState.RerouteMap,
			CacheStruct:             globalState.CacheStruct,
			OriginatorIndex:         globalState.OriginatorIndex,
			SuccessfulFound:         globalState.SuccessfulFound,
			FailedRequestsThreshold: globalState.FailedRequestsThreshold,
			FailedRequestsAccess:    globalState.FailedRequestsAccess,
			TimeStep:                globalState.TimeStep,
		}

		stateArray[atomic.LoadInt32(&globalState.TimeStep)] = newState

		stateChan <- &newState
	}
}

func main() {
	start := time.Now()
	globalState := MakeInitialState("./data/nodes_data_16_10000.txt")

	const iterations = 10000000
	numGoroutines := Constants.GetNumGoroutines()

	numLoops := iterations / numGoroutines
	stateArray := make([]State, iterations+1)
	stateArray[0] = globalState

	var wg sync.WaitGroup
	policyChan := make(chan Policy, numGoroutines)
	stateChan := make(chan *State, numGoroutines)
	origMutex := &sync.Mutex{}
	networkMutex := &sync.Mutex{}

	for i := 0; i < 5; i++ {
		go UpdateWorker(stateChan, policyChan, &globalState, stateArray, origMutex, networkMutex)
	}

	for j := 0; j < numGoroutines; j++ {
		wg.Add(1)
		go func(index int) {
			curState := &globalState
			for i := 0; i < numLoops; i++ {
				policyChan <- MakePolicyOutput(curState, index)
				// waiting for new state from UpdateWorker
				curState = <-stateChan
			}
			wg.Done()
		}(j)
	}
	wg.Wait()
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
	PrintState(globalState)
}

func PrintState(state State) {
	fmt.Println("SuccessfulFound: ", state.SuccessfulFound)
	fmt.Println("FailedRequestsThreshold: ", state.FailedRequestsThreshold)
	fmt.Println("FailedRequestsAccess: ", state.FailedRequestsAccess)
	fmt.Println("CacheHits:", state.CacheStruct.CacheHits)
	fmt.Println("TimeStep: ", state.TimeStep)
	fmt.Println("OriginatorIndex: ", state.OriginatorIndex)
	//fmt.Println("PendingMap: ", state.PendingMap)
	//fmt.Println("RerouteMap: ", state.RerouteMap)
	//fmt.Println("RouteLists: ", state.RouteLists)
	//fmt.Println("CacheMapArray: ", state.CacheStruct.CacheMapArray)
}
