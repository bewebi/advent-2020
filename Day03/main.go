package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type row []bool
type field []row

type slope struct {
	x, y int64
}

var (
	flagSlopeInput = flag.String("slope-input", "slope-input-1.txt", "input file with slopes for toboggan")
)

func main() {
	flag.Parse()

	fieldFile, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer fieldFile.Close()

	f, err := parseFieldFromFile(fieldFile)
	if err != nil {
		log.Fatalf("error parsing input file as field: %v", err)
	}

	slopeFile, err := os.Open(*flagSlopeInput)
	if err != nil {
		log.Fatalf("error opening slope input file: %v", err)
	}
	defer slopeFile.Close()

	slopes, err := parseSlopesFromFile(slopeFile)
	if err != nil {
		log.Fatalf("error parsing slope file: %v", err)
	}
	cnts := []int{}
	for _, s := range slopes {
		cnt := f.treeCount(int(s.x), int(s.y))
		log.Printf("encountered %d trees on the path with slope %d,%d\n", cnt, s.x, s.y)
		cnts = append(cnts, cnt)
	}

	log.Printf("product of tree counts from all slopes: %d", product(cnts))
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

func parseSlopesFromFile(file *os.File) ([]slope, error) {
	slopes := []slope{}
	s := bufio.NewScanner(file)
	for s.Scan() {
		slopeParts := strings.Split(s.Text(), ",")
		if len(slopeParts) != 2 {
			return nil, fmt.Errorf("unexpected number of slope parts %d", len(slopeParts))
		}
		x, err := strconv.ParseInt(slopeParts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse slope x part as int: %v", err)
		}
		y, err := strconv.ParseInt(slopeParts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse slope y part as int: %v", err)
		}
		slopes = append(slopes, slope{x, y})
	}
	return slopes, nil
}

func (f *field) treeCount(slopeX, slopeY int) int {
	cnt := 0
	x := 0
	for y, r := range *f {
		if y%slopeY != 0 {
			continue
		}
		if r[x%len(r)] {
			cnt++
		}
		x = x + slopeX
	}
	return cnt
}

func product(nums []int) int {
	ret := 1
	for _, i := range nums {
		ret *= i
	}
	return ret
}
