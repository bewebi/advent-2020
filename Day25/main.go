package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	s.Scan()
	cardKeyString := s.Text()
	cardKey, _ := strconv.ParseInt(cardKeyString, 10, 64)
	cardSecretLoop := findLoopSize(7, cardKey)
	log.Printf("card secret loop: %d", cardSecretLoop)

	s.Scan()
	doorKeyString := s.Text()
	doorKey, _ := strconv.ParseInt(doorKeyString, 10, 64)
	doorSecretLoop := findLoopSize(7, doorKey)
	log.Printf("door secret loop: %d", doorSecretLoop)

	doorEnc := generateEncryptionKey(doorKey, cardSecretLoop)
	cardEnc := generateEncryptionKey(cardKey, doorSecretLoop)

	log.Printf("the encrption key generated from the door key is %d, generated from the card key is %d", doorEnc, cardEnc)
}

func findLoopSize(subject, target int64) int64 {
	curKey := int64(1)
	for i := int64(1); true; i++ {
		curKey *= subject
		curKey = curKey % 20201227

		if curKey == target {
			return i
		}
	}
	return -1
}

func generateEncryptionKey(subject, loopSize int64) int64 {
	curKey := int64(1)
	for i := int64(0); i < loopSize; i++ {
		curKey *= subject
		curKey = curKey % 20201227
	}
	return curKey
}
