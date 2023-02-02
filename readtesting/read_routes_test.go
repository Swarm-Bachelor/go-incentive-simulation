package readtesting

import (
	"testing"
)

func TestReadRoutes(t *testing.T) {
	allReq, thresholdFails, requestsToBucketZero, rejectedBucketZero, rejectedFirstHop := ReadRoutes("routes_test.json")
	t.Log("all requests: ", allReq)
	t.Log("threshold_rejects: ", thresholdFails)
	t.Log("requests_routed_to_bucket_zero: ", requestsToBucketZero)
	t.Log("requests_rejected_in_bucket_zero: ", rejectedBucketZero)
	t.Log("requests_rejected_in_first_hop: ", rejectedFirstHop)
}
