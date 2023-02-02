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

func MakePolicyOutputOrig(state State) Policy {
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

func MakePolicyOutput(policyCh chan PolicyStruct, stateCh chan State, allDone chan struct{}) {
	fmt.Println("start of make initial policy")
	var state State
	var wg sync.WaitGroup
	for i := 0; i < iterations/numGoRoutines; i++ {
		state = <-stateCh
		policyStruct := PolicyStruct{
			Founds:                   make([]bool, numGoRoutines),
			Routes:                   make([]Route, numGoRoutines),
			ThresholdFailedListsList: make([][][]Threshold, numGoRoutines),
			OriginatorIndices:        make([]int, numGoRoutines),
			AccessFails:              make([]bool, numGoRoutines),
			PaymentListList:          make([][]Payment, numGoRoutines),
		}
		wg.Add(numGoRoutines)
		for j := 0; j < numGoRoutines; j++ {
			loop := j
			go func(int) {
				found, route, thresholdFailed, accessFailed, paymentsList := SendRequest(&state)
				policy := Policy{
					Found:                found,
					Route:                route,
					ThresholdFailedLists: thresholdFailed,
					OriginatorIndex:      state.OriginatorIndex,
					AccessFailed:         accessFailed,
					PaymentList:          paymentsList,
				}

				policyStruct.Founds[loop] = policy.Found
				policyStruct.Routes[loop] = policy.Route
				policyStruct.ThresholdFailedListsList[loop] = policy.ThresholdFailedLists
				policyStruct.OriginatorIndices[loop] = policy.OriginatorIndex
				policyStruct.AccessFails[loop] = policy.AccessFailed
				policyStruct.PaymentListList[loop] = policy.PaymentList

				wg.Done()
			}(loop)
		}
		wg.Wait()
		policyCh <- policyStruct
	}
	allDone <- struct{}{}
}

func UpdateState(policyCh chan PolicyStruct, stateCh chan State, firstState State) {
	state := firstState

	for policyStruct := range policyCh {

		state = UpdatePendingMap(state, policyStruct)
		state = UpdateRerouteMap(state, policyStruct)
		state = UpdateCacheMap(state, policyStruct)
		state = UpdateOriginatorIndex(state, policyStruct)
		state = UpdateSuccessfulFound(state, policyStruct)
		state = UpdateFailedRequestsThreshold(state, policyStruct)
		state = UpdateFailedRequestsAccess(state, policyStruct)
		state = UpdateRouteListAndFlush(state, policyStruct)
		state = UpdateNetwork(state, policyStruct)

		for i := 0; i < numGoRoutines; i++ {
			stateCh <- state
		}
	}
}

const iterations = 250000
const numGoRoutines = 10

func main() {

	start := time.Now()
	state := MakeInitialState("./data/nodes_data_8_10000.txt")
	//stateArray := []State{state}

	//wg := &sync.WaitGroup{}
	policyCh := make(chan PolicyStruct, numGoRoutines)
	stateCh := make(chan State, numGoRoutines)
	allDone := make(chan struct{})

	for i := 0; i < numGoRoutines; i++ {
		go MakePolicyOutput(policyCh, stateCh, allDone)
	}
	go UpdateState(policyCh, stateCh, state)

	stateCh <- state

	<-allDone

	state = <-stateCh
	//curState := State{
	//	Graph:                   state.Graph,
	//	Originators:             state.Originators,
	//	NodesId:                 state.NodesId,
	//	RouteLists:              state.RouteLists,
	//	PendingMap:              state.PendingMap,
	//	RerouteMap:              state.RerouteMap,
	//	CacheStruct:             state.CacheStruct,
	//	OriginatorIndex:         state.OriginatorIndex,
	//	SuccessfulFound:         state.SuccessfulFound,
	//	FailedRequestsThreshold: state.FailedRequestsThreshold,
	//	FailedRequestsAccess:    state.FailedRequestsAccess,
	//	TimeStep:                state.TimeStep}
	//stateArray = append(stateArray, curState)
	//PrintState(state)

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
