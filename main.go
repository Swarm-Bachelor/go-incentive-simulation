package main

import (
	"fmt"
	"go-incentive-simulation/config"
	"go-incentive-simulation/model/parts/types"
	"go-incentive-simulation/model/parts/workers"
	"go-incentive-simulation/model/state"
	"math"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	config.InitConfigs()
	network := fmt.Sprintf("./data/nodes_data_%d_10000txt", config.GetBinSize())
	//networkStart := time.Now()
	globalState := state.MakeInitialState(network)
	//networkElapsed := time.Since(networkStart)

	//iterations := config.GetIterations()
	//numTotalGoRoutines := config.GetNumGoroutines()
	numRoutingGoroutines := config.GetNumRoutingGoroutines()

	wgMain := &sync.WaitGroup{}
	wgOutput := &sync.WaitGroup{}
	requestChan := make(chan types.Request, numRoutingGoroutines)
	outputChan := make(chan types.OutputStruct, 100000)
	routeChan := make(chan types.RouteData, 100000)
	stateChan := make(chan types.StateSubset, 100000)
	pauseChan := make(chan bool, numRoutingGoroutines)
	continueChan := make(chan bool, numRoutingGoroutines)

	if config.IsWriteRoutesToFile() {
		wgOutput.Add(1)
		go workers.RouteFlushWorker(routeChan, wgOutput)
	}
	if config.IsWriteStatesToFile() {
		wgOutput.Add(1)
		go workers.StateFlushWorker(stateChan, wgOutput)
	}

	go workers.RequestWorker(pauseChan, continueChan, requestChan, &globalState, wgMain)
	wgMain.Add(1)

	if config.IsOutputEnabled() {
		go workers.OutputWorker(outputChan, wgOutput)
		wgOutput.Add(1)
	}

	for i := 0; i < numRoutingGoroutines; i++ {
		wgMain.Add(1)
		go workers.RoutingWorker(pauseChan, continueChan, requestChan, outputChan, routeChan, stateChan, &globalState, wgMain)
	}

	wgMain.Wait()
	close(outputChan)
	close(stateChan)
	close(routeChan)
	wgOutput.Wait()

	//fmt.Println("")
	//fmt.Println("end of main: ")
	elapsed := time.Since(start)
	new := strings.Trim(elapsed.String(), "s")
	fmt.Println(new)
	//fmt.Println("Number of Iterations: ", iterations)
	//fmt.Println("Number of Total Goroutines: ", numTotalGoRoutines)
	//fmt.Println("Number of Routing Goroutines: ", numRoutingGoroutines)
	// allReq, thresholdFails, requestsToBucketZero, rejectedBucketZero, rejectedFirstHop := ReadRoutes("routes.json")
	// fmt.Println("allReq: ", allReq)
	// fmt.Println("thresholdFails: ", thresholdFails)
	// fmt.Println("requestsToBucketZero: ", requestsToBucketZero)
	// fmt.Println("rejectedBucketZero: ", rejectedBucketZero)
	// fmt.Println("rejectedFirstHop: ", rejectedFirstHop)
	//PrintState(globalState)

	// TODO: Add this to another function in another file?
	// buf, err := ioutil.ReadFile("states.bin")
	// if err != nil {
	// 	panic(err)
	// }
	// stateSubsets := &protoGenerated.StateSubsets{}
	// err = proto.Unmarshal(buf, stateSubsets)
	// if err != nil {
	// 	panic(err)
	// }
	// // Access the subset field
	// count := 0
	// for _, subset := range stateSubsets.Subset {
	// 	count++
	// 	if count > 10 {
	// 		break
	// 	}
	// 	fmt.Printf("OriginatorIndex: %d\n", subset.OriginatorIndex)
	// 	fmt.Printf("PendingMap: %d\n", subset.PendingMap)
	// 	fmt.Printf("RerouteMap: %d\n", subset.RerouteMap)
	// 	fmt.Printf("CacheStruct: %d\n", subset.CacheStruct)
	// 	fmt.Printf("SuccessfulFound: %d\n", subset.SuccessfulFound)
	// 	fmt.Printf("FailedRequestsThreshold: %d\n", subset.FailedRequestsThreshold)
	// 	fmt.Printf("FailedRequestsAccess: %d\n", subset.FailedRequestsAccess)
	// 	fmt.Printf("TimeStep: %d\n", subset.TimeStep)
	// }
	// // read the binary protobuf message from the file
	// buf, err := ioutil.ReadFile("routes.bin")
	// if err != nil {
	// 	panic(err)
	// }

	// // unmarshal the binary protobuf message into a RouteData struct
	// routeData := &protoGenerated.RouteData{}
	// err = proto.Unmarshal(buf, routeData)
	// if err != nil {
	// 	panic(err)
	// }

	// // print the RouteData struct
	// fmt.Printf("TimeStep: %d\n", routeData.GetTimeStep())
	// count := 0
	// routedata := routeData.GetRoutes()
	// fmt.Println("length", len(routedata))
	// for _, route := range routeData.GetRoutes() {
	// 	if count == 10 {
	// 		break
	// 	}
	// 	fmt.Printf("Route: %v\n", route.GetWaypoints())
	// 	fmt.Printf("Length: %d\n", route.GetLength())
	// 	count++
	// }
}

func PrintState(state types.State) {
	total := float64(state.SuccessfulFound + state.FailedRequestsThreshold + state.FailedRequestsAccess)
	fmt.Println("SuccessfulFound: ", state.SuccessfulFound, "-->", math.Round(float64(state.SuccessfulFound)/total*1000000)/10000, "%")
	fmt.Println("ThresholdFail: ", state.FailedRequestsThreshold, "-->", math.Round(float64(state.FailedRequestsThreshold)/total*1000000)/10000, "%")
	fmt.Println("AccessFail: ", state.FailedRequestsAccess, "-->", math.Round(float64(state.FailedRequestsAccess)/total*1000000)/10000, "%")
	fmt.Println("TimeStep: ", state.TimeStep)
	fmt.Println("OriginatorIndex: ", state.OriginatorIndex)
	fmt.Println("CacheHits:", state.CacheHits)
	fmt.Println("UniqueRetryCounter: ", state.UniqueRetryCounter)
	fmt.Println("UniqueWaitingCounter: ", state.UniqueWaitingCounter)
	//fmt.Println("PendingMap: ", state.PendingStruct.PendingMap, state.PendingStruct.Counter)
	//fmt.Println("RerouteMap: ", state.RerouteStruct.RerouteMap)
	//fmt.Println("RouteLists: ", state.RouteLists)
}
