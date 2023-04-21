package state

import (
	"fmt"
	"go-incentive-simulation/config"
	"go-incentive-simulation/model/parts/types"
	"go-incentive-simulation/model/parts/utils"
	"math/rand"
)

func MakeInitialState(path string) types.State {
	// Initialize the state
	//fmt.Println("start of make initial state")
	rand.Seed(config.GetRandomSeed())
	network := types.Network{}
	network.Load(path)
	graph, err := utils.CreateGraphNetwork(&network)
	if err != nil {
		fmt.Println("create graph network returned an error: ", err)
	}
	//pendingStruct := types.PendingStruct{PendingMap: make(types.PendingMap, 0), PendingMutex: &sync.Mutex{}}
	//rerouteStruct := types.RerouteStruct{RerouteMap: make(types.RerouteMap, 0), RerouteMutex: &sync.Mutex{}}
	//cacheStruct := types.CacheStruct{CacheHits: 0, CacheMap: make(types.CacheMap), CacheMutex: &sync.Mutex{}}

	initialState := types.State{
		Graph:                   graph,
		Originators:             utils.CreateDownloadersList(graph),
		NodesId:                 utils.CreateNodesList(graph),
		RouteLists:              make([]types.RequestResult, 10000),
		CacheHits:               0,
		UniqueWaitingCounter:    0,
		UniqueRetryCounter:      0,
		OriginatorIndex:         0,
		SuccessfulFound:         0,
		FailedRequestsThreshold: 0,
		FailedRequestsAccess:    0,
		TimeStep:                0,
		Epoch:                   0,
	}
	return initialState
}
