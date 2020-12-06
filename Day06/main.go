package main

import (
	"bufio"
	"log"
	"os"
)

type groupAnswer map[rune]bool

func main() {
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	ga := groupAnswer{}
	gas := []groupAnswer{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() == "" {
			gas = append(gas, ga)
			ga = groupAnswer{}
			continue
		}
		for _, r := range []rune(s.Text()) {
			ga[r] = true
		}
	}
	gas = append(gas, ga)

	sum := 0
	for _, ga := range gas {
		sum += len(ga)
	}
	log.Printf("sum of group answers is %d", sum)
}
