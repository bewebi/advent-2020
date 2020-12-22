package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	flagRecurse = flag.Bool("recurse", false, "play a recursive game")

	player1 = []int64{}
	player2 = []int64{}
)

func main() {
	flag.Parse()

	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	s.Scan()
	for s.Scan() && s.Text() != "" {
		card, _ := strconv.ParseInt(s.Text(), 10, 64)
		player1 = append(player1, card)
	}

	s.Scan()
	for s.Scan() && s.Text() != "" {
		card, _ := strconv.ParseInt(s.Text(), 10, 64)
		player2 = append(player2, card)
	}

	if *flagRecurse {
		infinite := false
		player1, player2, infinite = playRecursiveGame()
		if infinite {
			player2 = []int64{}
		}
	} else {
		playGame()
	}

	winningScore := int64(0)
	if len(player2) > 0 {
		for i := 1; i <= len(player2); i++ {
			winningScore += int64(i) * player2[len(player2)-i]
		}
	} else {
		for i := 1; i <= len(player1); i++ {
			winningScore += int64(i) * player1[len(player1)-i]
		}
	}
	log.Printf("the winning score is %d", winningScore)
}

func playGame() {
	for len(player1) > 0 && len(player2) > 0 {
		if player1[0] > player2[0] {
			player1 = append(player1, player1[0])
			player1 = append(player1, player2[0])
		} else {
			player2 = append(player2, player2[0])
			player2 = append(player2, player1[0])
		}
		player1 = player1[1:]
		player2 = player2[1:]
	}
}

func playRecursiveGame() ([]int64, []int64, bool) {
	return recursiveGame(player1, player2)
}

func recursiveGame(player1, player2 []int64) ([]int64, []int64, bool) {
	previousP1Rounds := map[string]bool{}
	for len(player1) > 0 && len(player2) > 0 {
		p1Winner, recursed := false, false

		if player1[0] < int64(len(player1)) && player2[0] < int64(len(player2)) {
			recursed = true
			p1, _, infinite := recursiveGame(copy(player1[1:player1[0]+1]), copy(player2[1:player2[0]+1]))
			p1Winner = infinite || (len(p1) > 0)
		}
		if (recursed && p1Winner) || (!recursed && player1[0] > player2[0]) {
			player1 = append(player1, player1[0])
			player1 = append(player1, player2[0])
		} else {
			player2 = append(player2, player2[0])
			player2 = append(player2, player1[0])
		}
		player1 = player1[1:]
		player2 = player2[1:]

		if previousP1Rounds[toString(player1)] {
			return player1, player2, true
		}
		previousP1Rounds[toString(player1)] = true
	}
	return player1, player2, false
}

func copy(in []int64) []int64 {
	out := make([]int64, len(in))
	for i, n := range in {
		out[i] = n
	}
	return out
}

func toString(ns []int64) string {
	if len(ns) == 0 {
		return ""
	}
	s := fmt.Sprintf("%d", ns[0])
	for i := 1; i < len(ns); i++ {
		s = fmt.Sprintf("%s-%d", s, ns[i])
	}
	return s
}
