package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	numList := []int{0}

	s := bufio.NewScanner(f)
	for s.Scan() {
		n, err := strconv.ParseInt(s.Text(), 10, 64)
		if err != nil {
			log.Fatalf("cannot parse line %s as int", s.Text())
		}
		numList = append(numList, int(n))
	}

	sort.Ints(numList)
	ones, twos, threes := 0, 0, 0
	for i, num := range numList {
		if i == len(numList)-1 {
			threes++
			continue
		}
		switch numList[i+1] - num {
		case 1:
			ones++
		case 2:
			twos++
		case 3:
			threes++
		default:
			log.Fatalf("unsupported joltage difference")
		}
	}
	log.Printf("%d ones and %d threes multiply to %d", ones, threes, ones*threes)
	log.Printf("there are %d possible arrangements of adapters", countArrangements(numList))
}

func countArrangements(numList []int) int {
	countMap := map[int]int{}
	countMap[numList[len(numList)-1]] = 1
	for i := len(numList) - 2; i >= 0; i-- {
		cur := numList[i]
		countMap[cur] = 0
		if i+3 < len(numList) && (numList[i+3]-cur) <= 3 {
			countMap[cur] += countMap[numList[i+3]]
		}
		if i+2 < len(numList) && (numList[i+2]-cur) <= 3 {
			countMap[cur] += countMap[numList[i+2]]
		}
		if i+1 < len(numList) && (numList[i+1]-cur) <= 3 {
			countMap[cur] += countMap[numList[i+1]]
		}
	}
	return countMap[numList[0]]
}
