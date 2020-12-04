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

type passport struct {
	byr, iyr, eyr, hgt, hcl, ecl, pid, cid *string
	allFields, allValid                    bool
}
type field struct {
	key, val string
}

var (
	colorRE = regexp.MustCompile("^#[0-9a-f]{6}$")
	idRE    = regexp.MustCompile("^[0-9]{9}$")

	flagValidateFields = flag.Bool("validate-fields", false, "check validity of fields")
)

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
			if *flagValidateFields {
				if pass.allValid {
					validCnt++
				}
			} else {
				if pass.allFields {
					validCnt++
				}
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
	if *flagValidateFields {
		if pass.allValid {
			validCnt++
		}
	} else {
		if pass.allFields {
			validCnt++
		}
	}

	log.Printf("%d passports are valid of %d total", validCnt, len(passports))
}

func newPassportFromParts(fields []field) (*passport, error) {
	p := &passport{}
	allValid := true

	for _, f := range fields {
		switch f.key {
		case "byr":
			p.byr = &f.val
			byr, err := strconv.ParseInt(f.val, 10, 64)
			if err != nil {
				allValid = false
				continue
			}
			if byr < 1920 || byr > 2002 {
				allValid = false
			}
		case "iyr":
			p.iyr = &f.val
			iyr, err := strconv.ParseInt(f.val, 10, 64)
			if err != nil {
				allValid = false
				continue
			}
			if iyr < 2010 || iyr > 2020 {
				allValid = false
			}
		case "eyr":
			p.eyr = &f.val
			eyr, err := strconv.ParseInt(f.val, 10, 64)
			if err != nil {
				allValid = false
				continue
			}
			if eyr < 2020 || eyr > 2030 {
				allValid = false
			}
		case "hgt":
			p.hgt = &f.val
			if i := strings.Index(f.val, "cm"); i >= 0 {
				hgt, err := strconv.ParseInt(f.val[:i], 10, 64)
				if err != nil {
					allValid = false
					continue
				}
				if hgt < 150 || hgt > 193 {
					allValid = false
				}
			} else if i := strings.Index(f.val, "in"); i >= 0 {
				hgt, err := strconv.ParseInt(f.val[:i], 10, 64)
				if err != nil {
					allValid = false
					continue
				}
				if hgt < 59 || hgt > 76 {
					allValid = false
				}
			} else {
				allValid = false
			}
		case "hcl":
			p.hcl = &f.val
			if !colorRE.Match([]byte(f.val)) {
				allValid = false
			}
		case "ecl":
			p.ecl = &f.val
			switch f.val {
			case "amb", "blu", "brn", "gry", "grn", "hzl", "oth":
			default:
				allValid = false
			}
		case "pid":
			p.pid = &f.val
			if !idRE.Match([]byte(f.val)) {
				allValid = false
			}
		case "cid":
			p.cid = &f.val
		default:
			return nil, fmt.Errorf("invalid part key: %s", f.key)
		}
	}
	p.allFields = p.byr != nil && p.iyr != nil && p.eyr != nil && p.hgt != nil && p.hcl != nil && p.ecl != nil && p.pid != nil
	p.allValid = p.allFields && allValid

	return p, nil
}
