package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	ruleRE      = regexp.MustCompile(`^([0-9]+): ([0-9a-z |"]+)$`)
	quoteRuleRE = regexp.MustCompile(`"([a-z])"`)

	flagLoops = flag.Bool("loops", false, "use looping rules")

	rules = map[int64]string{}
)

func main() {
	flag.Parse()

	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() && s.Text() != "" {
		match := ruleRE.FindStringSubmatch(s.Text())
		num, _ := strconv.ParseInt(match[1], 10, 64)
		rules[num] = match[2]
	}

	messages := []string{}
	for s.Scan() {
		messages = append(messages, s.Text())
	}

	if *flagLoops {
		rules[int64(8)] = "42 | 42 8"
		rules[int64(11)] = "42 31 | 42 11 31"
	}

	count := 0
	for _, message := range messages {
		if followsRule(message, rules[int64(0)]) {
			count++
		}
	}
	log.Printf("%d messages follow rule 0", count)
}

func followsRule(message, rule string) bool {
	valid, remainder := followsRuleRemainder(message, rule)
	return valid && remainder == ""
}

func followsRuleRemainder(message, rule string) (bool, string) {
	if quoteMatch := quoteRuleRE.FindStringSubmatch(rule); len(quoteMatch) == 2 {
		if message == "" {
			return false, ""
		}
		return quoteMatch[1] == message[:1], message[1:]
	}

	if options := strings.Split(rule, " | "); len(options) > 1 {
		for _, option := range options {
			if valid, remainder := followsRuleRemainder(message, option); valid {
				return valid, remainder
			}
		}
		return false, ""
	}

	nextRulesS := strings.Split(rule, " ")
	for i := 0; i < len(nextRulesS); i++ {
		nextRule, _ := strconv.ParseInt(nextRulesS[i], 10, 64)
		if *flagLoops && nextRule == int64(8) && i+1 < len(nextRulesS) { // rule 8 repeats rule 42 1+ times
			nextNextRule, _ := strconv.ParseInt(nextRulesS[i+1], 10, 64)
			for {
				valid, remainder := followsRuleRemainder(message, rules[int64(42)]) // rule 8 repeats rule 42 1+ times
				if valid {
					nextValid, nextRemainder := followsRuleRemainder(remainder, rules[nextNextRule])
					if nextValid {
						i++
						message = nextRemainder
						break
					}
				}
				if remainder == "" {
					return false, ""
				}
				message = remainder
			}
			continue
		}

		valid, remainder := followsRuleRemainder(message, rules[nextRule])
		if !valid {
			return false, ""
		}
		message = remainder
	}
	return true, message
}
