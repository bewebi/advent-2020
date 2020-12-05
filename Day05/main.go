package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
)

func main() {
	flag.Parse()

	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	maxID := 0
	s := bufio.NewScanner(f)
	for s.Scan() {
		seatID, err := getSeatIDFromCode(s.Text())
		if err != nil {
			log.Fatalf("error getting seat ID from code: %v", err)
		}
		if seatID > maxID {
			maxID = seatID
		}
	}

	log.Printf("the highest seat ID is %d", maxID)
}

func getSeatIDFromCode(c string) (int, error) {
	if len(c) != 10 {
		return 0, fmt.Errorf("code has length %d, must be 10", len(c))
	}
	rowBs := make([]bool, 7)
	for i := 0; i < 7; i++ {
		switch c[i] {
		case 'F':
			rowBs[i] = false
		case 'B':
			rowBs[i] = true
		default:
			return 0, fmt.Errorf("malformed code %s, expected B or F at index %d", c, i)
		}
	}
	colBs := make([]bool, 3)
	for i := 7; i < 10; i++ {
		switch c[i] {
		case 'L':
			colBs[i-7] = false
		case 'R':
			colBs[i-7] = true
		default:
			return 0, fmt.Errorf("malformed code %s, expected L or R at index %d", c, i)
		}
	}
	row := bspFromBools(rowBs)
	col := bspFromBools(colBs)
	return (row * 8) + col, nil
}

func bspFromBools(bs []bool) int {
	sum := 0
	for i := 0; i < len(bs); i++ {
		if bs[i] {
			sum += int(math.Pow(2, float64((len(bs) - i - 1))))
		}
	}
	return sum
}
