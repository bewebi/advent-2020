package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type bag struct {
	color    string
	contents map[string]int64
}

const (
	BAGSCONTAIN = " bags contain "
	NOOTHERBAGS = "no other bags"
	BAG         = " bag"
)

var allBags map[string]*bag

func main() {
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	allBags = map[string]*bag{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		b, err := parseLineAsBag(s.Text())
		if err != nil {
			log.Fatalf("error parsing line as bag: %v", err)
		}
		allBags[b.color] = b
	}

	sum := 0
	for _, b := range allBags {
		if b.holdsColor("shiny gold") {
			sum++
		}
	}
	log.Printf("%d bags can hold shiny gold bags", sum)
}

func (b *bag) holdsColor(c string) bool {
	if _, ok := b.contents[c]; ok {
		return true
	}
	for bColor := range b.contents {
		if b, ok := allBags[bColor]; ok && b.holdsColor(c) {
			return true
		}
	}
	return false
}

func parseLineAsBag(l string) (*bag, error) {
	b := &bag{"", map[string]int64{}}
	i := strings.Index(l, BAGSCONTAIN)
	b.color = l[:i]
	contentPart := l[i+len(BAGSCONTAIN) : len(l)-1]
	for _, s := range strings.Split(contentPart, ", ") {
		i := strings.Index(s, NOOTHERBAGS)
		if i >= 0 {
			continue
		}
		numPart := s[:1] // should generalize to allow for multi-digit numbers
		num, err := strconv.ParseInt(numPart, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing numPart %s as int: %v", numPart, err)
		}
		i = strings.Index(s, BAG)
		colorPart := s[2:i]
		b.contents[colorPart] = num
	}
	return b, nil
}
