package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

var (
	flagDays     = flag.Int("days", 100, "number of days to flip")
	flagMaxTiles = flag.Int("max-tiles", 2000, "maximum size of tile grid")
)

func main() {
	flag.Parse()

	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	tiles := make([][]bool, *flagMaxTiles)
	for row := range tiles {
		tiles[row] = make([]bool, *flagMaxTiles)
	}

	for s.Scan() {
		steps := s.Text()
		curX, curY := 1000, 1000
		for steps != "" {
			nonCardinal := len(steps) >= 2
			if nonCardinal {
				switch steps[:2] {
				case "se":
					curX++
					curY++
					steps = steps[2:]
				case "sw":
					curX++
					curY--
					steps = steps[2:]
				case "ne":
					curX--
					curY++
					steps = steps[2:]
				case "nw":
					curX--
					curY--
					steps = steps[2:]
				default:
					nonCardinal = false
				}
			}
			if !nonCardinal {
				if steps[:1] == "e" {
					curY += 2
				} else if steps[:1] == "w" {
					curY -= 2
				} else {
					log.Fatalf("unexpected steps %s", steps)
				}
				steps = steps[1:]
			}
		}

		tiles[curX][curY] = !tiles[curX][curY]
	}
	blackCount := 0
	for _, row := range tiles {
		for _, tile := range row {
			if tile {
				blackCount++
			}
		}
	}
	log.Printf("after all moves there are %d black tiles", blackCount)

	for i := 0; i < *flagDays; i++ {
		tiles = dailyFlip(tiles)
	}
	blackCount = 0
	for _, row := range tiles {
		for _, tile := range row {
			if tile {
				blackCount++
			}
		}
	}
	log.Printf("after %d days of flipping there are %d black tiles", *flagDays, blackCount)
}

func dailyFlip(tiles [][]bool) [][]bool {
	nextTiles := make([][]bool, *flagMaxTiles)
	for row := range nextTiles {
		nextTiles[row] = make([]bool, *flagMaxTiles)
	}

	for x := 1; x < len(tiles)-1; x++ {
		for y := 2; y < len(tiles[0])-2; y++ {
			blackCount := 0
			if tiles[x-1][y-1] {
				blackCount++
			}
			if tiles[x-1][y+1] {
				blackCount++
			}
			if tiles[x][y-2] {
				blackCount++
			}
			if tiles[x][y+2] {
				blackCount++
			}
			if tiles[x+1][y-1] {
				blackCount++
			}
			if tiles[x+1][y+1] {
				blackCount++
			}
			if (tiles[x][y] && (blackCount == 1 || blackCount == 2)) || (!tiles[x][y] && blackCount == 2) {
				nextTiles[x][y] = true
			}
		}
	}
	return nextTiles
}
