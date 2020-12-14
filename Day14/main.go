package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	flagFloatingMask = flag.Bool("floating-mask", false, "use floating masks")
)

func main() {
	flag.Parse()

	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	mem := map[int64]int64{}
	var mask string

	s := bufio.NewScanner(f)
	for s.Scan() {
		inParts := strings.Split(s.Text(), " = ")
		if len(inParts) != 2 {
			log.Fatalf("input line %s is malformed", s.Text())
		}
		if inParts[0] == "mask" {
			mask = inParts[1]
			continue
		}
		memLoc, err := strconv.ParseInt(inParts[0][4:len(inParts[0])-1], 10, 64)
		if err != nil {
			log.Fatalf("cannot parse memLoc %s as int", inParts[0])
		}
		memVal, err := strconv.ParseInt(inParts[1], 10, 64)
		if err != nil {
			log.Fatalf("cannot parse memVal %s as int", inParts[1])
		}

		if *flagFloatingMask {
			for _, loc := range applyFloatingMask(mask, memLoc) {
				mem[loc] = memVal
			}
		} else {
			mem[memLoc] = applyMask(mask, memVal)
		}
	}

	var memSum int64
	for _, val := range mem {
		memSum += val
	}
	log.Printf("sum of memory is %d", memSum)
}

func applyMask(mask string, in int64) int64 {
	for i, m := range []rune(mask) {
		switch m {
		case '0':
			in = applyZero(in, 35-i)
		case '1':
			in = applyOne(in, 35-i)
		}
	}
	return in
}

func applyFloatingMask(mask string, in int64) []int64 {
	outs := []int64{in}
	for i, m := range []rune(mask) {
		switch m {
		case '1':
			for o, out := range outs {
				outs[o] = applyOne(out, 35-i)
			}
		case 'X':
			newOuts := make([]int64, len(outs)*2)
			for o, out := range outs {
				newOuts[o*2] = applyZero(out, 35-i)
				newOuts[o*2+1] = applyOne(out, 35-i)
			}
			outs = newOuts
		}
	}
	return outs
}

func applyZero(in int64, pos int) int64 {
	mask := int64(math.MaxInt64) - (1 << pos)
	return mask & int64(in)
}

func applyOne(in int64, pos int) int64 {
	mask := int64(1 << pos)
	return mask | in
}
