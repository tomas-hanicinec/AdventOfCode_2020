package main

import "fmt"

const subjectNumber = 7
const keyLimit = 20201227
const maxLoopSize = 1000000000

func main() {
	cardPublicKey := 5290733
	doorPublicKey := 15231938

	cardLoopSize := getLoopSize(cardPublicKey)
	doorLoopSize := getLoopSize(doorPublicKey)
	encryptionKey := loop(doorPublicKey, cardLoopSize)
	if encryptionKey != loop(cardPublicKey, doorLoopSize) {
		panic(fmt.Errorf("invalid loop sizes [%d, %d], door and card encryption keys do not match", cardLoopSize, doorLoopSize))
	}
	fmt.Printf("Loop sizes: [%d, %d], encryption key: %d\n", cardLoopSize, doorLoopSize, encryptionKey)
}

func getLoopSize(publicKey int) int {
	current := 1
	loopCounter := 0
	for loopCounter < maxLoopSize {
		if current == publicKey {
			return loopCounter
		}
		loopCounter++
		current *= subjectNumber
		if current >= keyLimit {
			current = current % keyLimit
		}
	}

	panic(fmt.Errorf("loop size not found (limit of %d reached)", maxLoopSize))
}

func loop(publicKey int, loopSize int) int {
	current := 1
	for i := 0; i < loopSize; i++ {
		current *= publicKey
		if current > keyLimit {
			current = current % keyLimit
		}
	}

	return current
}
