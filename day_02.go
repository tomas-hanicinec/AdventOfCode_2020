package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	input, err := getPasswords()
	if err != nil {
		panic(fmt.Errorf("failed to get input: %w", err))
	}

	validCounter := 0
	for _, val := range input {
		if val.isValidV2() {
			validCounter++
		}
	}

	fmt.Printf("%d out of %d passwords are valid\n", validCounter, len(input))
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

func newPassword(inputLine string) (*Password, error) {
	pattern := regexp.MustCompile("^([0-9]+)-([0-9]+) ([a-z]): ([a-z]+)$")
	matches := pattern.FindStringSubmatch(inputLine)
	if len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse input line")
	}

	a, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid first value [%s]: %w", matches[1], err)
	}
	b, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid second value [%s]: %w", matches[2], err)
	}

	return &Password{
		policy: PasswordPolicy{
			a:      a,
			b:      b,
			letter: matches[3][0],
		},
		password: matches[4],
	}, nil
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
		fmt.Printf("INVALID pass [%s] has %d letters [%s], can have %d-%d\n", p.password, count, string(p.policy.letter), p.policy.a, p.policy.b)
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
		fmt.Printf("INVALID pass [%s] has %d letters [%s] on positions %d,%d\n", p.password, count, string(p.policy.letter), p.policy.a, p.policy.b)
	}
	return result
}

func getPasswords() ([]*Password, error) {
	lines, err := getInputLines()
	if err != nil {
		return nil, err
	}
	result := make([]*Password, len(lines))
	for i, line := range lines {
		pp, err := newPassword(line)
		if err != nil {
			return nil, fmt.Errorf("failed to create password policy from [%s]: %w", line, err)
		}
		result[i] = pp
	}

	return result, nil
}

func getInputLines() ([]string, error) {
	data, err := ioutil.ReadFile("inputs/day_02")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return strings.Split(string(data), "\n"), nil
}
