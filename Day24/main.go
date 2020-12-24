package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	tiles := map[int]map[int]bool{}

	for s.Scan() {
		steps := s.Text()
		curX, curY := 0, 0
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
		if _, ok := tiles[curX]; !ok {
			tiles[curX] = map[int]bool{}
		}
		tiles[curX][curY] = !tiles[curX][curY]
	}
	blackCount, total := 0, 0
	for _, row := range tiles {
		for _, tile := range row {
			total++
			if tile {
				blackCount++
			}
		}
	}
	log.Printf("after all moves there are %d black tiles of %d total", blackCount, total)
}
