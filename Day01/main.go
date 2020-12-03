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
	target   = flag.Int64("target", 2020, "target sum")
	threeSum = flag.Bool("threeSum", false, "find three input lines that sum to target rather than two")
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must specify input file")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	nums, err := twoSum(s)
	if err != nil {
		log.Fatal(err)
	} else if nums == nil {
		log.Print("no solution found")
	}
	log.Printf("Inputs %v sum to %d! Their product is %d", nums, *target, product(nums))
}

func twoSum(s *bufio.Scanner) ([]int64, error) {
	encountered := map[int64]bool{}

	for s.Scan() {
		i, err := strconv.ParseInt(s.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing input as int64: %v, err")
		}
		if encountered[*target-i] {
			return []int64{*target - i, i}, nil
		}
		encountered[i] = true
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
