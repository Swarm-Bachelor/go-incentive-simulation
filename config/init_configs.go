package config

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Variables This is the one that gets changed in setup
var Variables = GetDefaultVariables()

func InitConfigs() {
	ymlData := ReadYamlFile()
	SetConfOptions(ymlData.ConfOptions)
	SetExperiment(ymlData)
}

func ReadYamlFile() Yml {
	yamlFile, err := os.ReadFile("config.yaml")

	var yamlData Yml

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &yamlData)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return yamlData
}

func SetExperiment(yml Yml) {

	switch yml.Experiment.Name {
	case "omega":
		fmt.Println("omega experiment chosen")
		OmegaExperiment()

		// TODO: Add more experiments here

	case "custom":
		//fmt.Println("custom experiment chosen")
		CustomExperiment(yml.CustomExperiment)

	default:
		fmt.Println("default experiment chosen")
	}
}

func SetConfOptions(configOptions confOptions) {
	Variables.confOptions.Iterations = configOptions.Iterations
	Variables.confOptions.Bits = configOptions.Bits
	Variables.confOptions.NetworkSize = configOptions.NetworkSize
	Variables.confOptions.BinSize = configOptions.BinSize
	Variables.confOptions.RangeAddress = configOptions.RangeAddress
	Variables.confOptions.Originators = configOptions.Originators
	Variables.confOptions.RefreshRate = configOptions.RefreshRate
	Variables.confOptions.Threshold = configOptions.Threshold
	Variables.confOptions.RandomSeed = configOptions.RandomSeed
	Variables.confOptions.MaxProximityOrder = configOptions.MaxProximityOrder
	Variables.confOptions.Price = configOptions.Price
	Variables.confOptions.RequestsPerSecond = configOptions.RequestsPerSecond
	Variables.confOptions.EdgeLock = configOptions.EdgeLock
	Variables.confOptions.SameOriginator = configOptions.SameOriginator
	Variables.confOptions.PrecomputeRespNodes = configOptions.PrecomputeRespNodes
	Variables.confOptions.WriteRoutesToFile = configOptions.WriteRoutesToFile
	Variables.confOptions.WriteStatesToFile = configOptions.WriteStatesToFile
	Variables.confOptions.IterationMeansUniqueChunk = configOptions.IterationMeansUniqueChunk
	Variables.confOptions.DebugPrints = configOptions.DebugPrints
	Variables.confOptions.DebugInterval = configOptions.DebugInterval
	Variables.confOptions.OutputEnabled = configOptions.OutputEnabled

	SetNumGoroutines(configOptions.NumGoroutines)

	Variables.confOptions.OutputOptions.MeanRewardPerForward = configOptions.OutputOptions.MeanRewardPerForward
	Variables.confOptions.OutputOptions.AverageNumberOfHops = configOptions.OutputOptions.AverageNumberOfHops
	Variables.confOptions.OutputOptions.AverageFractionOfTotalRewardsK8 = configOptions.OutputOptions.AverageFractionOfTotalRewardsK8
	Variables.confOptions.OutputOptions.AverageFractionOfTotalRewardsK16 = configOptions.OutputOptions.AverageFractionOfTotalRewardsK16
	Variables.confOptions.OutputOptions.RewardFairnessForForwardingAction = configOptions.OutputOptions.RewardFairnessForForwardingAction
	Variables.confOptions.OutputOptions.RewardFairnessForStoringAction = configOptions.OutputOptions.RewardFairnessForStoringAction
	Variables.confOptions.OutputOptions.RewardFairnessForAllActions = configOptions.OutputOptions.RewardFairnessForAllActions
	Variables.confOptions.OutputOptions.NegativeIncome = configOptions.OutputOptions.NegativeIncome

}

func SetNumGoroutines(numGoroutines int) {
	if numGoroutines == -1 {
		Variables.confOptions.NumGoroutines = runtime.NumCPU()
	} else {
		Variables.confOptions.NumGoroutines = numGoroutines
	}
}
