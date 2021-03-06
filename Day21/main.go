package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	ingredientsRE = regexp.MustCompile(`^([a-z ]+) \(contains ([a-z, ]+)\)$`)

	allergenList          = []string{}
	ingredientList        = []string{}
	allergenPossabilities = map[string][]string{}
	suspectedIngredients  = map[string][]string{}
	ingredientCount       = map[string]int{}
)

func main() {
	f, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("error opening input file: %v", err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() {
		match := ingredientsRE.FindStringSubmatch(s.Text())
		ingredients := strings.Split(match[1], " ")
		allergens := strings.Split(match[2], ", ")
		for _, ingredient := range ingredients {
			ingredientCount[ingredient]++
			if _, ok := suspectedIngredients[ingredient]; !ok {
				suspectedIngredients[ingredient] = []string{}
			}
		}

		for _, allergen := range allergens {
			if is, ok := allergenPossabilities[allergen]; ok {
				filteredIngredients := []string{}
				for _, i := range is {
					if inStrings(i, ingredients) {
						filteredIngredients = append(filteredIngredients, i)
					} else {
						suspectedIngredients[i] = removeString(allergen, suspectedIngredients[i])
					}
				}
				allergenPossabilities[allergen] = filteredIngredients
			} else {
				allergenPossabilities[allergen] = ingredients
				for _, ingredient := range ingredients {
					if sis, ok := suspectedIngredients[ingredient]; ok {
						suspectedIngredients[ingredient] = append(sis, allergen)
					} else {
						suspectedIngredients[ingredient] = []string{allergen}
					}
				}
			}
		}
	}

	nonAllergenCount := 0
	for ingredient, allergens := range suspectedIngredients {
		if len(allergens) == 0 {
			nonAllergenCount += ingredientCount[ingredient]
		}
	}
	log.Printf("ingredients that definitely have allergens appear %d times", nonAllergenCount)

	for {
		done := true
		for allergen, possabilities := range allergenPossabilities {
			if len(possabilities) == 1 {
				for a, ps := range allergenPossabilities {
					if a != allergen {
						allergenPossabilities[a] = removeString(possabilities[0], ps)
					}
				}
			} else {
				done = false
			}
		}
		if done {
			break
		}
	}
	for allergen := range allergenPossabilities {
		allergenList = append(allergenList, allergen)
	}
	sort.Strings(allergenList)

	for _, allergen := range allergenList {
		ingredientList = append(ingredientList, allergenPossabilities[allergen][0])
	}
	log.Printf("list of allergenic ingredients alphabetical by allergen: %s", strings.Join(ingredientList, ","))
}

func inStrings(s string, ss []string) bool {
	for _, elt := range ss {
		if s == elt {
			return true
		}
	}
	return false
}

func removeString(s string, ss []string) []string {
	for i, elt := range ss {
		if s == elt {
			return append(ss[:i], ss[i+1:]...)
		}
	}
	return ss
}
