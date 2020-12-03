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

var (
	flagPositional = flag.Bool("positional", false, "determine password validity by position of char rather than count")
)

type PassInfo struct {
	lBound, uBound int64
	char           rune
	password       string
}

func main() {
	flag.Parse()

	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	validCount, totalCount := 0, 0
	for s.Scan() {
		pi, err := parsePassInfo(s.Text())
		if err != nil {
			log.Fatalf("error parsing input line: %v", err)
		}
		totalCount++
		if *flagPositional {
			if pi.isValidByPosition() {
				validCount++
			}
		} else {
			if pi.isValidByCount() {
				validCount++
			}
		}
	}

	log.Printf("%d valid passwords out of %d parsed", validCount, totalCount)
}

func parsePassInfo(s string) (*PassInfo, error) {
	parts := strings.Split(s, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("improperly formated input")
	}
	bounds := strings.Split(parts[0], "-")
	if len(bounds) != 2 {
		return nil, fmt.Errorf("improperly formatted bounds part")
	}
	lBound, err := strconv.ParseInt(bounds[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse lower bound %s as int64: %v", bounds[0], err)
	}
	uBound, err := strconv.ParseInt(bounds[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse upper bound %s as int64: %v", bounds[1], err)
	}
	return &PassInfo{
		lBound,
		uBound,
		[]rune(parts[1])[0],
		parts[2],
	}, nil
}

func (pi *PassInfo) isValidByCount() bool {
	count := int64(strings.Count(pi.password, string(pi.char)))
	return count >= pi.lBound && count <= pi.uBound
}

func (pi *PassInfo) isValidByPosition() bool {
	passRunes := []rune(pi.password)
	return (passRunes[pi.lBound-1] == pi.char) != (passRunes[pi.uBound-1] == pi.char)
}
