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

const (
	N = "N"
	S = "S"
	E = "E"
	W = "W"
	L = "L"
	R = "R"
	F = "F"
)

var (
	directionRE = regexp.MustCompile(fmt.Sprintf(`(%s|%s|%s|%s|%s|%s|%s)([0-9]*)`, N, S, E, W, L, R, F))

	flagWaypoint = flag.Bool("waypoint", false, "use waypoint for directions")

	headings = map[int]string{270: N, 0: E, 90: S, 180: W}
)

type direction struct {
	action string
	value  int
}

func main() {
	flag.Parse()

	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	directions := []direction{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		match := directionRE.FindStringSubmatch(s.Text())
		val, _ := strconv.ParseInt(match[2], 10, 64)
		directions = append(directions, direction{match[1], int(val)})
	}

	var manhattan int
	if *flagWaypoint {
		manhattan = executeDirectionsWaypoint(directions)
	} else {
		manhattan = executeDirections(directions)
	}
	log.Printf("manhattan distance from origin is %d", manhattan)
}

func executeDirections(ds []direction) int {
	var north, east, header int
	for _, d := range ds {
		switch d.action {
		case N:
			north += d.value
		case S:
			north -= d.value
		case E:
			east += d.value
		case W:
			east -= d.value
		case L:
			header -= d.value
		case R:
			header += d.value
		case F:
			for header < 0 {
				header += 360
			}
			switch headings[header%360] {
			case N:
				north += d.value
			case S:
				north -= d.value
			case E:
				east += d.value
			case W:
				east -= d.value

			}
		}
	}
	if north < 0 {
		north *= -1
	}
	if east < 0 {
		east *= -1
	}
	return north + east
}

func executeDirectionsWaypoint(ds []direction) int {
	var north, east, wpNorth, wpEast int
	wpNorth, wpEast = 1, 10

	for _, d := range ds {
		switch d.action {
		case N:
			wpNorth += d.value
		case S:
			wpNorth -= d.value
		case E:
			wpEast += d.value
		case W:
			wpEast -= d.value
		case L:
			newWPNorth, newWPEast := wpNorth, wpEast
			switch d.value {
			case 90:
				newWPNorth = wpEast
				newWPEast = wpNorth * -1
			case 180:
				newWPNorth = wpNorth * -1
				newWPEast = wpEast * -1
			case 270:
				newWPNorth = wpEast * -1
				newWPEast = wpNorth
			}
			wpNorth, wpEast = newWPNorth, newWPEast
		case R:
			newWPNorth, newWPEast := wpNorth, wpEast
			switch d.value {
			case 270:
				newWPNorth = wpEast
				newWPEast = wpNorth * -1
			case 180:
				newWPNorth = wpNorth * -1
				newWPEast = wpEast * -1
			case 90:
				newWPNorth = wpEast * -1
				newWPEast = wpNorth
			}
			wpNorth, wpEast = newWPNorth, newWPEast
		case F:
			north += d.value * wpNorth
			east += d.value * wpEast
		}
	}
	if north < 0 {
		north *= -1
	}
	if east < 0 {
		east *= -1
	}
	return north + east
}
