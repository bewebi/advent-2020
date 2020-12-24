package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	flagInput    = flag.String("input", "219748365", "game input")
	flagMoves    = flag.Int("moves", 100, "number of moves")
	flagCups     = flag.Int("cups", 9, "number of cups")
	flagAfterOne = flag.Bool("after-one", false, "print product of cups after 1 instead of all cups")

	cups = map[int]*cupNode{}
)

type cupNode struct {
	val  int
	next *cupNode
}

func main() {
	flag.Parse()

	var prev, head *cupNode
	for _, c := range []rune(*flagInput) {
		if prev == nil {
			headNode := cupNode{int(c) - '0', nil}
			head = &headNode
			prev = head
		} else {
			newCup := cupNode{int(c) - '0', nil}
			prev.next = &newCup
			prev = &newCup
		}
		cups[prev.val] = prev
	}

	for i := len(*flagInput) + 1; i <= *flagCups; i++ {
		newCup := cupNode{i, nil}
		prev.next = &newCup
		prev = &newCup
		cups[prev.val] = prev
	}
	prev.next = head

	for i := 0; i < *flagMoves; i++ {
		head = doMove(head)
	}
	if *flagCups < 10 {
		log.Printf("after %d moves cups after '1' are %s", *flagMoves, cupsString(head))
	}
	log.Printf("two after 1 are %d, %d, which multiply to %d", cups[1].next.val, cups[1].next.next.val, cups[1].next.val*cups[1].next.next.val)
}

func doMove(cup *cupNode) *cupNode {
	held := map[int]bool{cup.next.val: true, cup.next.next.val: true, cup.next.next.next.val: true}
	nextVal := cup.val
	subSection := cup.next
	cup.next = cup.next.next.next.next

	for {
		nextVal--
		if nextVal <= 0 {
			nextVal = *flagCups
		}

		if !held[nextVal] {
			targetCup := cups[nextVal]
			subSection.next.next.next = targetCup.next
			targetCup.next = subSection
			return cup.next
		}
	}
}

func cupsString(cup *cupNode) string {
	foundOne := false
	ret := ""
	for {
		if cup.val == 1 {
			if foundOne {
				break
			}
			foundOne = true
		}
		if foundOne && cup.val != 1 {
			ret += fmt.Sprintf("%d", cup.val)
		}
		cup = cup.next
	}
	return ret
}
