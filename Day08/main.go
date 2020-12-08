package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type step struct {
	op  string
	num int64
}

const (
	ACC = "acc"
	JMP = "jmp"
	NOP = "nop"
)

var (
	stepRE = regexp.MustCompile(fmt.Sprintf(`(%s|%s|%s) ([+-][0-9]*)`, ACC, JMP, NOP))

	flagFixInfiniteLoop = flag.Bool("fix-infinite-loop", false, "change exactly one instruction to avoid an infinite loop")
)

func main() {
	flag.Parse()

	program, err := loadProgram(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error loading program: %v", err)
	}
	if *flagFixInfiniteLoop {
		for i, step := range program {
			switch step.op {
			case NOP:
				tmp := step
				tmp.op = JMP
				program[i] = tmp
				acc, err := executeProgram(program, true)
				if err == nil {
					log.Printf("fixed step %d, program result: %d", i, acc)
					return
				}
				program[i] = step
			case JMP:
				tmp := step
				tmp.op = NOP
				program[i] = tmp
				acc, err := executeProgram(program, true)
				if err == nil {
					log.Printf("fixed step %d, program result: %d", i, acc)
					return
				}
				program[i] = step
			}
		}
	}

	acc, err := executeProgram(program, false)
	if err != nil {
		log.Fatalf("error executing program: %v", err)
	}
	log.Printf("program result: %d", acc)
}

func loadProgram(filename string) ([]step, error) {
	program := []step{}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		match := stepRE.FindStringSubmatch(s.Text())
		if len(match) != 3 {
			return nil, fmt.Errorf("cannot load improperly formatted step '%s'", s.Text())
		}
		n, err := strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing num part '%s' as int: %v", match[2], err)
		}
		program = append(program, step{match[1], n})
	}
	return program, nil
}

func executeProgram(program []step, infiniteErr bool) (int64, error) {
	acc, stepN := int64(0), int64(0)
	l := len(program)
	encountered := make([]bool, l)

	for stepN < int64(l) {
		if encountered[stepN] {
			if infiniteErr {
				return -1, fmt.Errorf("infinite loop")
			}
			break
		}

		encountered[stepN] = true
		step := program[stepN]
		switch step.op {
		case ACC:
			acc += step.num
			stepN++
		case JMP:
			stepN += step.num
		case NOP:
			stepN++
		default:
			return -1, fmt.Errorf("encountered unsupported operation '%s'", step.op)
		}
	}

	return acc, nil
}
