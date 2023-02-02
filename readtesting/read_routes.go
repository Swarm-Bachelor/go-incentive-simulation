package readtesting

import (
	"encoding/json"
	"fmt"
	. "go-incentive-simulation/model/parts/types"
	"io/ioutil"
)

type RouteData struct {
	Timestep int     `json:"timestep"`
	Routes   []Route `json:"routes"`
}

func ReadRoutes(filePath string) (int, int, int, int, int) {
	fmt.Println("HERE")
	var data RouteData
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("error reading file: ", err)
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("error unmarshalling file: ", err)
	}
	fmt.Println("HERE2")
	routes := data.Routes
	allReq := 0
	rejectedBucketZero := 0
	rejectedFirstHop := 0
	requestsToBucketZero := 0
	thresholdFails := 0
	for _, route := range routes {
		fmt.Println("Inside loop")
		fmt.Println(route)
		allReq++
		first := route[0]
		chunk := route[len(route)-1]
		if contains(route, -3) {
			chunk = route[len(route)-2]
		}
		if bitLength(first^chunk) == 16 {
			requestsToBucketZero++
		}
		if contains(route, -1) {
			thresholdFails++
			if bitLength(first^chunk) == 16 && indexOf(route, -1) == 1 {
				rejectedBucketZero++
			}
			if indexOf(route, -1) == 1 {
				rejectedFirstHop++
			}
		}
	}

	fmt.Printf("all requests: %d\n", allReq)
	fmt.Printf("threshold_rejects: %d\n", thresholdFails)
	fmt.Printf("requests_routed_to_bucket_zero: %d\n", requestsToBucketZero)
	fmt.Printf("requests_rejected_in_bucket_zero: %d\n", rejectedBucketZero)
	fmt.Printf("requests_rejected_in_first_hop: %d\n", rejectedFirstHop)
	return allReq, thresholdFails, requestsToBucketZero, rejectedBucketZero, rejectedFirstHop
}

func bitLength(n int) int {
	count := 0
	for n != 0 {
		count++
		n >>= 1
	}
	return count
}

func contains(list []int, ele int) bool {
	for _, v := range list {
		if v == ele {
			return true
		}
	}
	return false
}

func indexOf(list []int, ele int) int {
	for i, v := range list {
		if v == ele {
			return i
		}
	}
	return -1
}
