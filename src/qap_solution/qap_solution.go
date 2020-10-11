package qap_solution

import "fmt"

type QAPSolution struct {
	Solution []int
	Result int
}

func (s *QAPSolution) PrintGenerationResult(generation_num int) {
	if generation_num >= 0 {
		fmt.Println("GENERATION ", generation_num)
	}
	fmt.Println("- Solution: ", s.Solution)
	fmt.Println("- Result: ", s.Result)
}

func (s *QAPSolution) CreateChild() QAPSolution {
	var solution []int = make([]int, len(s.Solution))
	
	for i:= 0; i < len(s.Solution); i++ {
		solution[i] = -1
	}
	
	return QAPSolution{solution, 0}
}

func GiveBestSolution(arr *[]QAPSolution) QAPSolution {
	array := *arr
	var best_result QAPSolution = array[0]
	for i := 1; i < len(array); i++ {
		if best_result.Result > array[i].Result {
			best_result = array[i]
		}
	}
	return best_result
}