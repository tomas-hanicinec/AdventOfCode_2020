package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	passports, err := readPassports()
	if err != nil {
		panic(fmt.Errorf("failed to crete passports: %w", err))
	}

	validCounter := 0
	for _, passport := range passports {
		if passport.isValid() {
			validCounter++
		}
	}

	fmt.Printf("Total valid passports: %d\n", validCounter)
}

type Passport struct {
	number int
	fields map[string]*PassportField
}

func readPassports() ([]*Passport, error) {
	data, err := ioutil.ReadFile("inputs/day_04")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	lines := strings.Split(string(data), "\n")

	result := make([]*Passport, 0)
	lastNewlineIndex := -1
	passportCounter := 0
	for i, line := range lines {
		if line == "" {
			// newline, process the previous passport
			passportCounter++
			passport, err := newPassport(passportCounter, lines[lastNewlineIndex+1:i])
			if err != nil {
				return nil, fmt.Errorf("failed to create passport %d: %w", passportCounter, err)
			}
			result = append(result, passport)
			lastNewlineIndex = i
		}
	}

	return result, nil
}

func newPassport(number int, lines []string) (*Passport, error) {
	// Get passport skeleton
	passport := getEmptyPassport(number)
	// Fill fields
	for _, line := range lines {
		for _, fieldChunk := range strings.Split(line, " ") {
			field := strings.Split(fieldChunk, ":")
			err := passport.setFieldValue(strings.TrimSpace(field[0]), strings.TrimSpace(field[1]))
			if err != nil {
				return nil, fmt.Errorf("failed to set field [%s]: %w", field[0], err)
			}
		}
	}

	return passport, nil
}

func getEmptyPassport(number int) *Passport {
	fieldValidators := getValidators()
	fields := make(map[string]*PassportField)
	for code, fieldValidator := range fieldValidators {
		fields[code] = &PassportField{
			code:      code,
			validator: fieldValidator,
			value:     "",
		}
	}
	return &Passport{
		number: number,
		fields: fields,
	}
}

func (p *Passport) setFieldValue(fieldCode string, value string) error {
	if _, exists := p.fields[fieldCode]; !exists {
		return fmt.Errorf("unknown field [%s]", fieldCode)
	}

	p.fields[fieldCode].value = value
	return nil
}

func (p *Passport) isValid() bool {
	for _, field := range p.fields {
		if err := field.validate(); err != nil {
			fmt.Printf("Passport #%d not valid. Field [%s] validation error: %s\n", p.number, field.code, err)
			return false
		}
	}

	return true
}

type PassportField struct {
	code      string
	value     string
	validator Validator
}

func (pf PassportField) validate() error {
	return pf.validator(pf.value)
}

type Validator func(val string) error

func getValidators() map[string]Validator {
	return map[string]Validator{
		"byr": func(val string) error {
			// (Birth Year)
			if err := mandatoryValidator(val); err != nil {
				return err
			}
			return numericRangeValidator(val, 1920, 2002)
		},
		"iyr": func(val string) error {
			// (Issue Year)
			if err := mandatoryValidator(val); err != nil {
				return err
			}
			return numericRangeValidator(val, 2010, 2020)
		},
		"eyr": func(val string) error {
			// (Expiration Year)
			if err := mandatoryValidator(val); err != nil {
				return err
			}
			return numericRangeValidator(val, 2020, 2030)
		},
		"hgt": func(val string) error {
			// (Height)
			if err := mandatoryValidator(val); err != nil {
				return err
			}
			pattern := regexp.MustCompile("^([0-9]+)(cm|in)$")
			matches := pattern.FindStringSubmatch(val)
			if len(matches) != 3 {
				return fmt.Errorf("height [%s] in invalid format", val)
			}
			if matches[2] == "cm" {
				return numericRangeValidator(matches[1], 150, 193)
			} else {
				return numericRangeValidator(matches[1], 59, 76)
			}
		},
		"hcl": func(val string) error {
			// (Hair Color)
			if err := mandatoryValidator(val); err != nil {
				return err
			}
			return regexpValidator(val, "^#[0-9a-f]{6}$")
		},
		"ecl": func(val string) error {
			// (Eye Color)
			if err := mandatoryValidator(val); err != nil {
				return err
			}
			return regexpValidator(val, "^amb|blu|brn|gry|grn|hzl|oth$")
		},
		"pid": func(val string) error {
			// (Passport ID)
			if err := mandatoryValidator(val); err != nil {
				return err
			}
			return regexpValidator(val, "^[0-9]{9}$")
		},
		"cid": func(val string) error {
			// (Country ID)
			return nil // not mandatory
		},
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
