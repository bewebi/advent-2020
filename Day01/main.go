package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	flagTarget   = flag.Int64("target", 2020, "target sum")
	flagThreeSum = flag.Bool("three-sum", false, "find three input lines that sum to target rather than two")
)

func main() {
	flag.Parse()
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var nums []int64
	if *flagThreeSum {
		nums, err = threeSum(s, *flagTarget)
	} else {
		nums, err = twoSum(s, *flagTarget)
	}
	if err != nil {
		log.Fatal(err)
	} else if nums == nil {
		log.Print("no solution found")
	}
	log.Printf("Inputs %v sum to %d! Their product is %d", nums, *flagTarget, product(nums))
}

func twoSum(s *bufio.Scanner, target int64) ([]int64, error) {
	encountered := map[int64]bool{}

	for s.Scan() {
		i, err := strconv.ParseInt(s.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing input as int64: %v, err")
		}
		if encountered[target-i] {
			return []int64{target - i, i}, nil
		}
		encountered[i] = true
	}
	return nil, nil
}

func threeSum(s *bufio.Scanner, target int64) ([]int64, error) {
	inputs := map[int64]bool{}
	twoSums := map[int64][]int64{}

	for s.Scan() {
		i, err := strconv.ParseInt(s.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing input as int64: %v, err")
		}
		if nums, ok := twoSums[target-i]; ok {
			return append(nums, i), nil
		}
		for n := range inputs {
			twoSums[i+n] = []int64{i, n}
		}
		inputs[i] = true
	}
	return nil, nil
}

func product(nums []int64) int64 {
	ret := int64(1)
	for _, i := range nums {
		ret *= i
	}
	return ret
}
