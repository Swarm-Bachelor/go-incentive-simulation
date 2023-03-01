package workers

import (
	. "go-incentive-simulation/model/constants"
	. "go-incentive-simulation/model/parts/types"
	. "go-incentive-simulation/model/parts/update"
	"math/rand"
)

func RequestWorker(newStateChan chan bool, requestChan chan Request, globalState *State, iterations int) {

	//curState := globalState
	counter := 0
	for counter < iterations {
		if len(requestChan) < Constants.GetNumGoroutines() {
			UpdateOriginatorIndex(globalState)

			//curState = <-stateChan
			//<-newStateChan

			chunkId := rand.Intn(Constants.GetRangeAddress() - 1)

			if Constants.IsPreferredChunksEnabled() {
				var random float32
				numPreferredChunks := 1000
				random = rand.Float32()
				if float32(random) <= 0.5 {
					chunkId = rand.Intn(numPreferredChunks)
				} else {
					chunkId = rand.Intn(Constants.GetRangeAddress()-numPreferredChunks) + numPreferredChunks
				}
			}

			responsibleNodes := globalState.Graph.FindResponsibleNodes(chunkId)
			originatorId := globalState.Originators[globalState.OriginatorIndex]
			//originatorId := prevState.Originators[rand.Intn(Constants.GetOriginators())]

			if _, ok := globalState.PendingMap[originatorId]; ok {
				chunkId = globalState.PendingMap[originatorId]
				responsibleNodes = globalState.Graph.FindResponsibleNodes(chunkId)
			}

			if _, ok := globalState.RerouteMap[originatorId]; ok {
				chunkId = globalState.RerouteMap[originatorId][len(globalState.RerouteMap[originatorId])-1]
				responsibleNodes = globalState.Graph.FindResponsibleNodes(chunkId)
			}

			request := Request{
				OriginatorId: originatorId,
				ChunkId:      chunkId,
				RespNodes:    responsibleNodes,
			}
			requestChan <- request
			counter++
		}
	}
}