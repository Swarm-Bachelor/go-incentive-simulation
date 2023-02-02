package main

import (
	"fmt"
	. "go-incentive-simulation/model/parts/policy"
	. "go-incentive-simulation/model/parts/types"
	. "go-incentive-simulation/model/parts/update"
	. "go-incentive-simulation/model/state"
	"sync"
	"time"
	//. "go-incentive-simulation/model/constants"
)

func MakePolicyOutput(state State) Policy {
	fmt.Println("start of make initial policy")
	found, route, thresholdFailed, accessFailed, paymentsList := SendRequest(&state)
	policy := Policy{
		Found:                found,
		Route:                route,
		ThresholdFailedLists: thresholdFailed,
		OriginatorIndex:      state.OriginatorIndex,
		AccessFailed:         accessFailed,
		PaymentList:          paymentsList,
	}
	return policy
}

func main() {
	start := time.Now()
	state := MakeInitialState("./data/nodes_data_8_10000.txt")
	stateArray := []State{state}
	iterations := 250000
	numGoRoutines := 100
	for i := 0; i < iterations/numGoRoutines; i++ {
		policyStruct := PolicyStruct{
			Founds:                   make([]bool, numGoRoutines),
			Routes:                   make([]Route, numGoRoutines),
			ThresholdFailedListsList: make([][][]Threshold, numGoRoutines),
			OriginatorIndices:        make([]int, numGoRoutines),
			AccessFails:              make([]bool, numGoRoutines),
			PaymentListList:          make([][]Payment, numGoRoutines),
		}
		var wg sync.WaitGroup
		wg.Add(numGoRoutines)
		fmt.Println("before PolicyOutput: ", time.Since(start))
		for j := 0; j < numGoRoutines; j++ {
			loop := j
			go func(int) {
				policyOutput := MakePolicyOutput(state)

				policyStruct.Founds[loop] = policyOutput.Found
				policyStruct.Routes[loop] = policyOutput.Route
				policyStruct.ThresholdFailedListsList[loop] = policyOutput.ThresholdFailedLists
				policyStruct.OriginatorIndices[loop] = policyOutput.OriginatorIndex
				policyStruct.AccessFails[loop] = policyOutput.AccessFailed
				policyStruct.PaymentListList[loop] = policyOutput.PaymentList
				wg.Done()
			}(loop)
		}
		wg.Wait()
		fmt.Println("middle: ", time.Since(start))

		state = UpdatePendingMap(state, policyStruct)
		state = UpdateRerouteMap(state, policyStruct)
		state = UpdateCacheMap(state, policyStruct)
		state = UpdateOriginatorIndex(state, policyStruct)
		state = UpdateSuccessfulFound(state, policyStruct)
		state = UpdateFailedRequestsThreshold(state, policyStruct)
		state = UpdateFailedRequestsAccess(state, policyStruct)
		state = UpdateRouteListAndFlush(state, policyStruct)
		state = UpdateNetwork(state, policyStruct)

		fmt.Println("after Updates: ", time.Since(start))

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
			TimeStep:                state.TimeStep}
		stateArray = append(stateArray, curState)
		//PrintState(state)
	}
	PrintState(state)
	fmt.Print("end of main: ")
	end := time.Since(start)
	fmt.Println(end)
}

func PrintState(state State) {
	fmt.Println("SuccessfulFound: ", state.SuccessfulFound)
	fmt.Println("FailedRequestsThreshold: ", state.FailedRequestsThreshold)
	fmt.Println("FailedRequestsAccess: ", state.FailedRequestsAccess)
	fmt.Println("CacheHits:", state.CacheStruct.CacheHits)
	fmt.Println("TimeStep: ", state.TimeStep)
	fmt.Println("OriginatorIndex: ", state.OriginatorIndex)
	fmt.Println("PendingMap: ", state.PendingMap)
	fmt.Println("RerouteMap: ", state.RerouteMap)
	//fmt.Println("RouteLists: ", state.RouteLists)
	//fmt.Println("CacheMap: ", state.CacheStruct.CacheMap)
}
