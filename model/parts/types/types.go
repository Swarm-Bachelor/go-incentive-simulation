package types

type Request struct {
	TimeStep        int
	Epoch           int
	OriginatorIndex int
	OriginatorId    NodeId
	ChunkId         ChunkId
	RespNodes       [4]NodeId
}

type RequestResult struct {
	Route           []NodeId
	PaymentList     []Payment
	ChunkId         ChunkId
	Found           bool
	AccessFailed    bool
	ThresholdFailed bool
	FoundByCaching  bool
}

type Payment struct {
	FirstNodeId  NodeId
	PayNextId    NodeId
	ChunkId      ChunkId
	IsOriginator bool
}

func (p Payment) IsNil() bool {
	if p.PayNextId == 0 && p.FirstNodeId == 0 && p.ChunkId == 0 {
		return true
	} else {
		return false
	}
}

type Threshold [2]NodeId

type StateSubset struct {
	WaitingCounter          int64
	RetryCounter            int64
	CacheHits               int64
	ChunkId                 int
	OriginatorIndex         int64
	SuccessfulFound         int64
	FailedRequestsThreshold int64
	FailedRequestsAccess    int64
	TimeStep                int64
	Epoch                   int
}

type RouteData struct {
	Epoch           int      `json:"e"`
	Route           []NodeId `json:"r"`
	ChunkId         ChunkId  `json:"c"`
	Found           bool     `json:"f"`
	ThresholdFailed bool     `json:"t"`
	AccessFailed    bool     `json:"a"`
}

//type StateData struct {
//	TimeStep int         `json:"t"`
//	State    StateSubset `json:"s"`
//}

type State struct {
	Graph                   *Graph
	Originators             []NodeId
	NodesId                 []NodeId
	RouteLists              []RequestResult
	UniqueWaitingCounter    int64
	UniqueRetryCounter      int64
	CacheHits               int64
	OriginatorIndex         int64
	SuccessfulFound         int64
	FailedRequestsThreshold int64
	FailedRequestsAccess    int64
	TimeStep                int64
	Epoch                   int
}

func (s *State) GetOriginatorId(originatorIndex int) NodeId {
	return s.Originators[originatorIndex]
}

type NodePairWithPrice struct {
	RequesterNode NodeId
	ProviderNode  NodeId
	Price         int
}

type PaymentWithPrice struct {
	Payment Payment
	Price   int
}

type Output struct {
	RouteWithPrices    []NodePairWithPrice
	PaymentsWithPrices []PaymentWithPrice
}

type Outputs struct {
	Outputs []Output
}
