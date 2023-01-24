package model

type PolicyTest struct {
	found           bool
	route           []int
	thresholdFailed []int
	originatorIndex int
	accessFailed    bool
	paymentList     []int
}
