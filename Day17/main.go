package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

var (
	flagCycles = flag.Int("cycles", 6, "number of cycles in boot sequence")
)

func main() {
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	initialGrid := [][]bool{}
	inY := 0
	for s.Scan() {
		initialGrid = append(initialGrid, []bool{})
		for _, r := range []rune(s.Text()) {
			if r == '#' {
				initialGrid[inY] = append(initialGrid[inY], true)
			} else {
				initialGrid[inY] = append(initialGrid[inY], false)
			}
		}
		inY++
	}

	cubeSize := *flagCycles*2 + 1
	cubes := make([][][]bool, cubeSize+len(initialGrid))
	for x := range cubes {
		cubes[x] = make([][]bool, cubeSize+len(initialGrid[0]))
		for y := range cubes[x] {
			cubes[x][y] = make([]bool, cubeSize)
		}
	}
	hypercubes := make([][][][]bool, cubeSize+len(initialGrid))
	for x := range hypercubes {
		hypercubes[x] = make([][][]bool, cubeSize+len(initialGrid[0]))
		for y := range hypercubes[x] {
			hypercubes[x][y] = make([][]bool, cubeSize)
			for z := range hypercubes[x][y] {
				hypercubes[x][y][z] = make([]bool, cubeSize)
			}
		}
	}

	for x, ys := range initialGrid {
		for y, active := range ys {
			cubes[x+*flagCycles][y+*flagCycles][*flagCycles] = active
			hypercubes[x+*flagCycles][y+*flagCycles][*flagCycles][*flagCycles] = active
		}
	}

	for i := 0; i < *flagCycles; i++ {
		ncs := newCubes(cubes)
		for x := 0; x < len(cubes); x++ {
			for y := 0; y < len(cubes[0]); y++ {
				for z := 0; z < len(cubes[0][0]); z++ {
					ncs[x][y][z] = shouldCubeBeActive(x, y, z, cubes)
				}
			}
		}
		cubes = ncs
	}
	log.Printf("after %d cycles there are %d active cubes", *flagCycles, countActiveCubes(cubes))

	for i := 0; i < *flagCycles; i++ {
		nhs := newHypercubes(hypercubes)
		for x := 0; x < len(hypercubes); x++ {
			for y := 0; y < len(hypercubes[0]); y++ {
				for z := 0; z < len(hypercubes[0][0]); z++ {
					for w := 0; w < len(hypercubes[0][0][0]); w++ {
						nhs[x][y][z][w] = shouldHypercubeBeActive(x, y, z, w, hypercubes)
					}
				}
			}
		}
		hypercubes = nhs
	}
	log.Printf("after %d cycles there are %d active hypercubes", *flagCycles, countActiveHypercubes(hypercubes))
}

func shouldCubeBeActive(x, y, z int, cubes [][][]bool) bool {
	activeNeighbors := 0
	for xDiff := -1; xDiff <= 1; xDiff++ {
		if x+xDiff < 0 || x+xDiff >= len(cubes) {
			continue
		}
		for yDiff := -1; yDiff <= 1; yDiff++ {
			if y+yDiff < 0 || y+yDiff >= len(cubes[0]) {
				continue
			}
			for zDiff := -1; zDiff <= 1; zDiff++ {
				if z+zDiff < 0 || z+zDiff >= len(cubes[0][0]) {
					continue
				}
				if xDiff == 0 && yDiff == 0 && zDiff == 0 {
					continue
				}
				if cubes[x+xDiff][y+yDiff][z+zDiff] {
					activeNeighbors++
				}
				if activeNeighbors > 3 {
					return false
				}
			}
		}
	}
	if activeNeighbors == 3 || (activeNeighbors == 2 && cubes[x][y][z]) {
		return true
	}
	return false
}

func newCubes(cubes [][][]bool) [][][]bool {
	newCubes := make([][][]bool, len(cubes))
	for x := range newCubes {
		newCubes[x] = make([][]bool, len(cubes[0]))
		for y := range newCubes[x] {
			newCubes[x][y] = make([]bool, len(cubes[0][0]))
		}
	}
	return newCubes
}

func countActiveCubes(cubes [][][]bool) int {
	active := 0
	for x := 0; x < len(cubes); x++ {
		for y := 0; y < len(cubes[0]); y++ {
			for z := 0; z < len(cubes[0][0]); z++ {
				if cubes[x][y][z] {
					active++
				}
			}
		}
	}
	return active
}

func shouldHypercubeBeActive(x, y, z, w int, hypercubes [][][][]bool) bool {
	activeNeighbors := 0
	for xDiff := -1; xDiff <= 1; xDiff++ {
		if x+xDiff < 0 || x+xDiff >= len(hypercubes) {
			continue
		}
		for yDiff := -1; yDiff <= 1; yDiff++ {
			if y+yDiff < 0 || y+yDiff >= len(hypercubes[0]) {
				continue
			}
			for zDiff := -1; zDiff <= 1; zDiff++ {
				if z+zDiff < 0 || z+zDiff >= len(hypercubes[0][0]) {
					continue
				}
				for wDiff := -1; wDiff <= 1; wDiff++ {
					if w+wDiff < 0 || w+wDiff >= len(hypercubes[0][0][0]) {
						continue
					}
					if xDiff == 0 && yDiff == 0 && zDiff == 0 && wDiff == 0 {
						continue
					}
					if hypercubes[x+xDiff][y+yDiff][z+zDiff][w+wDiff] {
						activeNeighbors++
					}
					if activeNeighbors > 3 {
						return false
					}
				}
			}
		}
	}
	if activeNeighbors == 3 || (activeNeighbors == 2 && hypercubes[x][y][z][w]) {
		return true
	}
	return false
}

func newHypercubes(hypercubes [][][][]bool) [][][][]bool {
	newHypercubes := make([][][][]bool, len(hypercubes))
	for x := range newHypercubes {
		newHypercubes[x] = make([][][]bool, len(hypercubes[0]))
		for y := range newHypercubes[x] {
			newHypercubes[x][y] = make([][]bool, len(hypercubes[0][0]))
			for z := range newHypercubes[x][y] {
				newHypercubes[x][y][z] = make([]bool, len(hypercubes[0][0][0]))
			}
		}
	}
	return newHypercubes
}

func countActiveHypercubes(hypercubes [][][][]bool) int {
	active := 0
	for x := 0; x < len(hypercubes); x++ {
		for y := 0; y < len(hypercubes[0]); y++ {
			for z := 0; z < len(hypercubes[0][0]); z++ {
				for w := 0; w < len(hypercubes[0][0][0]); w++ {
					if hypercubes[x][y][z][w] {
						active++
					}
				}
			}
		}
	}
	return active
}
