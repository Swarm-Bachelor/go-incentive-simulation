package utils

import (
	"go-incentive-simulation/config"
	"go-incentive-simulation/model/general"
	"go-incentive-simulation/model/parts/types"
)

// returns the next node in the route, which is the closest node to the route in the previous nodes adjacency list
func getNext(request types.Request, firstNodeId types.NodeId, prevNodePaid bool, graph *types.Graph) (types.NodeId, bool, bool, bool, types.Payment) {
	var nextNodeId types.NodeId = 0
	var payNextId types.NodeId = 0
	var payment types.Payment
	var thresholdFailed bool
	var accessFailed bool
	mainOriginatorId := request.OriginatorId
	chunkId := request.ChunkId
	lastDistance := firstNodeId.ToInt() ^ chunkId.ToInt()
	//fmt.Printf("\n last distance is : %d, chunk is: %d, first is: %d", lastDistance, chunkId, firstNodeId)
	//fmt.Printf("\n which bucket: %d \n", 16-BitLength(chunkId^firstNodeId))

	currDist := lastDistance
	payDist := lastDistance

	//var lockedEdges []types.NodeId

	bin := config.GetBits() - general.BitLength(firstNodeId.ToInt()^chunkId.ToInt())
	firstNodeAdjIds := graph.GetNodeAdj(firstNodeId)

	for _, nodeId := range firstNodeAdjIds[bin] {
		dist := nodeId.ToInt() ^ chunkId.ToInt()
		if general.BitLength(dist) >= general.BitLength(lastDistance) {
			continue
		}
		if dist >= currDist {
			continue
		}
		// This means the node is now actively trying to communicate with the other node
		if config.IsEdgeLock() {
			graph.LockEdge(firstNodeId, nodeId)
		}
		
		if !isThresholdFailed(firstNodeId, nodeId, graph, request) {
			thresholdFailed = false
			if config.IsRetryWithAnotherPeer() {
				rerouteStruct := graph.GetNode(mainOriginatorId).RerouteStruct
				if rerouteStruct.Reroute.RejectedNodes != nil {
					if general.Contains(rerouteStruct.Reroute.RejectedNodes, nodeId) {
						if config.IsEdgeLock() {
							graph.UnlockEdge(firstNodeId, nodeId)
						}
						continue // skips node that's been part of a failed route before
					}
				}
			}

			if config.IsEdgeLock() {
				if !nextNodeId.IsNil() {
					graph.UnlockEdge(firstNodeId, nextNodeId)
				}
				if !payNextId.IsNil() {
					graph.UnlockEdge(firstNodeId, payNextId)
					payNextId = 0 // IMPORTANT!
				}
			}
			currDist = dist
			nextNodeId = nodeId

		} else {
			thresholdFailed = true
			if config.GetPaymentEnabled() {
				if dist < payDist && nextNodeId.IsNil() {
					if config.IsEdgeLock() && !payNextId.IsNil() {
						graph.UnlockEdge(firstNodeId, payNextId)
					}
					payDist = dist
					payNextId = nodeId
				} else {
					if config.IsEdgeLock() {
						graph.UnlockEdge(firstNodeId, nodeId)
					}
				}
			} else {
				if config.IsEdgeLock() {
					graph.UnlockEdge(firstNodeId, nodeId)
				}
			}
		}
	}

	// unlocks all nodes except the nextNodeId lock
	//if constants.GetEdgeLock() {
	//	for _, nodeId := range lockedEdges {
	//		if nodeId.ToInt32() != nextNodeId.ToInt32() {
	//			graph.UnlockEdge(firstNodeId, nodeId)
	//			unlockedEdges = append(unlockedEdges, nodeId)
	//		}
	//	}
	//}

	if !nextNodeId.IsNil() {
		thresholdFailed = false
		accessFailed = false
	} else if !thresholdFailed {
		accessFailed = true
	} else {
	}

	if config.GetPaymentEnabled() && !payNextId.IsNil() {
		accessFailed = false

		if firstNodeId == mainOriginatorId {
			payment.IsOriginator = true
		}

		if config.IsOnlyOriginatorPays() {
			// Only set payment if the firstNode is the MainOriginator
			if payment.IsOriginator {
				payment.FirstNodeId = firstNodeId
				payment.PayNextId = payNextId
				payment.ChunkId = chunkId
				nextNodeId = payNextId
			}

		} else if config.IsPayIfOrigPays() {
			// Pay if the originator pays or if the previous node has paid
			if payment.IsOriginator || prevNodePaid {
				payment.FirstNodeId = firstNodeId
				payment.PayNextId = payNextId
				payment.ChunkId = chunkId
				nextNodeId = payNextId
				thresholdFailed = false

			} else {
				thresholdFailed = true
			}

		} else {
			// Always pays
			payment.FirstNodeId = firstNodeId
			payment.PayNextId = payNextId
			payment.ChunkId = chunkId
			nextNodeId = payNextId
			thresholdFailed = false
		}
	}

	if !payment.IsNil() {
		prevNodePaid = true
	} else {
		prevNodePaid = false
	}

	if !nextNodeId.IsNil() {
		// fmt.Println("Next node is: ", nextNodeId)
	}

	return nextNodeId, thresholdFailed, accessFailed, prevNodePaid, payment
}
