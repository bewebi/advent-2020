package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
)

var (
	flagInput = flag.String("input", "12,1,16,3,11,0", "starting numbers")
	flagLimit = flag.Int("limit", 2020, "number of turns")
)

func main() {
	flag.Parse()

	numLog := map[int][]int{}
	var last int

	inStrings := strings.Split(*flagInput, ",")
	for i, in := range inStrings {
		sn, err := strconv.ParseInt(in, 10, 64)
		if err != nil {
			log.Fatalf("cannot parse input %s as int", in)
		}
		numLog[int(sn)] = []int{i + 1}
		last = int(sn)
	}

	for turn := len(inStrings) + 1; turn <= *flagLimit; turn++ {
		next := 0
		turns, _ := numLog[last]
		if len(turns) > 1 {
			next = turns[1] - turns[0]
		}
		if turns, ok := numLog[next]; ok {
			numLog[next] = []int{turns[len(turns)-1], turn}
		} else {
			numLog[next] = []int{turn}
		}
		last = next
	}
	log.Printf("on turn %d, the elves said %d", *flagLimit, last)
}
