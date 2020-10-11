package main

import (
	"qap"
	"utils"
	"os"
	"fmt"
	"time"
	"qap_solution"
)

func main() {
	// Take execution parameters
	data, pob_size, gen_num, tour_size := utils.TakeParameters(os.Args[1:])

	// Obtain QAP data
	fmt.Println("Initializing data....")
	qap_object := qap.CreateQAPObject(&data)

	// Execute QAP algorithm
	fmt.Println("Starting algorithm execution")
	start := time.Now()
	var best_solution qap_solution.QAPSolution = qap_object.ExecuteAlgorithm(pob_size, gen_num, tour_size)
	end := time.Now()

	// Show results
	fmt.Println("-------------------------------------------------")
	fmt.Println("FINAL RESULTS: ")
	best_solution.PrintGenerationResult(-1)
	fmt.Println("- Execution time: ", end.Sub(start))
}