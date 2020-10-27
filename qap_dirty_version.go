package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"strconv"
	"time"
	"sort"
	"math/rand"
	//"sync"
	//"runtime"
)

type QAPSolution struct {
	solution []int
	result int
}

type QAP struct {
	size int
	flows [][]int
	distances [][]int
	poblation []QAPSolution
}

func IsIn(elem int, arr *[]int) int {
	a := *arr
	for i := 0; i < len(a); i++ {
		if a[i] == elem {
			return i
		}
	}
	return -1
}

func readFile(filename string) []byte {
	fmt.Println("Reading file: ", filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file")
		os.Exit(0)
	}
	return data
}

func getNumber(elem string, args []string, arg_pos int) int {
	var number int = 0

	if len(args) > arg_pos {
		number, _ = strconv.Atoi(args[arg_pos])
	}

	return number
}

func takeParameters(args []string) ([]byte, int, int, int){
	// 1 - Filename
	data := readFile(args[0])

	// 2 - Poblation size
	pob_size := getNumber("poblation size", args, 1)
	fmt.Println("Poblation size: ", pob_size)

	// 3 - Number of generations
	gen_num := getNumber("the number of generations", args, 2)
	fmt.Println("Number of generations: ", gen_num)

	// 4 - Tournament size
	tour_size := getNumber("tournament size", args, 3)
	fmt.Println("Tournament size: ", tour_size)

	return data, pob_size, gen_num, tour_size
}

func obtainFileData(data *[]byte) []string {
	var text string = string(*data)
	return strings.Fields(text)
}

func obtainMatrixSize(number string) int {
	var size int
	size, err := strconv.Atoi(number)

	if err != nil {
		fmt.Println("We expect that matrix size is ubicated in first position")
		return 0
	} else {
		return size
	}
}

func obtainMatrix(matrix *[][]int, size int, subfileData []string) {
	m := make([][]int, size)

	for i := 0; i < size; i++ {
		m[i] = make([]int, size)
		for j := 0; j < size; j++ {
			m[i][j], _ = strconv.Atoi(subfileData[i*size + j])
		}
	}

	*matrix = m
}

func createQAPObject(data *[]byte) QAP {
	var fileData []string = obtainFileData(data)

	// QAP object
	var size int = obtainMatrixSize(fileData[0])
	if size == 0 {
		return QAP{}
	} else {
		qap := QAP{size: size}
		obtainMatrix(&qap.flows, size, fileData[1:(size*size)+1])
		obtainMatrix(&qap.distances, size, fileData[(size*size)+1:])
		return qap
	}
}

func giveBestSolution(arr *[]QAPSolution) QAPSolution {
	array := *arr
	var best_result QAPSolution = array[0]
	for i := 1; i < len(array); i++ {
		if best_result.result > array[i].result {
			best_result = array[i]
		}
	}
	return best_result
}

func obtainUniqueRandomNumbers(max int) (int, int) {
	var a, b int

	a = rand.Intn(max)
	b = rand.Intn(max)

	for a == b {
		b = rand.Intn(max)
	}
	
	if a > b {
		aux := a
		a = b
		b = aux
	}

	return a, b
}

func (s *QAPSolution) printGenerationResult(generation_num int) {
	if generation_num >= 0 {
		fmt.Println("GENERATION ", generation_num)
	}
	fmt.Println("- Solution: ", s.solution)
	fmt.Println("- Result: ", s.result)
}

func (s *QAPSolution) createChild() QAPSolution {
	var solution []int = make([]int, len(s.solution))
	
	for i:= 0; i < len(s.solution); i++ {
		solution[i] = -1
	}
	
	return QAPSolution{solution, 0}
}

/*
TARDA 4min 53 s
func (qap *QAP) calculateResult(qs *QAPSolution) { 
	s := *qs

	var wg sync.WaitGroup
	wg.Add(6)
	res := make(chan int)

	for i := 0; i < 6; i++ {
		go func (i int, r chan int) {
			sum := 0
			start := (qap.size / 6) * i
			end := start + (qap.size / 6)
			for j := start; j < end; j ++ {
				for k := 0; k < qap.size; k++ {
					s.result = s.result + (qap.flows[j][k] * qap.distances[s.solution[j]][s.solution[k]])
				}
			}

			res <- sum
			wg.Done()
		}(i, res)
		s.result = s.result + <-res
	}
	wg.Wait()
	*qs = s
}*/

/*
TARDA 4min
func (qap *QAP) calculateResult(qs *QAPSolution) {
	s := *qs

	//n := runtime.GOMAXPROCS(0)
	res := make(chan int)

	for i := 0; i < 6; i++ {
		go func (i int, r chan int) {
			sum := 0
			start := (qap.size / 6) * i
			end := start + (qap.size / 6)
			for j := start; j < end; j ++ {
				for k := 0; k < qap.size; k++ {
					s.result = s.result + (qap.flows[j][k] * qap.distances[s.solution[j]][s.solution[k]])
				}
			}

			res <- sum
		}(i, res)
		s.result = s.result + <-res
	}

	*qs = s
}*/

func (qap *QAP) calculateResult(qs *QAPSolution) {
	s := *qs
	for i := 0; i < qap.size; i++ {
		for j := 0; j < qap.size; j++ {
			s.result = s.result + (qap.flows[i][j] * qap.distances[s.solution[i]][s.solution[j]])
		}
	}
	*qs = s
}

func (qap *QAP) makeChilds(ind1 *QAPSolution, ind2 *QAPSolution) (QAPSolution, QAPSolution) {
	// Create references
	elem1 := *ind2
	elem2 := *ind1

	// Create childs
	var child1 QAPSolution = elem1.createChild()
	var child2 QAPSolution = elem2.createChild()

	var start, end int = obtainUniqueRandomNumbers(qap.size)
	var pos1, pos2 int

	// Take data from first parent
	for i := start; i < end; i++ {
		child1.solution[i] = elem1.solution[i]
		child2.solution[i] = elem2.solution[i]
	}

	if start == 0 {
		pos1, pos2 = end, end
	}

	// Take data from second parent
	for i := 0; i < qap.size; i++ {
		if IsIn(elem1.solution[i], &child2.solution) == -1 {
			child2.solution[pos1] = elem1.solution[i]
			pos1 = pos1 + 1
			if pos1 == start {
				pos1 = end
			}
		}

		if IsIn(elem2.solution[i], &child1.solution) == -1 {
			child1.solution[pos2] = elem2.solution[i]
			pos2 = pos2 + 1
			if pos2 == start {
				pos2 = end
			}
		}
	}

	// Calculate result
	qap.calculateResult(&child1)
	qap.calculateResult(&child2)

	return child1, child2
}

/*func (s *QAPSolution) searchBestSolution(qap QAP) {
	qs := s

	var wg sync.WaitGroup
	wg.Add(6)

	res := make(chan QAPSolution)
	res <- *qs

	for k := 0; k < 6; k++ {
		go func (k int, r chan QAPSolution) {
			start := (len(s.solution) / 6) * k
			end := start + (len(s.solution) / 6)

			qp := <-r

			for i := start; i < end; i++ {
				for j := (i+1); j < len(qp.solution); j++ {
					// Do swap
					aux := qp.solution[j]
					qp.solution[j] = qp.solution[i]
					qp.solution[i] = aux

					qap.calculateResult(&*qs)

					if s.result > qp.result {
						s = &qp
					}

					// Undo swap
					aux = qs.solution[j]
					qs.solution[j] = qs.solution[i]
					qs.solution[i] = aux
				}
			}

			wg.Done()
		}(k, res)
	}

	wg.Wait()
}*/

func (qap *QAP) obtainInitialPoblation(pob_size int) {
	qap.poblation = make([]QAPSolution, pob_size)
	
	for i := 0; i < pob_size; i++ {
		s := QAPSolution{}

		// Initialize a random solution
		s.solution = rand.Perm(qap.size)
		
		// Calculate result
		qap.calculateResult(&s)

		// Search best solution
		//s.searchBestSolution(*qap)

		// Add to poblation
		qap.poblation[i] = s
	}
}

func (qap *QAP) evaluate() QAPSolution {
	return giveBestSolution(&qap.poblation)
}

func (qap *QAP) tournamentSelect(pob_size int, tour_size int) QAPSolution {
	var selected []int = rand.Perm(pob_size)[:tour_size]
	var chosen_array []QAPSolution = make([]QAPSolution, tour_size)

	for i := 0; i < tour_size; i++ {
		chosen_array[i] = qap.poblation[selected[i]]
	}

	return giveBestSolution(&chosen_array)
}

func cmp(arr1, arr2 *[]int) bool {
	array1 := *arr1
	array2 := *arr2

	for i := 0; i < len(array1); i++ {
		if array1[i] != array2[i] {
			return false
		}
	}

	return true
}

func (qap *QAP) obtainElements(pob_size int, tour_size int) (QAPSolution, QAPSolution) {
	var elem1, elem2 QAPSolution
	elem1 = qap.tournamentSelect(pob_size, tour_size)

	for {
		elem2 = qap.tournamentSelect(pob_size, tour_size)
		if !cmp(&elem1.solution, &elem2.solution) {
			break
		}
	}

	return elem1, elem2
}

func (qap *QAP) cross(pob_size int, tour_size int) {
	var elem1, elem2, child1, child2 QAPSolution
	for i := 0; i < pob_size; i += 2 {
		elem1, elem2 = qap.obtainElements(pob_size, tour_size)

		child1, child2 = qap.makeChilds(&elem1, &elem2)

		qap.poblation = append(qap.poblation, child1)
		qap.poblation = append(qap.poblation, child2)
	}

	// Sort array
	sort.SliceStable(qap.poblation, func(i, j int) bool {
		return qap.poblation[i].result < qap.poblation[j].result
	})

	// Remain best solutions
	qap.poblation = qap.poblation[:pob_size]
}

func (qap *QAP) mutate(prob_mutation float64) {
	var rand_number float64
	for i := 0; i < len(qap.poblation); i++ {
		rand_number = float64(rand.Intn(100)) / 100
		if rand_number < prob_mutation {
			// Swap 2 elements
			pos_to_change := rand.Perm(qap.size)[:2]
			ind := qap.poblation[i]
			aux := ind.solution[pos_to_change[0]]
			ind.solution[pos_to_change[0]] = ind.solution[pos_to_change[1]]
			ind.solution[pos_to_change[1]] = aux
		}
	}
}

func (qap *QAP) executeAlgorithm(pob_size int, gen_num int, tour_size int) QAPSolution {
	// Obtain mutation probability
	//var prob_mutation float64 = float64(1 / pob_size)
	//var prob_mutation float64 = 0.2
	var prob_mutation float64 = 1

	// Obtain initial poblation
	qap.obtainInitialPoblation(pob_size)
	
	// Evaluate initial poblation
	var best_solution QAPSolution = qap.evaluate()

	// Print solution
	best_solution.printGenerationResult(0)

	for i := 1; i < gen_num; i++ {
		qap.cross(pob_size, tour_size)
		qap.mutate(prob_mutation)
		best_solution = qap.evaluate()
		best_solution.printGenerationResult(i)
	}

	return best_solution
}

func main() {
	// Take execution parameters
	data, pob_size, gen_num, tour_size := takeParameters(os.Args[1:])

	// Obtain QAP data
	fmt.Println("Initializing data....")
	qap := createQAPObject(&data)

	// Execute QAP algorithm
	fmt.Println("Starting algorithm execution")
	start := time.Now()
	var best_solution QAPSolution = qap.executeAlgorithm(pob_size, gen_num, tour_size)
	end := time.Now()

	// Show results
	fmt.Println("-------------------------------------------------")
	fmt.Println("FINAL RESULTS: ")
	best_solution.printGenerationResult(-1)
	fmt.Println("- Execution time: ", end.Sub(start))
}