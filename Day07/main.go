package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"
)

type bag struct {
	color    string
	contents map[string]int64
}

var (
	lineRE    = regexp.MustCompile(`([a-z ]*) bags contain ([0-9a-z, ]*)\.`)
	contentRE = regexp.MustCompile(`([0-9]*) ([a-z ]*) bag`)

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
	tcBag, ok := allBags[*flagTargetColor]
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
	lMatch := lineRE.FindStringSubmatch(l) // format: [full, color, content]
	b.color = lMatch[1]
	contentPart := lMatch[2]
	for _, cMatch := range contentRE.FindAllStringSubmatch(contentPart, -1) { // format: [full, num, color]
		num, _ := strconv.ParseInt(cMatch[1], 10, 64) // can ignore error since we're guaranteed numeric input from regex match
		b.contents[cMatch[2]] = num
	}
	return b, nil
}
