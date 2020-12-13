package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	passports := readPassports()
	validCounterV1, validCounterV2 := 0, 0
	for _, passport := range passports {
		if passport.isValid(ValidateV1) {
			validCounterV1++
		}
		if passport.isValid(ValidateV2) {
			validCounterV2++
		}
	}

	fmt.Printf("V1 total valid passports: %d\n", validCounterV1) // Part I.
	fmt.Printf("V2 total valid passports: %d\n", validCounterV2) // Part II.
}

type Passport struct {
	number int
	fields map[string]string
}

func readPassports() []*Passport {
	lines := ReadLines("inputs/day_04.txt")

	result := make([]*Passport, 0)
	lastNewlineIndex := -1
	passportCounter := 0
	for i, line := range lines {
		if line == "" {
			// newline, process the previous passport
			passportCounter++
			result = append(result, newPassport(passportCounter, lines[lastNewlineIndex+1:i]))
			lastNewlineIndex = i
		}
	}

	return result
}

func newPassport(number int, lines []string) *Passport {
	// Get passport skeleton
	passport := getEmptyPassport(number)
	// Fill fields from input
	for _, line := range lines {
		for _, fieldChunk := range strings.Split(line, " ") {
			field := strings.Split(fieldChunk, ":")
			passport.fields[strings.TrimSpace(field[0])] = strings.TrimSpace(field[1])
		}
	}

	return passport
}

func getEmptyPassport(number int) *Passport {
	return &Passport{
		number: number,
		fields: map[string]string{
			"byr": "",
			"iyr": "",
			"eyr": "",
			"hgt": "",
			"hcl": "",
			"ecl": "",
			"pid": "",
			"cid": "",
		},
	}
}

func (p *Passport) setField(code string, value string) {
	p.fields[code] = value
}

func (p *Passport) isValid(validator Validator) bool {
	for code, value := range p.fields {
		if err := validator(code, value); err != nil {
			//fmt.Printf("Passport #%d not valid. Field [%s] validation error: %s\n", p.number, code, err)
			return false
		}
	}

	return true
}

type Validator func(string, string) error

func ValidateV1(fieldCode string, fieldValue string) error {
	if fieldCode != "cid" {
		return mandatoryValidator(fieldValue)
	}
	return nil
}

func ValidateV2(fieldCode string, fieldValue string) error {
	switch fieldCode {

	case "byr":
		// (Birth Year)
		if err := mandatoryValidator(fieldValue); err != nil {
			return err
		}
		return numericRangeValidator(fieldValue, 1920, 2002)

	case "iyr":
		// (Issue Year)
		if err := mandatoryValidator(fieldValue); err != nil {
			return err
		}
		return numericRangeValidator(fieldValue, 2010, 2020)
	case "eyr":
		// (Expiration Year)
		if err := mandatoryValidator(fieldValue); err != nil {
			return err
		}
		return numericRangeValidator(fieldValue, 2020, 2030)
	case "hgt":
		// (Height)
		if err := mandatoryValidator(fieldValue); err != nil {
			return err
		}
		pattern := regexp.MustCompile("^([0-9]+)(cm|in)$")
		matches := pattern.FindStringSubmatch(fieldValue)
		if len(matches) != 3 {
			return fmt.Errorf("height [%s] in invalid format", fieldValue)
		}
		if matches[2] == "cm" {
			return numericRangeValidator(matches[1], 150, 193)
		} else {
			return numericRangeValidator(matches[1], 59, 76)
		}
	case "hcl":
		// (Hair Color)
		if err := mandatoryValidator(fieldValue); err != nil {
			return err
		}
		return regexpValidator(fieldValue, "^#[0-9a-f]{6}$")
	case "ecl":
		// (Eye Color)
		if err := mandatoryValidator(fieldValue); err != nil {
			return err
		}
		return regexpValidator(fieldValue, "^amb|blu|brn|gry|grn|hzl|oth$")
	case "pid":
		// (Passport ID)
		if err := mandatoryValidator(fieldValue); err != nil {
			return err
		}
		return regexpValidator(fieldValue, "^[0-9]{9}$")
	default:
		return nil // (Country ID - "cid") is not mandatory, unsupported fields are not important
	}
}

func mandatoryValidator(val string) error {
	if val == "" {
		return fmt.Errorf("empty mandatory value")
	}
	return nil
}

func numericRangeValidator(val string, min int, max int) error {
	numVal, err := strconv.Atoi(val)
	if err != nil {
		return fmt.Errorf("value is not a number")
	}

	if numVal < min || numVal > max {
		return fmt.Errorf("value [%d] outside of allowed range [%d-%d]", numVal, min, max)
	}

	return nil
}

func regexpValidator(val string, pattern string) error {
	patternObject := regexp.MustCompile(pattern)
	if !patternObject.MatchString(val) {
		return fmt.Errorf("value [%s] not satysfying pattern [%s]", val, pattern)
	}

	return nil
}
