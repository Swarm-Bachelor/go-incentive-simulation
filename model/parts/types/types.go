package types

type Request struct {
	OriginatorId int
	ChunkId      int
}

type PendingMap map[int]int

type RerouteMap map[int][]int

type CacheMap map[*Node]map[int]int

type CacheStruct struct {
	CacheHits int
	CacheMap  CacheMap
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
	OriginatorIndex         int
	SuccessfulFound         int
	FailedRequestsThreshold int
	FailedRequestsAccess    int
	TimeStep                int
	PriceList               []int
}

type Policy struct {
	Found                bool
	Route                Route
	ThresholdFailedLists [][]Threshold
	AccessFailed         bool
	PaymentList          []Payment
}
