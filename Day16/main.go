package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	fieldRE = regexp.MustCompile(`^([a-z ]+): ([0-9]+-[0-9]+) or ([0-9]+-[0-9]+)$`)
)

func main() {
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	fields := map[string][][]int64{}
	for s.Scan() && s.Text() != "" {
		match := fieldRE.FindStringSubmatch(s.Text())
		c0, c1 := make([]int64, 2), make([]int64, 2)
		c0Strings := strings.Split(match[2], "-")
		for i, str := range c0Strings {
			num, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				log.Fatalf("error parsing constraint %s as int", str)
			}
			c0[i] = num
		}
		c1Strings := strings.Split(match[3], "-")
		for i, str := range c1Strings {
			num, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				log.Fatalf("error parsing constraint %s as int", str)
			}
			c1[i] = num
		}
		fields[match[1]] = [][]int64{c0, c1}
	}

	s.Scan() // your ticket:
	s.Scan()
	myTicketStrings := strings.Split(s.Text(), ",")
	myTicket := make([]int64, len(myTicketStrings))
	for i, str := range myTicketStrings {
		num, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			log.Fatalf("error parsing ticket field %s as int", str)
		}
		myTicket[i] = num
	}

	s.Scan()
	s.Scan() // nearby tickets:
	nearbyTickets := [][]int64{}
	for s.Scan() {
		newTicketStrings := strings.Split(s.Text(), ",")
		newTicket := make([]int64, len(myTicketStrings))
		for i, str := range newTicketStrings {
			num, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				log.Fatalf("error parsing ticket field %s as int", str)
			}
			newTicket[i] = num
		}
		nearbyTickets = append(nearbyTickets, newTicket)
	}

	errorRate := int64(0)
	validTickets := [][]int64{}
	for _, ticket := range nearbyTickets {
		valid := true
		for _, field := range ticket {
			if !possibleValid(field, fields) {
				errorRate += field
				valid = false
			}
		}
		if valid {
			validTickets = append(validTickets, ticket)
		}
	}
	log.Printf("error rate of nearby tickets is %d", errorRate)

	fieldPossabilities := map[string]map[int]bool{}
	for field := range fields {
		fieldPossabilities[field] = map[int]bool{}
		for i := 0; i < len(validTickets[0]); i++ {
			fieldPossabilities[field][i] = true
		}
	}

	for _, ticket := range validTickets {
		for i, ticketField := range ticket {
			for field, constraints := range fields {
				if !((ticketField >= constraints[0][0] && ticketField <= constraints[0][1]) || (ticketField >= constraints[1][0] && ticketField <= constraints[1][1])) {
					fieldPossabilities[field][i] = false
				}
			}
		}
	}

	fieldMappingPossabilities := map[string][]int{}
	for field, possabilities := range fieldPossabilities {
		fieldMappingPossabilities[field] = []int{}
		for i, b := range possabilities {
			if b {
				fieldMappingPossabilities[field] = append(fieldMappingPossabilities[field], i)
			}
		}
	}

	fieldMappings := map[string]int{}
	for len(fieldMappings) != len(fieldMappingPossabilities) {
		for field, possabilities := range fieldMappingPossabilities {
			if len(possabilities) == 1 {
				fieldMappings[field] = possabilities[0]
				for f, ps := range fieldMappingPossabilities {
					fieldMappingPossabilities[f] = removeInt(possabilities[0], ps)
				}
			}
		}
	}

	myTicketDepartureProduct := int64(1)
	for field, mapping := range fieldMappings {
		if strings.HasPrefix(field, "departure") {
			myTicketDepartureProduct *= myTicket[mapping]
		}
	}
	log.Printf("my ticket's departure product is %d", myTicketDepartureProduct)
}

func possibleValid(i int64, fields map[string][][]int64) bool {
	for _, constraints := range fields {
		if (i > constraints[0][0] && i < constraints[0][1]) || (i > constraints[1][0] && i < constraints[1][1]) {
			return true
		}
	}
	return false
}

func removeInt(rm int, ns []int) []int {
	for i, n := range ns {
		if rm == n {
			return append(ns[:i], ns[i+1:]...)
		}
	}
	return []int{}
}
