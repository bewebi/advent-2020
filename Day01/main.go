package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
)

var (
	target = flag.Int64("target", 2020, "target sum")
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

	encountered := map[int64]bool{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		i, err := strconv.ParseInt(s.Text(), 10, 64)
		if err != nil {
			log.Fatal("error parsing input as int64")
		}
		if encountered[*target-i] {
			log.Printf("Inputs %d and %d sum to %d! Their product is %d", *target-i, i, *target, (*target-i)*i)
			return
		}
		encountered[i] = true
	}
	log.Print("no solution found")
}
