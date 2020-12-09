package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
)

type sumIndices struct {
	a, b int
}

type sumMap map[int64][]*sumIndices

func main() {
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	numList := []int64{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		n, err := strconv.ParseInt(s.Text(), 10, 64)
		if err != nil {
			log.Fatalf("cannot parse line %s as int", s.Text())
		}
		numList = append(numList, n)
	}
	sm := sumMap{}
	for i := 0; i < 25; i++ {
		sm.add25Sums(numList, i)
	}

	nonSum := int64(0)
	for i := 25; i < len(numList); i++ {
		if si, ok := sm[numList[i]]; !ok || si == nil {
			nonSum = numList[i]
			break
		}
		sm.remove25Sums(i - 25)
		sm.add25Sums(numList, i)
	}
	log.Printf("the first nonSum number is %d", nonSum)
	startIndex, endIndex := findContiguousSumList(numList, nonSum)
	if startIndex < 0 || endIndex < 0 {
		log.Fatalf("did not find contiguous sum list")
	}
	log.Printf("the sequence in the list from indexes %d to %d sum to %d", startIndex, endIndex, nonSum)
	min, max := minMax(numList[startIndex:endIndex])
	log.Printf("the min in that sequence is %d, the max is %d, and they sum to %d", min, max, min+max)
}

func (sm sumMap) add25Sums(numList []int64, index int) {
	floor := index - 24 // passed index is one of the 25
	if floor < 0 {
		floor = 0
	}
	for i := floor; i < index; i++ {
		sum := numList[i] + numList[index]
		if _, ok := sm[sum]; !ok {
			sm[sum] = []*sumIndices{}
		}
		sm[sum] = append(sm[sum], &sumIndices{i, index})
	}
}

func (sm sumMap) remove25Sums(index int) {
	for key, sis := range sm {
		for i, si := range sis {
			if si.a == index || si.b == index {
				sm[key] = append(sis[:i], sis[i+1:]...)
			}
		}
	}
}

func findContiguousSumList(numList []int64, sum int64) (int, int) {
	sums := make([]int64, len(numList))
	for i := 0; i < len(numList); i++ {
		for startIndex, sumSoFar := range sums {
			if sumSoFar < 0 || startIndex >= i {
				continue
			}
			if sumSoFar+numList[i] == sum {
				return startIndex, i
			}
			if sumSoFar+numList[i] > sum {
				sums[startIndex] = -1
			}
			if newSum := sumSoFar + numList[i]; newSum < sum {
				sums[startIndex] = newSum
			}
		}
		sums[i] = numList[i]
	}
	return -1, -1
}

func minMax(numList []int64) (int64, int64) {
	min, max := int64(math.MaxInt64), int64(math.MinInt64)
	for _, n := range numList {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return min, max
}
