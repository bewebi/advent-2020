package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	tileRE = regexp.MustCompile(`^Tile ([0-9]+):$`)

	tiles = map[int64]tile{}
	edges = map[string][]int64{}
)

type tile struct {
	top, right, bottom, left string
	rotation                 int
	rows                     []string
}

func main() {
	flag.Parse()

	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() {
		match := tileRE.FindStringSubmatch(s.Text())
		tileNum, _ := strconv.ParseInt(match[1], 10, 64)
		t := tile{}

		// top
		s.Scan()
		row := s.Text()
		t.top = row
		t.left = row[:1]
		t.right = row[len(row)-1:]
		t.rows = []string{row}

		for s.Scan() && s.Text() != "" {
			row = s.Text()
			t.left = row[:1] + t.left
			t.right += row[len(row)-1:]
			t.rows = append(t.rows, row)
		}
		t.bottom = reverseString(row)

		tiles[tileNum] = t

		addEdge(tileNum, t.top)
		addEdge(tileNum, t.right)
		addEdge(tileNum, t.left)
		addEdge(tileNum, t.bottom)
	}

	corners := []int64{}
	for num, t := range tiles {
		if matchedEdges(num, t) == 2 {
			corners = append(corners, num)
		}
	}
	if len(corners) == 4 {
		log.Printf("have four corners %v which multiply to %d", corners, corners[0]*corners[1]*corners[2]*corners[3])
	} else {
		log.Fatalf("wrong number of corners: %d", len(corners))
	}
}

func addEdge(tNum int64, edge string) {
	if ts, ok := edges[edge]; ok {
		edges[edge] = append(ts, tNum)
	} else {
		edges[edge] = []int64{tNum}
	}
}

func matchEdge(tNum int64, edge string) bool {
	if ts, ok := edges[edge]; ok {
		for _, t := range ts {
			if t != tNum {
				return true
			}
		}
	}
	if ts, ok := edges[reverseString(edge)]; ok {
		for _, t := range ts {
			if t != tNum {
				return true
			}
		}
	}
	return false
}

func matchedEdges(tNum int64, t tile) int {
	matched := 0
	if matchEdge(tNum, t.top) {
		matched++
	}
	if matchEdge(tNum, t.right) {
		matched++
	}
	if matchEdge(tNum, t.left) {
		matched++
	}
	if matchEdge(tNum, t.bottom) {
		matched++
	}
	return matched
}

func reverseString(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
