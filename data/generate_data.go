package main

import (
	"fmt"
	"go-incentive-simulation/config"
	"go-incentive-simulation/model/parts/types"
	"math/rand"
	"os"
	"time"
)

func main() {
	config.InitConfigs()
	start := time.Now()
	binSize := config.GetBinSize()
	bits := config.GetBits()
	networkSize := config.GetNetworkSize()
	rand.Seed(config.GetRandomSeed())
	network := types.Network{Bits: bits, Bin: binSize}
	filename := fmt.Sprintf("test_nodes_data_%d_%d.txt", binSize, networkSize)
	filepath := "./data/" + filename
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	err = network.GenerateConcurrently(networkSize, file)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Successfully generated new network file: ", filename)
	}
	fmt.Println("Time since start: ", time.Since(start))
}
