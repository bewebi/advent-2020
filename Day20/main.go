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

	tiles = map[int64][]string{}
	edges = map[string][]int64{}
)

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
		tile := make([]string, 4)

		// top
		s.Scan()
		row := s.Text()
		tile[0] = row
		tile[1] = row[:1]
		tile[3] = row[len(row)-1:]

		for s.Scan() && s.Text() != "" {
			row = s.Text()
			tile[1] += row[:1]
			tile[3] += row[len(row)-1:]
		}
		tile[2] = row

		tiles[tileNum] = tile
		for _, edge := range tile {
			if ts, ok := edges[edge]; ok {
				edges[edge] = append(ts, tileNum)
			} else {
				edges[edge] = []int64{tileNum}
			}
		}
	}

	corners := []int64{}
	for num, tile := range tiles {
		if matchedEdges(num, tile) == 2 {
			corners = append(corners, num)
		}
	}
	if len(corners) == 4 {
		log.Printf("have four corners %v which multiply to %d", corners, corners[0]*corners[1]*corners[2]*corners[3])
	} else {
		log.Fatalf("wrong number of corners: %d", len(corners))
	}
}

func matchedEdges(tNum int64, tile []string) int {
	matched := 0
	for _, edge := range tile {
		if ts, ok := edges[edge]; ok {
			found := false
			for _, t := range ts {
				if t != tNum {
					matched++
					break
				}
			}
			if found {
				continue
			}
		}
		if _, ok := edges[reverseString(edge)]; ok {
			matched++
		}
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
