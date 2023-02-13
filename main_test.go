package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

type Data struct {
	Timestep int     `json:"timestep"`
	Routes   [][]int `json:"routes"`
}

func TestReadRoutes(t *testing.T) {
	err := CountHopsInRoutesFile("routes.json")
	if err != nil {
		t.Error(err)
	}
}

func CountHopsInRoutesFile(filename string) error {
	dataBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var data Data
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return err
	}

	var totalHops int
	for _, route := range data.Routes {
		hops := len(route) - 1
		if Contains(route, -3) {
			hops = len(route) - 2
		}
		totalHops += hops
		fmt.Printf("Route with %d hops: %v\n", hops, route)
	}

	avgHops := float64(totalHops) / float64(len(data.Routes))
	fmt.Printf("Average hops in routes: %f\n", avgHops)

	return nil
}

func Contains(slice []int, search int) bool {
	for _, s := range slice {
		if s == search {
			return true
		}
	}
	return false
}
