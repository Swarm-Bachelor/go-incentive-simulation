package types

import (
	"sync"
)

type Request struct {
	OriginatorId int
	ChunkId      int
}

type PendingMap map[int]int

type RerouteMap map[int][]int

type CacheMap map[int]map[int]int

type CacheStruct struct {
	CacheHits  int
	CacheMap   CacheMap
	CacheMutex *sync.Mutex
}

// TODO: cache is now slower on than off because of the concurrency

func (c *CacheStruct) Contains(nodeId int, chunkId int) bool {
	c.CacheMutex.Lock()
	defer c.CacheMutex.Unlock()
	nodeMap := c.CacheMap[nodeId]
	if nodeMap != nil && nodeMap[chunkId] > 0 {
		return true
	}
	return false
}

func (c *CacheStruct) AddToCache(nodeId int, chunkId int) {
	c.CacheMutex.Lock()
	defer c.CacheMutex.Unlock()
	nodeMap := c.CacheMap[nodeId]
	if nodeMap != nil {
		if _, ok2 := nodeMap[chunkId]; ok2 {
			nodeMap[chunkId]++
		} else {
			nodeMap[chunkId] = 1
		}
	} else {
		c.CacheMap[nodeId] = map[int]int{chunkId: 1}
	}
	return
}

type Route []int

type Payment struct {
	FirstNodeId  int
	PayNextId    int
	ChunkId      int
	IsOriginator bool
}

type Threshold [2]int

type State struct {
	Graph                   *Graph
	Originators             []int
	NodesId                 []int
	RouteLists              []Route
	PendingMap              PendingMap
	RerouteMap              RerouteMap
	CacheStruct             CacheStruct
	OriginatorIndex         int32
	SuccessfulFound         int32
	FailedRequestsThreshold int32
	FailedRequestsAccess    int32
	TimeStep                int32
}

type Policy struct {
	Found                bool
	Route                Route
	ThresholdFailedLists [][]Threshold
	AccessFailed         bool
	PaymentList          []Payment
}
