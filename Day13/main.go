package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	s.Scan()
	start, err := strconv.ParseInt(s.Text(), 10, 64)
	if err != nil {
		log.Fatalf("error parsing first input line as int: %v", err)
	}

	s.Scan()
	busses := strings.Split(s.Text(), ",")
	busIDs := make([]int64, len(busses))

	min, busID := int64(math.MaxInt64), int64(0)
	for i, b := range busses {
		if b == "x" {
			busIDs[i] = int64(-1)
			continue
		}
		id, err := strconv.ParseInt(b, 10, 64)
		if err != nil {
			log.Fatalf("error parsing bus %s as int: %v", b, err)
		}
		busIDs[i] = id

		timeSinceBus := start % id
		if timeSinceBus == 0 {
			log.Printf("bus %s will come right now!", b)
			return
		}
		timeToBus := id - start%id
		if timeToBus < min {
			min = timeToBus
			busID = id
		}
	}
	log.Printf("next bus after %d is %d, arriving in %d minutes; id*wait: %d", start, busID, min, min*busID)

	base := busIDs[0]
	found := map[int64]bool{busIDs[0]: true}
	attempt, prodFound := int64(0), int64(1)
	start = 0
	success := false
	for !success {
		attempt += prodFound
		success = true
		start = base * attempt
		for i, id := range busIDs {
			if id == -1 {
				continue
			}
			if (start+int64(i))%id != 0 {
				success = false
				break
			} else if !found[id] {
				prodFound *= id
				found[id] = true
			}
		}
	}
	log.Printf("the earliest timestamp for successive busses is %d", start)
}
