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

func ProduceSendRequest(policyCh chan Policy, stateCh chan State, doneDoneCh chan struct{}) {
	for true {
		select {
		case <-doneDoneCh:
			return
		case state := <-stateCh:
			found, route, thresholdFailed, accessFailed, paymentsList := SendRequest(&state)
			policy := Policy{
				Found:                found,
				Route:                route,
				ThresholdFailedLists: thresholdFailed,
				OriginatorIndex:      state.OriginatorIndex,
				AccessFailed:         accessFailed,
				PaymentList:          paymentsList,
			}
			policyCh <- policy
		default:
		}
	}
}

func MakePolicyOutput(policyStructCh chan PolicyStruct, policyCh chan Policy, allDone chan struct{}) {
	fmt.Println("start of make initial policy")
	var policy Policy

	for i := 0; i < iterations/numGoRoutines; i++ {
		policyStruct := PolicyStruct{
			Founds:                   make([]bool, numGoRoutines),
			Routes:                   make([]Route, numGoRoutines),
			ThresholdFailedListsList: make([][][]Threshold, numGoRoutines),
			OriginatorIndices:        make([]int, numGoRoutines),
			AccessFails:              make([]bool, numGoRoutines),
			PaymentListList:          make([][]Payment, numGoRoutines),
		}
		loop := 0
		for policy = range policyCh {

			policyStruct.Founds[loop] = policy.Found
			policyStruct.Routes[loop] = policy.Route
			policyStruct.ThresholdFailedListsList[loop] = policy.ThresholdFailedLists
			policyStruct.OriginatorIndices[loop] = policy.OriginatorIndex
			policyStruct.AccessFails[loop] = policy.AccessFailed
			policyStruct.PaymentListList[loop] = policy.PaymentList
			loop++
			if loop == numGoRoutines {
				break
			}
		}
		policyStructCh <- policyStruct
	}
	allDone <- struct{}{}
}

func UpdateState(policyStructCh chan PolicyStruct, stateCh chan State, firstState State) {
	state := firstState

	for policyStruct := range policyStructCh {

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
const numGoRoutines = 50

func main() {

	start := time.Now()
	state := MakeInitialState("./data/nodes_data_8_10000.txt")
	//stateArray := []State{state}

	//wg := &sync.WaitGroup{}
	policyStructCh := make(chan PolicyStruct, numGoRoutines)
	policyCh := make(chan Policy, numGoRoutines)
	stateCh := make(chan State, numGoRoutines)
	allDoneCh := make(chan struct{})
	doneDoneCh := make(chan struct{})

	go MakePolicyOutput(policyStructCh, policyCh, allDoneCh)
	go UpdateState(policyStructCh, stateCh, state)

	for i := 0; i < numGoRoutines; i++ {
		go ProduceSendRequest(policyCh, stateCh, doneDoneCh)
	}
	for i := 0; i < numGoRoutines; i++ {
		stateCh <- state
	}

	<-allDoneCh

	for i := 0; i < numGoRoutines; i++ {
		doneDoneCh <- struct{}{}
	}

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
