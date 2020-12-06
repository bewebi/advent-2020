package main

import (
	"bufio"
	"log"
	"os"
)

type groupAnswer struct {
	answers map[rune]int
	count   int
}

func main() {
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	ga := groupAnswer{map[rune]int{}, 0}
	gas := []groupAnswer{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() == "" {
			gas = append(gas, ga)
			ga = groupAnswer{map[rune]int{}, 0}
			continue
		}
		ga.count++
		for _, r := range []rune(s.Text()) {
			if _, ok := ga.answers[r]; ok {
				ga.answers[r]++
				continue
			}
			ga.answers[r] = 1
		}
	}
	gas = append(gas, ga)

	sum := 0
	uSum := 0
	for _, ga := range gas {
		sum += len(ga.answers)
		for _, c := range ga.answers {
			if c == ga.count {
				uSum++
			}
		}
	}
	log.Printf("sum of group answers is %d and the sum of unanimous group answers is %d", sum, uSum)
}
