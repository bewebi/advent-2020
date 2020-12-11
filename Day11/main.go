package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

type layout [][]rune

var (
	flagFarSight = flag.Bool("far-sight", false, "passangers observe seats beyond their immediate vicinity")
	flagNCount   = flag.Int("ncount", 4, "number of adjacent neighbors passengers will tolerate")
)

func main() {
	flag.Parse()

	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	l := layout{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		l = append(l, []rune(s.Text()))
	}

	prev := 0
	for {
		seatCount := l.applyRules(*flagNCount, *flagFarSight)
		if seatCount == prev {
			break
		}
		prev = seatCount
	}

	log.Printf("seat count stabalized at %d", prev)
}

func (l layout) applyRules(nCount int, farSight bool) int {
	tmp := layout{}
	for i, row := range l {
		tmp = append(tmp, []rune{})
		for _, seat := range row {
			tmp[i] = append(tmp[i], seat)
		}
	}
	seated := 0

	for r, row := range l {
		for c, seat := range row {
			switch seat {
			case '.':
				continue
			case 'L':
				if tmp.occupiedSeat(r, c, -1, -1, farSight) ||
					tmp.occupiedSeat(r, c, -1, 0, farSight) ||
					tmp.occupiedSeat(r, c, -1, 1, farSight) ||
					tmp.occupiedSeat(r, c, 0, -1, farSight) ||
					tmp.occupiedSeat(r, c, 0, 1, farSight) ||
					tmp.occupiedSeat(r, c, 1, -1, farSight) ||
					tmp.occupiedSeat(r, c, 1, 0, farSight) ||
					tmp.occupiedSeat(r, c, 1, 1, farSight) {
					continue
				}
				seated++
				l[r][c] = '#'
			case '#':
				neighbors := 0
				if tmp.occupiedSeat(r, c, -1, -1, farSight) {
					neighbors++
				}
				if tmp.occupiedSeat(r, c, -1, 0, farSight) {
					neighbors++
				}
				if tmp.occupiedSeat(r, c, -1, 1, farSight) {
					neighbors++
				}
				if tmp.occupiedSeat(r, c, 0, -1, farSight) {
					neighbors++
				}
				if tmp.occupiedSeat(r, c, 0, 1, farSight) {
					neighbors++
				}
				if tmp.occupiedSeat(r, c, 1, -1, farSight) {
					neighbors++
				}
				if tmp.occupiedSeat(r, c, 1, 0, farSight) {
					neighbors++
				}
				if tmp.occupiedSeat(r, c, 1, 1, farSight) {
					neighbors++
				}
				if neighbors >= nCount {
					l[r][c] = 'L'
					continue
				}
				seated++
			}
		}
	}

	return seated
}

func (l layout) occupiedSeat(myR, myC, rDiff, cDiff int, farSight bool) bool {
	checkR, checkC := myR+rDiff, myC+cDiff
	for checkR >= 0 && checkR < len(l) && checkC >= 0 && checkC < len(l[0]) {
		if seat := l[checkR][checkC]; seat != '.' {
			return seat == '#'
		}
		if !farSight {
			return false
		}
		checkR += rDiff
		checkC += cDiff
	}
	return false
}

func (l layout) print() {
	for _, row := range l {
		for _, seat := range row {
			fmt.Printf("%s", string(seat))
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}
