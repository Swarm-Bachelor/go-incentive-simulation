#!/bin/bash

goProgram="main.go"
outputFilePath="output.txt"
configFilePath="config.yaml"

numRuns=5
endGr=30
startGr=10

rm "$outputFilePath"

for (( i=$startGr; i<=$endGr; i++ )); do
    echo "$i" >> "$outputFilePath"
    sed -i "s/NumGoroutines: .*/NumGoroutines: $i/" "$configFilePath"
    for (( j=1; j<=$numRuns; j++ )); do
        go run "$goProgram" >> "$outputFilePath"
    done
done
