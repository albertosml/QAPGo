package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"math/rand"
	"strings"
)

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

func TakeParameters(args []string) ([]byte, int, int, int){
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

// To make a function exportable, we need to put the first letter in uppercase mode
func ObtainFileData(data *[]byte) []string {
	var text string = string(*data)
	return strings.Fields(text)
}

func ObtainMatrixSize(number string) int {
	var size int
	size, err := strconv.Atoi(number)

	if err != nil {
		fmt.Println("We expect that matrix size is ubicated in first position")
		return 0
	} else {
		return size
	}
}

func ObtainMatrix(matrix *[][]int, size int, subfileData []string) {
	m := make([][]int, size)

	for i := 0; i < size; i++ {
		m[i] = make([]int, size)
		for j := 0; j < size; j++ {
			m[i][j], _ = strconv.Atoi(subfileData[i*size + j])
		}
	}

	*matrix = m
}

func ObtainUniqueRandomNumbers(max int) (int, int) {
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

func Cmp(arr1, arr2 *[]int) bool {
	array1 := *arr1
	array2 := *arr2

	for i := 0; i < len(array1); i++ {
		if array1[i] != array2[i] {
			return false
		}
	}

	return true
}