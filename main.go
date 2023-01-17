package main

import (
	"fmt"
	"os"
)

func main() {
	LoadNodes("/home/voldemort/bachelor/go-incentive-simulation/model/parts/utils/input_test.txt")

}

func LoadNodes(path string) {
	fileContent, _ := os.Open(path)
	defer fileContent.Close()

	fmt.Println(fileContent)



}
