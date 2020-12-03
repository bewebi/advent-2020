package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

type row []bool
type field []row

var (
	flagSlope = flag.Int("slope", 3, "slope of toboggan")
)

func main() {
	flag.Parse()

	file, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer file.Close()

	f, err := parseFieldFromFile(file)
	if err != nil {
		log.Fatalf("error parsing input field as field: %v", err)
	}

	cnt := f.treeCount(3)
	fmt.Printf("encountered %d trees on the path with slope %d\n", cnt, *flagSlope)
}

func parseFieldFromFile(file *os.File) (field, error) {
	f := field{}
	rowCnt := 0

	s := bufio.NewScanner(file)
	for s.Scan() {
		rowRunes := []rune(s.Text())
		f = append(f, row{})
		for _, r := range rowRunes {
			switch r {
			case '.':
				f[rowCnt] = append(f[rowCnt], false)
			case '#':
				f[rowCnt] = append(f[rowCnt], true)
			default:
				return nil, fmt.Errorf("unexpected input char %v", r)
			}
		}
		rowCnt++
	}
	return f, nil
}

func (f *field) treeCount(slope int) int {
	cnt := 0
	x := 0
	for _, r := range *f {
		if r[x%len(r)] {
			cnt++
		}
		x = x + 3
	}
	return cnt
}
