package workers

import (
	"bufio"
	"fmt"
	"go-incentive-simulation/config"
	"go-incentive-simulation/model/parts/output"
	"go-incentive-simulation/model/parts/types"
	"os"
	"sync"
)

func OutputWorker(outputChan chan types.OutputStruct, wg *sync.WaitGroup) {
	defer wg.Done()
	var outputStruct types.OutputStruct
	counter := 0
	var meanRewardPerForward output.MeanRewardPerForward
	var avgNumberOfHops output.AvgNumberOfHops
	var fractions output.Fractions
	var rewardFairnessForStoringAction output.RewardFairnessForStoringAction
	var rewardFairnessForAllActions output.RewardFairnessForAllActions
	var rewardFairnessForForwardingAction output.RewardFairnessForForwardingActions
	var negativeIncome output.NegativeIncome
	negativeIncome.IncomeMap = make(map[int]int)

	//filePath := "./results/output.txt"
	//err := os.Remove(filePath)
	//if err != nil {
	//	fmt.Println("Could not remove the file", filePath)
	//}
	//file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	panic(err)
	//}
	//defer func(file *os.File) {
	//	err1 := file.Close()
	//	if err1 != nil {
	//		fmt.Println("Couldn't close the file with filepath: ", filePath)
	//	}
	//}(file)
	//
	//writer := bufio.NewWriter(file) // default writer size is 4096 bytes
	////writer = bufio.NewWriterSize(writer, 1048576) // 1MiB
	//defer func(writer *bufio.Writer) {
	//	err1 := writer.Flush()
	//	if err1 != nil {
	//		fmt.Println("Couldn't flush the remaining buffer in the writer for output")
	//	}
	//}(writer)

	if config.GetMeanRewardPerForward() {
		file, filePath := output.MakeMeanRewardPerForwardFile()
		defer func(file *os.File) {
			err1 := file.Close()
			if err1 != nil {
				fmt.Println("Couldn't close the file with filepath: ", filePath)
			}
		}(file)
		meanRewardPerForward.Writer = bufio.NewWriter(file)
		defer func(writer *bufio.Writer) {
			err1 := writer.Flush()
			if err1 != nil {
				fmt.Println("Couldn't flush the remaining buffer in the writer for output")
			}
		}(meanRewardPerForward.Writer)
	}
	if config.GetAverageNumberOfHops() {
		file2, filePath2 := output.MakeAvgNumberOfHopsFile()
		defer func(file2 *os.File) {
			err1 := file2.Close()
			if err1 != nil {
				fmt.Println("Couldn't close the file with filepath: ", filePath2)
			}
		}(file2)
		avgNumberOfHops.Writer = bufio.NewWriter(file2)
		defer func(writer *bufio.Writer) {
			err1 := writer.Flush()
			if err1 != nil {
				fmt.Println("Couldn't flush the remaining buffer in the writer for output")
			}
		}(avgNumberOfHops.Writer)
	}
	if config.GetAverageFractionOfTotalRewardsK16() {
		file3, filePath3 := output.MakeFractionOfRewardsFile()
		defer func(file3 *os.File) {
			err1 := file3.Close()
			if err1 != nil {
				fmt.Println("Couldn't close the file with filepath: ", filePath3)
			}
		}(file3)
		fractions.Writer = bufio.NewWriter(file3)
		defer func(writer *bufio.Writer) {
			err1 := writer.Flush()
			if err1 != nil {
				fmt.Println("Couldn't flush the remaining buffer in the writer for output")
			}
		}(fractions.Writer)
	}
	if config.GetRewardFairnessForStoringAction() {
		file4, filePath4 := output.MakeRewardFairnessForStoringActionFile()
		defer func(file4 *os.File) {
			err1 := file4.Close()
			if err1 != nil {
				fmt.Println("Couldn't close the file with filepath: ", filePath4)
			}
		}(file4)
		rewardFairnessForStoringAction.Writer = bufio.NewWriter(file4)
		defer func(writer *bufio.Writer) {
			err1 := writer.Flush()
			if err1 != nil {
				fmt.Println("Couldn't flush the remaining buffer in the writer for output")
			}
		}(rewardFairnessForStoringAction.Writer)
	}
	if config.GetRewardFairnessForAllActions() {
		file5, filePath5 := output.MakeRewardFairnessForAllActionsFile()
		defer func(file5 *os.File) {
			err1 := file5.Close()
			if err1 != nil {
				fmt.Println("Couldn't close the file with filepath: ", filePath5)
			}
		}(file5)
		rewardFairnessForAllActions.Writer = bufio.NewWriter(file5)
		defer func(writer *bufio.Writer) {
			err1 := writer.Flush()
			if err1 != nil {
				fmt.Println("Couldn't flush the remaining buffer in the writer for output")
			}
		}(rewardFairnessForAllActions.Writer)
	}
	if config.GetRewardFairnessForForwardingAction() {
		file6, filePath6 := output.MakeRewardFairnessForForwardingActionFile()
		defer func(file6 *os.File) {
			err1 := file6.Close()
			if err1 != nil {
				fmt.Println("Couldn't close the file with filepath: ", filePath6)
			}
		}(file6)
		rewardFairnessForForwardingAction.Writer = bufio.NewWriter(file6)
		defer func(writer *bufio.Writer) {
			err1 := writer.Flush()
			if err1 != nil {
				fmt.Println("Couldn't flush the remaining buffer in the writer for output")
			}
		}(rewardFairnessForForwardingAction.Writer)
	}
	if config.GetNegativeIncome() {
		file7, filePath7 := output.MakeNegativeIncomeFile()
		defer func(file7 *os.File) {
			err1 := file7.Close()
			if err1 != nil {
				fmt.Println("Couldn't close the file with filepath: ", filePath7)
			}
		}(file7)
		negativeIncome.Writer = bufio.NewWriter(file7)
		defer func(writer *bufio.Writer) {
			err1 := writer.Flush()
			if err1 != nil {
				fmt.Println("Couldn't flush the remaining buffer in the writer for output")
			}
		}(negativeIncome.Writer)
	}
	for outputStruct = range outputChan {
		counter++
		//fmt.Println("Route: ", outputStruct.RouteWithPrices)
		if config.GetMeanRewardPerForward() {
			for i := range outputStruct.RouteWithPrices {
				if i == len(outputStruct.RouteWithPrices)-1 {
					break
				}
				reward := outputStruct.RouteWithPrices[i].Price - outputStruct.RouteWithPrices[i+1].Price
				meanRewardPerForward.AllRewards = append(meanRewardPerForward.AllRewards, reward)
				meanRewardPerForward.SumRewards += reward
			}
			if counter%1_000_000 == 0 {
				mean := meanRewardPerForward.CalculateMeanRewardPerForward()
				_, err := meanRewardPerForward.Writer.WriteString(fmt.Sprintf("Mean reward per forward: %f \n", mean))
				if err != nil {
					panic(err)
				}
				meanRewardPerForward.AllRewards = nil
				meanRewardPerForward.SumRewards = 0
			}
		}

		if config.GetAverageNumberOfHops() {
			avgNumberOfHops.TotalNumberOfHops += len(outputStruct.RouteWithPrices)
			avgNumberOfHops.NumberOfRoutes++
			if counter%100_000 == 0 {
				hops := avgNumberOfHops.CalculateAverageNumberOfHops()
				_, err := avgNumberOfHops.Writer.WriteString(fmt.Sprintf("Average number of hops: %f \n", hops))
				if err != nil {
					panic(err)
				}
				avgNumberOfHops.TotalNumberOfHops = 0
				avgNumberOfHops.NumberOfRoutes = 0
			}
		}

		if config.GetAverageFractionOfTotalRewardsK16() && config.GetMaxProximityOrder() == 16 {
			var FractionOfRewardsK16 output.FractionOfRewardsK16
			if len(outputStruct.RouteWithPrices) == 2 {
				FractionOfRewardsK16.RouteRewards = append(FractionOfRewardsK16.RouteRewards, outputStruct.RouteWithPrices[0].Price-outputStruct.RouteWithPrices[1].Price)
				FractionOfRewardsK16.RouteRewards = append(FractionOfRewardsK16.RouteRewards, outputStruct.RouteWithPrices[1].Price)
				FractionOfRewardsK16.SumRouteRewards += outputStruct.RouteWithPrices[0].Price - outputStruct.RouteWithPrices[1].Price
				FractionOfRewardsK16.SumRouteRewards += outputStruct.RouteWithPrices[1].Price
				FractionOfRewardsK16.Hop1 = float64(FractionOfRewardsK16.RouteRewards[0]) / float64(FractionOfRewardsK16.SumRouteRewards)
				FractionOfRewardsK16.Hop2 = float64(FractionOfRewardsK16.RouteRewards[1]) / float64(FractionOfRewardsK16.SumRouteRewards)
				fractions.Fractions = append(fractions.Fractions, FractionOfRewardsK16)
				FractionOfRewardsK16.RouteRewards = nil
				FractionOfRewardsK16.SumRouteRewards = 0
			}
			if len(outputStruct.RouteWithPrices) == 3 {
				FractionOfRewardsK16.RouteRewards = append(FractionOfRewardsK16.RouteRewards, outputStruct.RouteWithPrices[0].Price-outputStruct.RouteWithPrices[1].Price)
				FractionOfRewardsK16.RouteRewards = append(FractionOfRewardsK16.RouteRewards, outputStruct.RouteWithPrices[1].Price-outputStruct.RouteWithPrices[2].Price)
				FractionOfRewardsK16.RouteRewards = append(FractionOfRewardsK16.RouteRewards, outputStruct.RouteWithPrices[2].Price)
				FractionOfRewardsK16.SumRouteRewards += outputStruct.RouteWithPrices[0].Price - outputStruct.RouteWithPrices[1].Price
				FractionOfRewardsK16.SumRouteRewards += outputStruct.RouteWithPrices[1].Price - outputStruct.RouteWithPrices[2].Price
				FractionOfRewardsK16.SumRouteRewards += outputStruct.RouteWithPrices[2].Price
				FractionOfRewardsK16.Hop1 = float64(FractionOfRewardsK16.RouteRewards[0]) / float64(FractionOfRewardsK16.SumRouteRewards)
				FractionOfRewardsK16.Hop2 = float64(FractionOfRewardsK16.RouteRewards[1]) / float64(FractionOfRewardsK16.SumRouteRewards)
				FractionOfRewardsK16.Hop3 = float64(FractionOfRewardsK16.RouteRewards[2]) / float64(FractionOfRewardsK16.SumRouteRewards)
				fractions.Fractions = append(fractions.Fractions, FractionOfRewardsK16)
				FractionOfRewardsK16.RouteRewards = nil
				FractionOfRewardsK16.SumRouteRewards = 0
			}
			if counter%100_000 == 0 {
				hop1, hop2, hop3 := fractions.CalculateFractionOfRewards()
				_, err := fractions.Writer.WriteString(fmt.Sprintf("hop 1: %f, hop 2: %f, hop 3: %f \n", hop1, hop2, hop3))
				if err != nil {
					panic(err)
				}
				fractions.Fractions = nil
			}
		}

		if config.GetRewardFairnessForStoringAction() {
			route := outputStruct.RouteWithPrices
			if route != nil {
				reward := route[len(route)-1].Price
				rewardFairnessForStoringAction.AllStoringRewards = append(rewardFairnessForStoringAction.AllStoringRewards, reward)
				rewardFairnessForStoringAction.SumAllStoringRewards += reward
			}
			if counter == config.GetIterations() {
				fairness := rewardFairnessForStoringAction.CalculateRewardFairnessForStoringAction()
				_, err := rewardFairnessForStoringAction.Writer.WriteString(fmt.Sprintf("Reward fairness for storing action: %f \n", fairness))
				if err != nil {
					panic(err)
				}
				rewardFairnessForStoringAction.AllStoringRewards = nil
				rewardFairnessForStoringAction.SumAllStoringRewards = 0
			}
		}

		if config.GetRewardFairnessForAllActions() {
			route := outputStruct.RouteWithPrices
			if route != nil {
				for i := range route {
					if i == len(route)-1 {
						reward := route[i].Price
						rewardFairnessForAllActions.AllRewards = append(rewardFairnessForAllActions.AllRewards, reward)
						rewardFairnessForAllActions.SumAllRewards += reward
						break
					}
					reward := route[i].Price - route[i+1].Price
					rewardFairnessForAllActions.AllRewards = append(rewardFairnessForAllActions.AllRewards, reward)
					rewardFairnessForAllActions.SumAllRewards += reward
				}
				if counter == 100_000 {
					fairness := rewardFairnessForAllActions.CalculateRewardFairnessForAllActions()
					_, err := rewardFairnessForAllActions.Writer.WriteString(fmt.Sprintf("Reward fairness for all actions: %f \n", fairness))
					if err != nil {
						panic(err)
					}
				}
			}
		}

		if config.GetRewardFairnessForForwardingAction() {
			route := outputStruct.RouteWithPrices
			if route != nil {
				for i := range route {
					if i == len(route)-1 {
						break
					}
					reward := route[i].Price - route[i+1].Price
					rewardFairnessForForwardingAction.AllForwardingRewards = append(rewardFairnessForForwardingAction.AllForwardingRewards, reward)
					rewardFairnessForForwardingAction.SumAllForwardingRewards += reward
				}
				if counter == config.GetIterations() {
					fairness := rewardFairnessForForwardingAction.CalculateRewardFairnessForForwardingAction()
					_, err := rewardFairnessForAllActions.Writer.WriteString(fmt.Sprintf("Reward fairness for forwarding action: %f \n", fairness))
					if err != nil {
						panic(err)
					}
				}
			}
		}

		// payment enabled, forgiveness enabled, threshold enabled, k = 8
		if config.GetNegativeIncome() && config.GetPaymentEnabled() && config.IsForgivenessEnabled() { //Payment enabled can use output.Payment
			route := outputStruct.RouteWithPrices
			if route != nil {
				for _, nodePair := range route {
					//if i == len(payments)-1 {
					//	break
					//}
					payer := int(nodePair.RequesterNode)
					payee := int(nodePair.ProviderNode)
					value := nodePair.Price
					valPayer, ok := negativeIncome.IncomeMap[payer]
					if !ok {
						negativeIncome.IncomeMap[payer] = 0
					}
					valPayee, ok := negativeIncome.IncomeMap[payee]
					if !ok {
						negativeIncome.IncomeMap[payee] = 0
					}
					negativeIncome.IncomeMap[payer] = valPayer - value
					negativeIncome.IncomeMap[payee] = valPayee + value
				}
			}
			// if counter%500_000==0 or counter==100_000 {
			if counter%100_000 == 0 {
				negativeIncomeRes := negativeIncome.CalculateNegativeIncome()
				_, err := negativeIncome.Writer.WriteString(fmt.Sprintf("Negative income: %f %% \n", negativeIncomeRes*100))
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
