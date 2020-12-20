package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	tileRE = regexp.MustCompile(`^Tile ([0-9]+):$`)

	tiles             = map[int64]tile{}
	edges             = map[string][]int64{}
	tileGrid          = [][]tile{}
	bitGrid           = [][]bool{}
	seamonsterMap     = [][]bool{}
	seamonsterPattern = [][]bool{
		[]bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, false},
		[]bool{true, false, false, false, false, true, true, false, false, false, false, true, true, false, false, false, false, true, true, true},
		[]bool{false, true, false, false, true, false, false, true, false, false, true, false, false, true, false, false, true, false, false},
	}
)

type tile struct {
	id           int64
	edges        [4]string
	rotation     int
	flipH, flipV bool
	rows         []string
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
		t := tile{id: tileNum}

		// top
		s.Scan()
		row := s.Text()
		t.edges[0] = row
		t.edges[3] = row[:1]
		t.edges[1] = row[len(row)-1:]
		t.rows = []string{row}

		for s.Scan() && s.Text() != "" {
			row = s.Text()
			t.edges[3] = row[:1] + t.edges[3]
			t.edges[1] += row[len(row)-1:]
			t.rows = append(t.rows, row)
		}
		t.edges[2] = reverseString(row)

		tiles[tileNum] = t

		for _, edge := range t.edges {
			addEdge(tileNum, edge)
		}
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

	buildPicture(corners)
	generateBitGrid()
	log.Printf("roughness of the sea is %d", getRoughness())
}

func buildPicture(corners []int64) {
	ulCorner := tiles[corners[0]]

	for orientation, edge := range ulCorner.edges {
		if match := edgeMatch(ulCorner.id, edge); match >= 0 {
			switch orientation {
			case 0:
				if match := edgeMatch(ulCorner.id, ulCorner.edges[3]); match >= 0 {
					ulCorner.rotate180()
				} else {
					ulCorner.rotate90()
				}
			case 2:
				ulCorner.rotate270()
			}
			tileGrid = [][]tile{[]tile{ulCorner}}
			break
		}
	}

	// first row ends with corner
	c := 1
	for {
		prevTile := tileGrid[0][c-1]
		nextTile := tiles[edgeMatch(prevTile.id, prevTile.edges[1])]
		leftO := nextTile.findEdge(prevTile.edges[1])
		switch leftO {
		case 0:
			nextTile.rotate270()
		case 1:
			nextTile.rotate180()
		case 2:
			nextTile.rotate90()
		}
		if nextTile.edges[3] == prevTile.edges[1] {
			nextTile.flipVertical()
		}

		tileGrid[0] = append(tileGrid[0], nextTile)
		if isMember(nextTile.id, corners) {
			break
		}
		c++
	}

	r := 1
	for {
		prevRowStart := tileGrid[r-1][0]
		nextRowStart := tiles[edgeMatch(prevRowStart.id, prevRowStart.edges[2])]

		topO := nextRowStart.findEdge(prevRowStart.edges[2])
		switch topO {
		case 1:
			nextRowStart.rotate270()
		case 2:
			nextRowStart.rotate180()
		case 3:
			nextRowStart.rotate90()
		}
		if nextRowStart.edges[0] == prevRowStart.edges[2] {
			nextRowStart.flipHorizontal()
		}

		tileGrid = append(tileGrid, []tile{nextRowStart})

		for c := 1; c < len(tileGrid[0]); c++ {
			prevTile := tileGrid[r][c-1]
			nextTile := tiles[edgeMatch(prevTile.id, prevTile.edges[1])]

			leftO := nextTile.findEdge(prevTile.edges[1])
			switch leftO {
			case 0:
				nextTile.rotate270()
			case 1:
				nextTile.rotate180()
			case 2:
				nextTile.rotate90()
			}
			if nextTile.edges[3] == prevTile.edges[1] {
				nextTile.flipVertical()
			}
			tileGrid[r] = append(tileGrid[r], nextTile)
		}

		if isMember(tileGrid[r][0].id, corners) {
			break
		}
		r++
	}
}

func generateBitGrid() {
	rowLen := len(tileGrid[0][0].rows[0]) - 2
	bitGrid = make([][]bool, len(tileGrid)*rowLen)
	for r := range bitGrid {
		bitGrid[r] = make([]bool, len(tileGrid[0])*rowLen)
	}

	for r, row := range tileGrid {
		for c, tile := range row {
			tileBits := tile.bits()
			for tileR, tileRow := range tileBits {
				for tileC, bit := range tileRow {
					bitGrid[r*rowLen+tileR][c*rowLen+tileC] = bit
				}
			}
		}
	}
}

func getRoughness() int {
	foundMonster := false

	for flipTry := 0; flipTry < 4; flipTry++ {
		switch flipTry {
		case 1, 3: // flip H
			for r := 0; r < len(bitGrid); r++ {
				for c := 0; c < len(bitGrid[0])/2; c++ {
					bitGrid[r][c], bitGrid[r][len(bitGrid[0])-c-1] = bitGrid[r][len(bitGrid[0])-c-1], bitGrid[r][c]
				}
			}
		case 2: // flip V
			for r := 0; r < len(bitGrid)/2; r++ {
				for c := 0; c < len(bitGrid[0]); c++ {
					bitGrid[r][c], bitGrid[len(bitGrid)-r-1][c] = bitGrid[len(bitGrid)-r-1][c], bitGrid[r][c]
				}
			}
		}

		for rotateTry := 0; rotateTry < 4; rotateTry++ {
			seamonsterMap = make([][]bool, len(bitGrid))
			for r := range seamonsterMap {
				seamonsterMap[r] = make([]bool, len(bitGrid[0]))
			}

			for r, row := range bitGrid {
				for c, bit := range row {
					seamonsterMap[r][c] = bit
				}
			}

			for r := 0; r+len(seamonsterPattern) < len(bitGrid); r++ {
				for c := 0; c+len(seamonsterPattern[0]) < len(bitGrid[0]); c++ {
					foundMonster = findAndMarkSeamonster(r, c) || foundMonster
				}
			}
			if foundMonster {
				return roughCount()
			}
			rotateBitGrid90()
		}

	}
	return -1
}

func findAndMarkSeamonster(r, c int) bool {
	for sR, row := range seamonsterPattern {
		for sC, bit := range row {
			if bit && !bitGrid[r+sR][c+sC] {
				return false
			}
		}
	}

	for sR, row := range seamonsterPattern {
		for sC, bit := range row {
			if bit {
				seamonsterMap[r+sR][c+sC] = false
			}
		}
	}
	return true
}

func roughCount() int {
	roughCount := 0
	for _, row := range seamonsterMap {
		for _, bit := range row {
			if bit {
				roughCount++
			}
		}
	}
	return roughCount
}

func addEdge(tNum int64, edge string) {
	if ts, ok := edges[edge]; ok {
		edges[edge] = append(ts, tNum)
	} else {
		edges[edge] = []int64{tNum}
	}
}

func (t *tile) rotate90() {
	t.edges = [4]string{t.edges[3], t.edges[0], t.edges[1], t.edges[2]}
	t.rotation = (t.rotation + 90) % 360
}

func (t *tile) rotate180() {
	t.edges = [4]string{t.edges[2], t.edges[3], t.edges[0], t.edges[1]}
	t.rotation = (t.rotation + 180) % 360
}

func (t *tile) rotate270() {
	t.edges = [4]string{t.edges[1], t.edges[2], t.edges[3], t.edges[0]}
	t.rotation = (t.rotation + 270) % 360
}

func (t *tile) flipHorizontal() {
	t.edges = [4]string{reverseString(t.edges[0]), reverseString(t.edges[3]), reverseString(t.edges[2]), reverseString(t.edges[1])}
	t.flipH = !t.flipH
}

func (t *tile) flipVertical() {
	t.edges = [4]string{reverseString(t.edges[2]), reverseString(t.edges[1]), reverseString(t.edges[0]), reverseString(t.edges[3])}
	t.flipV = !t.flipV
}

func (t *tile) findEdge(edge string) int {
	for i, e := range t.edges {
		if edge == e || reverseString(edge) == e {
			return i
		}
	}
	return -1
}

func (t *tile) bits() [][]bool {
	rowLen := len(t.rows[0]) - 2
	bitmap := make([][]bool, rowLen)
	for r := range bitmap {
		bitmap[r] = make([]bool, rowLen)
	}

	for r := 0; r < rowLen; r++ {
		runes := []rune(t.rows[r+1])
		for c := 0; c < rowLen; c++ {
			if runes[c+1] == '#' {
				bitmap[r][c] = true
			}
		}
	}
	rotatedBitmap := make([][]bool, rowLen)
	for r := range rotatedBitmap {
		rotatedBitmap[r] = make([]bool, rowLen)
	}
	switch t.rotation {
	case 0:
		for r := 0; r < rowLen; r++ {
			for c := 0; c < rowLen; c++ {
				rotatedBitmap[r][c] = bitmap[r][c]
			}
		}
	case 90:
		for r := 0; r < rowLen; r++ {
			for c := 0; c < rowLen; c++ {
				rotatedBitmap[c][rowLen-r-1] = bitmap[r][c]
			}
		}
	case 180:
		for r := 0; r < rowLen; r++ {
			for c := 0; c < rowLen; c++ {
				rotatedBitmap[rowLen-r-1][rowLen-c-1] = bitmap[r][c]
			}
		}
	case 270:
		for r := 0; r < rowLen; r++ {
			for c := 0; c < rowLen; c++ {
				rotatedBitmap[rowLen-c-1][r] = bitmap[r][c]
			}
		}
	}
	if t.flipH {
		for r := 0; r < rowLen; r++ {
			for c := 0; c < rowLen/2; c++ {
				rotatedBitmap[r][c], rotatedBitmap[r][rowLen-c-1] = rotatedBitmap[r][rowLen-c-1], rotatedBitmap[r][c]
			}
		}
	}
	if t.flipV {
		for r := 0; r < rowLen/2; r++ {
			for c := 0; c < rowLen; c++ {
				rotatedBitmap[r][c], rotatedBitmap[rowLen-r-1][c] = rotatedBitmap[rowLen-r-1][c], rotatedBitmap[r][c]
			}
		}
	}

	return rotatedBitmap
}

func rotateBitGrid90() {
	rotated := make([][]bool, len(bitGrid[0]))
	for r := range rotated {
		rotated[r] = make([]bool, len(bitGrid))
	}
	for r, row := range bitGrid {
		for c, bit := range row {
			rotated[c][len(bitGrid)-r-1] = bit
		}
	}
	bitGrid = rotated
}

func edgeMatch(tNum int64, edge string) int64 {
	if ts, ok := edges[edge]; ok {
		for _, t := range ts {
			if t != tNum {
				return t
			}
		}
	}
	if ts, ok := edges[reverseString(edge)]; ok {
		for _, t := range ts {
			if t != tNum {
				return t
			}
		}
	}
	return -1
}

func matchedEdges(tNum int64, t tile) int {
	matched := 0
	for _, edge := range t.edges {
		if edgeMatch(tNum, edge) >= 0 {
			matched++
		}
	}
	return matched
}

// utils
func reverseString(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func isMember(i int64, ns []int64) bool {
	for _, n := range ns {
		if i == n {
			return true
		}
	}
	return false
}

// debug helpers
func printGrid(grid [][]bool) {
	for _, row := range grid {
		for _, bit := range row {
			if bit {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}

func printTiles(grid [][]tile) {
	for _, row := range grid {
		for _, tile := range row {
			fmt.Printf(" %d ", tile.id)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}
