package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func main() {
	passwords := getPasswords()
	validV1Counter, validV2Counter := 0, 0
	for _, val := range passwords {
		if val.isValidV1() {
			validV1Counter++
		}
		if val.isValidV2() {
			validV2Counter++
		}
	}

	fmt.Printf("V1: %d out of %d passwords are valid\n", validV1Counter, len(passwords)) // Part I.
	fmt.Printf("V2: %d out of %d passwords are valid\n", validV2Counter, len(passwords)) // Part II.
}

type PasswordPolicy struct {
	a      int
	b      int
	letter byte
}

type Password struct {
	policy   PasswordPolicy
	password string
}

func newPassword(inputLine string) *Password {
	pattern := regexp.MustCompile("^([0-9]+)-([0-9]+) ([a-z]): ([a-z]+)$")
	matches := pattern.FindStringSubmatch(inputLine)
	if len(matches) < 5 {
		panic(fmt.Errorf("failed to parse input line [%s]", inputLine))
	}

	a, err := strconv.Atoi(matches[1])
	if err != nil {
		panic(fmt.Errorf("invalid first value [%s]: %w", matches[1], err))
	}
	b, err := strconv.Atoi(matches[2])
	if err != nil {
		panic(fmt.Errorf("invalid second value [%s]: %w", matches[2], err))
	}

	return &Password{
		policy: PasswordPolicy{
			a:      a,
			b:      b,
			letter: matches[3][0],
		},
		password: matches[4],
	}
}

func (p Password) isValidV1() bool {
	count := 0
	for i := range p.password {
		if p.password[i] == p.policy.letter {
			count++
		}
	}

	result := count >= p.policy.a && count <= p.policy.b
	if !result {
		//fmt.Printf("INVALID pass [%s] has %d letters [%s], can have %d-%d\n", p.password, count, string(p.policy.letter), p.policy.a, p.policy.b)
	}
	return result
}

func (p Password) isValidV2() bool {
	count := 0
	if p.password[p.policy.a-1] == p.policy.letter {
		count++
	}
	if p.password[p.policy.b-1] == p.policy.letter {
		count++
	}
	result := count == 1
	if !result {
		//fmt.Printf("INVALID pass [%s] has %d letters [%s] on positions %d,%d\n", p.password, count, string(p.policy.letter), p.policy.a, p.policy.b)
	}
	return result
}

func getPasswords() []*Password {
	lines := ReadLines("inputs/day_02.txt")
	result := make([]*Password, len(lines))
	for i, line := range lines {
		result[i] = newPassword(line)
	}

	return result
}
