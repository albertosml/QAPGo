package qap

import (
	"utils"
	"qap_solution"
	"math/rand"
	"sort"
)

type QAP struct {
	size int
	flows [][]int
	distances [][]int
	poblation []qap_solution.QAPSolution
}

func CreateQAPObject(data *[]byte) QAP {
	var fileData []string = utils.ObtainFileData(data)

	// QAP object
	var size int = utils.ObtainMatrixSize(fileData[0])
	if size == 0 {
		return QAP{}
	} else {
		qap := QAP{size: size}
		utils.ObtainMatrix(&qap.flows, size, fileData[1:(size*size)+1])
		utils.ObtainMatrix(&qap.distances, size, fileData[(size*size)+1:])
		return qap
	}
}

func (qap *QAP) calculateResult(qs *qap_solution.QAPSolution) {
	s := *qs
	for i := 0; i < qap.size; i++ {
		for j := 0; j < qap.size; j++ {
			s.Result = s.Result + (qap.flows[i][j] * qap.distances[s.Solution[i]][s.Solution[j]])
		}
	}
	*qs = s
}

func (qap *QAP) makeChilds(ind1 *qap_solution.QAPSolution, ind2 *qap_solution.QAPSolution) (qap_solution.QAPSolution, qap_solution.QAPSolution) {
	// Create references
	elem1 := *ind2
	elem2 := *ind1

	// Create childs
	var child1 qap_solution.QAPSolution = elem1.CreateChild()
	var child2 qap_solution.QAPSolution = elem2.CreateChild()

	var start, end int = utils.ObtainUniqueRandomNumbers(qap.size)
	var pos1, pos2 int

	// Take data from first parent
	for i := start; i < end; i++ {
		child1.Solution[i] = elem1.Solution[i]
		child2.Solution[i] = elem2.Solution[i]
	}

	if start == 0 {
		pos1, pos2 = end, end
	}

	// Take data from second parent
	for i := 0; i < qap.size; i++ {
		if utils.IsIn(elem1.Solution[i], &child2.Solution) == -1 {
			child2.Solution[pos1] = elem1.Solution[i]
			pos1 = pos1 + 1
			if pos1 == start {
				pos1 = end
			}
		}

		if utils.IsIn(elem2.Solution[i], &child1.Solution) == -1 {
			child1.Solution[pos2] = elem2.Solution[i]
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

func (qap *QAP) obtainInitialpoblation(pob_size int) {
	qap.poblation = make([]qap_solution.QAPSolution, pob_size)
	
	for i := 0; i < pob_size; i++ {
		s := qap_solution.QAPSolution{}

		// Initialize a random solution
		s.Solution = rand.Perm(qap.size)
		
		// Calculate result
		qap.calculateResult(&s)

		// Add to poblation
		qap.poblation[i] = s
	}
}

func (qap *QAP) evaluate() qap_solution.QAPSolution {
	return qap_solution.GiveBestSolution(&qap.poblation)
}

func (qap *QAP) tournamentSelect(pob_size int, tour_size int) qap_solution.QAPSolution {
	var selected []int = rand.Perm(pob_size)[:tour_size]
	var chosen_array []qap_solution.QAPSolution = make([]qap_solution.QAPSolution, tour_size)

	for i := 0; i < tour_size; i++ {
		chosen_array[i] = qap.poblation[selected[i]]
	}

	return qap_solution.GiveBestSolution(&chosen_array)
}

func (qap *QAP) obtainElements(pob_size int, tour_size int) (qap_solution.QAPSolution, qap_solution.QAPSolution) {
	var elem1, elem2 qap_solution.QAPSolution
	elem1 = qap.tournamentSelect(pob_size, tour_size)

	for {
		elem2 = qap.tournamentSelect(pob_size, tour_size)
		if !utils.Cmp(&elem1.Solution, &elem2.Solution) {
			break
		}
	}

	return elem1, elem2
}

func (qap *QAP) cross(pob_size int, tour_size int) {
	var elem1, elem2, child1, child2 qap_solution.QAPSolution
	for i := 0; i < pob_size; i += 2 {
		elem1, elem2 = qap.obtainElements(pob_size, tour_size)

		child1, child2 = qap.makeChilds(&elem1, &elem2)

		qap.poblation = append(qap.poblation, child1)
		qap.poblation = append(qap.poblation, child2)
	}

	// Sort array
	sort.SliceStable(qap.poblation, func(i, j int) bool {
		return qap.poblation[i].Result < qap.poblation[j].Result
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
			aux := ind.Solution[pos_to_change[0]]
			ind.Solution[pos_to_change[0]] = ind.Solution[pos_to_change[1]]
			ind.Solution[pos_to_change[1]] = aux
		}
	}
}

func (qap *QAP) ExecuteAlgorithm(pob_size int, gen_num int, tour_size int) qap_solution.QAPSolution {
	// Obtain mutation probability
	var prob_mutation float64 = 1

	// Obtain initial poblation
	qap.obtainInitialpoblation(pob_size)
	
	// Evaluate initial poblation
	var best_solution qap_solution.QAPSolution = qap.evaluate()

	// Print solution
	best_solution.PrintGenerationResult(0)

	for i := 1; i < gen_num; i++ {
		qap.cross(pob_size, tour_size)
		qap.mutate(prob_mutation)
		best_solution = qap.evaluate()
		best_solution.PrintGenerationResult(i)
	}

	return best_solution
}