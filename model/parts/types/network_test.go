package types

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGenerateAndLoad(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	bits := 16
	bin := 4
	size := 10000
	network := &Network{Bits: bits, Bin: bin}
	nodes := network.Generate(size)

	filename := fmt.Sprintf("nodes_data_%d_%d.txt", bin, size)
	network.Dump(filename)

	network2 := Network{}
	bits2, bin2, nodes2 := network2.Load(filename)

	//Check if bits2, bin2, nodes2 are the same as bits, bin, nodes
	if bits2 != bits {
		t.Error("Bits are different")
	}
	if bin2 != bin {
		t.Error("Bin are different")
	}
	if len(nodes2) != len(nodes) {
		t.Error("Nodes are different")
	}

}
