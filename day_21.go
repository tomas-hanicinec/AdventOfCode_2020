package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	foodList := getFoodList()
	allergenInIngredient, ingredientHasAllergen := foodList.getAllergenMapping()

	// Part I.
	counter := 0
	for _, food := range foodList.foods {
		for _, ingredient := range food.ingredients {
			if _, ok := ingredientHasAllergen[ingredient]; !ok {
				counter++
			}
		}
	}
	fmt.Printf("Sum of occurences of ingredients without allergens: %d\n", counter)

	// Part II.
	allergens := make([]string, 0)
	for allergen := range allergenInIngredient {
		allergens = append(allergens, allergen)
	}
	sort.Strings(allergens)
	result := make([]string, len(allergens))
	for i, allergen := range allergens {
		result[i] = allergenInIngredient[allergen]
	}
	fmt.Printf("Sorted list of ingredients with allergens: [%s]\n", strings.Join(result, ","))
}

func getFoodList() FoodList {
	lines := ReadLines("inputs/day_21.txt")
	foods := make([]Food, len(lines))
	allergenInFood := make(map[string][]int)
	for i := range lines {
		foods[i] = NewFood(lines[i])
		for _, al := range foods[i].allergens {
			if _, ok := allergenInFood[al]; ok {
				allergenInFood[al] = append(allergenInFood[al], i)
			} else {
				allergenInFood[al] = []int{i}
			}
		}
	}

	return FoodList{
		foods:          foods,
		allergenInFood: allergenInFood,
	}
}

type FoodList struct {
	allergenInFood map[string][]int
	foods          []Food
}

func (fl FoodList) getAllergenMapping() (map[string]string, map[string]string) {
	alToIng, ingToAl := make(map[string]string), make(map[string]string)
	counter := 1
	for len(alToIng) < len(fl.allergenInFood) {
		for allergen, foodIndexes := range fl.allergenInFood {
			commonIngredients := make([]string, 0)
			for _, ingredient := range fl.getCommonIngredients(foodIndexes) {
				if _, ok := ingToAl[ingredient]; ok {
					continue // this ingredient already has an assigned allergen
				}
				commonIngredients = append(commonIngredients, ingredient)
			}
			if len(commonIngredients) == 1 {
				// unique match for allergen found
				alToIng[allergen] = commonIngredients[0]
				ingToAl[commonIngredients[0]] = allergen
			}
		}
		if counter > len(fl.allergenInFood) {
			panic(fmt.Errorf("no unique allergen mapping found")) // to prevent infinite loop
		}
		counter++
	}

	return alToIng, ingToAl
}

func (fl FoodList) getCommonIngredients(foodIndexes []int) []string {
	countMap := make(map[string]int)
	for _, i := range foodIndexes {
		for _, ingredient := range fl.foods[i].ingredients {
			if _, ok := countMap[ingredient]; ok {
				countMap[ingredient]++
			} else {
				countMap[ingredient] = 1
			}
		}
	}

	result := make([]string, 0)
	for ingredient, count := range countMap {
		if count == len(foodIndexes) {
			result = append(result, ingredient)
		}
	}

	return result
}

type Food struct {
	ingredients []string
	allergens   []string
}

func NewFood(line string) Food {
	parts := strings.Split(line, " (contains ")
	allergensString := parts[1][:len(parts[1])-1] // strip the trailing ")"
	return Food{
		ingredients: strings.Split(strings.TrimSpace(parts[0]), " "),
		allergens:   strings.Split(allergensString, ", "),
	}
}
