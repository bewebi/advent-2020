package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type passport struct {
	byr, iyr, eyr, hgt, hcl, ecl, pid, cid *string
}
type field struct {
	key, val string
}

func main() {
	flag.Parse()

	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	passports := []*passport{}
	fields := []field{}
	validCnt := 0

	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() == "" {
			pass, err := newPassportFromParts(fields)
			if err != nil {
				log.Fatalf("error creating passport from parts: %v", err)
			}
			passports = append(passports, pass)
			if pass.isValid() {
				validCnt++
			}
			fields = []field{}
			continue
		}
		parts := strings.Split(s.Text(), " ")
		for _, p := range parts {
			kv := strings.Split(p, ":")
			if len(kv) != 2 {
				log.Fatalf("part %s is malformed", p)
			}
			fields = append(fields, field{kv[0], kv[1]})
		}
	}
	pass, err := newPassportFromParts(fields)
	if err != nil {
		log.Fatalf("error creating passport from parts: %v", err)
	}
	passports = append(passports, pass)
	if pass.isValid() {
		validCnt++
	}

	log.Printf("%d passports are valid of %d total", validCnt, len(passports))
}

func newPassportFromParts(fields []field) (*passport, error) {
	p := &passport{}

	for _, f := range fields {
		switch f.key {
		case "byr":
			p.byr = &f.val
		case "iyr":
			p.iyr = &f.val
		case "eyr":
			p.eyr = &f.val
		case "hgt":
			p.hgt = &f.val
		case "hcl":
			p.hcl = &f.val
		case "ecl":
			p.ecl = &f.val
		case "pid":
			p.pid = &f.val
		case "cid":
			p.cid = &f.val
		default:
			return nil, fmt.Errorf("invalid part key: %s", f.key)
		}
	}
	return p, nil
}

func (p *passport) isValid() bool {
	return p.byr != nil && p.iyr != nil && p.eyr != nil && p.hgt != nil && p.hcl != nil && p.ecl != nil && p.pid != nil
}
