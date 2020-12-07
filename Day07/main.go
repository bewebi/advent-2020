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

type bag struct {
	color    string
	contents map[string]int64
}

const (
	BAGSCONTAIN = " bags contain "
	NOOTHERBAGS = "no other bags"
	BAG         = " bag"
)

var (
	allBags         map[string]*bag
	flagTargetColor = flag.String("target-color", "shiny gold", "color of bag to print info for")
)

func main() {
	flag.Parse()

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
		if b.holdsColor(*flagTargetColor) {
			sum++
		}
	}
	log.Printf("%d bags can hold %s bags", sum, *flagTargetColor)
	tcBag, ok := allBags["shiny gold"]
	if !ok {
		log.Fatalf("no bag of color \"%s\" found", *flagTargetColor)
	}
	log.Printf("%s bags hold %d other bags", *flagTargetColor, tcBag.countContents())
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

func (b *bag) countContents() int64 {
	var sum int64
	for cBagColor, i := range b.contents {
		sum += i
		if cBag, ok := allBags[cBagColor]; ok {
			sum += i * cBag.countContents()
		}
	}
	return sum
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
