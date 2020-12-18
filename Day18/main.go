package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	flagAdditionFirst = flag.Bool("addition-first", false, "addition first math")

	parenRE = regexp.MustCompile(`\([0-9+* ]+\)`)
)

func main() {
	flag.Parse()

	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var result int64
	for s.Scan() {
		if *flagAdditionFirst {
			result += evaluateExpressionAdditionFirst(s.Text())
		} else {
			result += evaluateExpression(s.Text())
		}
	}
	log.Printf("sum of all expressions is %d", result)
}

func evaluateExpression(exp string) int64 {
	for {
		parenExp := parenRE.FindString(exp)
		if parenExp == "" {
			break
		}
		exp = strings.Replace(exp, parenExp, fmt.Sprintf("%d", evaluateExpression(parenExp[1:len(parenExp)-1])), -1)
	}
	parts := strings.Split(exp, " ")

	cur, _ := strconv.ParseInt(parts[0], 10, 64)
	for i := 1; i < len(parts); i += 2 {
		next, _ := strconv.ParseInt(parts[i+1], 10, 64)
		switch parts[i] {
		case "+":
			cur += next
		case "*":
			cur *= next
		default:
			log.Fatalf("unexpected character %s", parts[i])
		}
	}
	return cur
}

func evaluateExpressionAdditionFirst(exp string) int64 {
	for {
		parenExp := parenRE.FindString(exp)
		if parenExp == "" {
			break
		}
		exp = strings.Replace(exp, parenExp, fmt.Sprintf("%d", evaluateExpressionAdditionFirst(parenExp[1:len(parenExp)-1])), -1)
	}
	parts := strings.Split(exp, " ")
	for {
		plusFound := false
		for i, p := range parts {
			if p == "+" {
				prev, _ := strconv.ParseInt(parts[i-1], 10, 64)
				next, _ := strconv.ParseInt(parts[i+1], 10, 64)
				sum := fmt.Sprintf("%d", prev+next)
				newParts := append(parts[:i-1], sum)
				parts = append(newParts, parts[i+2:]...)
				plusFound = true
				break
			}
		}
		if !plusFound {
			break
		}
	}
	res := int64(1)
	for i := 0; i < len(parts); i += 2 {
		next, _ := strconv.ParseInt(parts[i], 10, 64)
		res *= next
	}
	return res
}
