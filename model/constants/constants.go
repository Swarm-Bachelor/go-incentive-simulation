package constants

import "math/rand"

type constant struct {
	runs        int
	bits        int
	networkSize int
	binSize     int
	rangeAddress                     int
	originators                      int
	refreshRate                      int
	threshold                        int
	randomSeed                       int
	maxProximityOrder                int
	price                            int
	chunks                           int
	requestsPerSecond                int
	thresholdEnabled                 bool
	paymentEnabled                   bool
	maxPoCheckEnabled                bool
	waitingEnabled                   bool
	onlyOriginatorPays               bool
	payOnlyForCurrentRequest         bool
	payIfOrigPays                    bool
	forwarderPayForceOriginatorToPay bool
	retryWithAnotherPeer             bool
	cacheIsEnabled                   bool
}

var Constants = constant{
	runs:        1,
	bits:        16,
	networkSize: 1000,
	binSize:     8,
	rangeAddress: 0, // 2 * *Bits
	originators : 1, // int(0.001 * NetworkSize)
	refreshRate: 8,
	threshold: 16,
	randomSeed : rand.Int() * 100,
	maxProximityOrder: 16,
	price: 1,
	chunks: 1000,
	requestsPerSecond: 12500,
	thresholdEnabled: true,
	paymentEnabled: true,
	maxPoCheckEnabled: false,
	waitingEnabled: false,
	onlyOriginatorPays: false,
	payOnlyForCurrentRequest: false,
	payIfOrigPays: false,
	forwarderPayForceOriginatorToPay: false,
	retryWithAnotherPeer: false,
	cacheIsEnabled: true,
}

// func CreateRangeAddress(c *constant){
// 	c.rangeAddress = 2 * c.bits
// }

// func (c *constant) CreateOriginators(){
// 	c.originators = int(0.001 * float64(c.networkSize))
// }

func (c *constant) IsCacheEnabled() bool {
	return c.cacheIsEnabled
}

func (c *constant) IsRetryWithAnotherPeer() bool {
	return c.retryWithAnotherPeer
}

func (c *constant) IsForwarderPayForceOriginatorToPay() bool {
	return c.forwarderPayForceOriginatorToPay
}

func (c *constant) IsPayIfOrigPays() bool {
	return c.payIfOrigPays
}

func (c *constant) IsPayOnlyForCurrentRequest() bool {
	return c.payOnlyForCurrentRequest
}

func (c *constant) IsOnlyOriginatorPays() bool {
	return c.onlyOriginatorPays
}

func (c *constant) IsWaitingEnabled() bool {
	return c.waitingEnabled
}

func (c *constant) GetMaxPoCheckEnabled() bool {
	return c.maxPoCheckEnabled
}

func (c *constant) GetThresholdEnabled() bool {
	return c.thresholdEnabled
}

func (c *constant) GetPaymentEnabled() bool {
	return c.paymentEnabled
}

func (c *constant) GetRequestsPerSecond() int {
	return c.requestsPerSecond
}

func (c *constant) GetChunks() int {
	return c.chunks
}

func (c *constant) GetBits() int {
	return c.bits
}

func (c *constant) GetNetworkSize() int {
	return c.networkSize
}

func (c *constant) GetBinSize() int {
	return c.binSize
}

func GetSimulationRuns() int {
	return 125000
}

func (c *constant) GetRangeAddress() int {
	return c.rangeAddress
}

func (c *constant) GetOriginators() int {
	return c.originators
}

func (c *constant) GetRefreshRate() int {
	return c.refreshRate
}

func (c *constant) GetThreshold() int {
	return c.threshold
}

func (c *constant) GetRandomSeed() int {
	return c.randomSeed
}

func (c *constant) GetMaxProximityOrder() int {
	return c.maxProximityOrder
}

func (c *constant) GetPrice() int {
	return c.price
}

// np.random.seed(get_random_seed())



