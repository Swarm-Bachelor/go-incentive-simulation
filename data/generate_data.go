package main

import (
	"fmt"
	"go-incentive-simulation/config"
	"go-incentive-simulation/model/parts/types"
	"math/rand"
	"time"
)

func main() {
	start := time.Now()
	binSize := 16        //config.GetBinSize()
	bits := 32           //config.GetBits()
	networkSize := 40000 //config.GetNetworkSize()
	rand.Seed(config.GetRandomSeed())
	network := types.Network{Bits: bits, Bin: binSize}
	network.Generate(networkSize)
	filename := fmt.Sprintf("nodes_data_%d_%d_32bits.txt", binSize, networkSize)
	err := network.Dump(filename)
	if err != nil {
		return
	}
	elapsed := time.Since(start)
	fmt.Printf("Time elapsed: %s\n", elapsed)
}
