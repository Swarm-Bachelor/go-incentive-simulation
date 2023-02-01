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

func MakePolicyOutput(wg *sync.WaitGroup, mutex *sync.Mutex, stateStruct *StateStruct, policyCh chan Policy, StateCh chan State) {
	fmt.Println("start of make initial policy")
	defer wg.Done()
	//state := stateArray[0]
	for i := 0; i < iterations; i++ {
		state := <-StateCh
		//mutex.Lock()
		//state = stateArray[len(stateArray)-1]
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
		//mutex.Unlock()

		found, route, thresholdFailed, accessFailed, paymentsList := SendRequest(&curState)

		policy := Policy{
			Found:                found,
			Route:                route,
			ThresholdFailedLists: thresholdFailed,
			OriginatorIndex:      state.OriginatorIndex,
			AccessFailed:         accessFailed,
			PaymentList:          paymentsList,
		}
		policyCh <- policy
	}
	close(policyCh)

}

func UpdateState(wg *sync.WaitGroup, mutex *sync.Mutex, statestruct *StateStruct, policyCh chan Policy, stateCh chan State) {
	defer wg.Done()
	state := statestruct.curState

	for policyOutput := range policyCh {
		//mutex.Lock()
		state = statestruct.curState
		newState := State{
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
		//mutex.Unlock()

		newState = UpdatePendingMap(newState, policyOutput)
		newState = UpdateRerouteMap(newState, policyOutput)
		newState = UpdateCacheMap(newState, policyOutput)
		newState = UpdateOriginatorIndex(newState, policyOutput)
		newState = UpdateSuccessfulFound(newState, policyOutput)
		newState = UpdateFailedRequestsThreshold(newState, policyOutput)
		newState = UpdateFailedRequestsAccess(newState, policyOutput)
		newState = UpdateRouteListAndFlush(newState, policyOutput)
		newState = UpdateNetwork(newState, policyOutput)

		stateCh <- newState

		//mutex.Lock()
		statestruct.curState = newState
		statestruct.stateArr = append(statestruct.stateArr, newState)
		//mutex.Unlock()
	}
}

const iterations = 250000
const numWorkers = 10

type StateStruct = struct {
	stateArr []State
	curState State
}

func main() {
	start := time.Now()
	state := MakeInitialState("./data/nodes_data_8_10000.txt")
	stateArray := []State{state}
	stateStruct := &StateStruct{
		stateArray,
		state,
	}

	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	policyCh := make(chan Policy, numWorkers)
	stateCh := make(chan State, numWorkers)

	wg.Add(2)
	go MakePolicyOutput(wg, mutex, stateStruct, policyCh, stateCh)
	go UpdateState(wg, mutex, stateStruct, policyCh, stateCh)
	stateCh <- state
	wg.Wait()

	PrintState(stateStruct.curState)
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
