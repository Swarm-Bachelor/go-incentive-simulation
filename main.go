package main

import (
	"fmt"
	. "go-incentive-simulation/model/parts/policy"
	. "go-incentive-simulation/model/parts/types"
	. "go-incentive-simulation/model/parts/update"
	. "go-incentive-simulation/model/state"
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
	for i := 0; i < iterations; {
		policyStruct := PolicyStruct{}
		for j := 0; j < 5; j++ {
			policyOutput := MakePolicyOutput(state)

			policyStruct.Founds = append(policyStruct.Founds, policyOutput.Found)
			policyStruct.Routes = append(policyStruct.Routes, policyOutput.Route)
			policyStruct.ThresholdFailedListsList = append(policyStruct.ThresholdFailedListsList, policyOutput.ThresholdFailedLists)
			policyStruct.OriginatorIndices = append(policyStruct.OriginatorIndices, policyOutput.OriginatorIndex)
			policyStruct.AccessFails = append(policyStruct.AccessFails, policyOutput.AccessFailed)
			policyStruct.PaymentListList = append(policyStruct.PaymentListList, policyOutput.PaymentList)
			i++
		}
		
		state = UpdatePendingMap(state, policyStruct)
		state = UpdateRerouteMap(state, policyStruct)
		state = UpdateCacheMap(state, policyStruct)
		state = UpdateOriginatorIndex(state, policyStruct)
		state = UpdateSuccessfulFound(state, policyStruct)
		state = UpdateFailedRequestsThreshold(state, policyStruct)
		state = UpdateFailedRequestsAccess(state, policyStruct)
		state = UpdateRouteListAndFlush(state, policyStruct)
		state = UpdateNetwork(state, policyStruct)

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
