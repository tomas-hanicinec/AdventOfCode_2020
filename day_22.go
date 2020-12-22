package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Part I.
	game := getGame()
	winner := game.playV1()
	fmt.Printf("V1 game winner is Player [%d] with the score od %d\n", winner, game.decks[winner-1].getScore())

	// Part II.
	game = getGame()
	winner = game.playV2()
	fmt.Printf("V2 game winner is Player [%d] with the score od %d\n", winner, game.decks[winner-1].getScore())
}

func getGame() Game {
	lines := ReadLines("inputs/day_22.txt")
	return NewGame(getDeck("Player 1:", lines), getDeck("Player 2:", lines))
}

func getDeck(headline string, lines []string) Deck {
	deckStart := 0
	for i, line := range lines {
		if line == headline {
			deckStart = i + 1
		}
		if deckStart > 0 && line == "" {
			return NewDeck(lines[deckStart:i])
		}
	}
	panic(fmt.Errorf("invalid input, deck end not reached"))
}

type Game struct {
	decks  [2]Deck
	memory [2]DeckMemory
}

func NewGame(deck1 Deck, deck2 Deck) Game {
	return Game{
		decks:  [2]Deck{deck1, deck2},
		memory: [2]DeckMemory{make(DeckMemory), make(DeckMemory)},
	}
}

const noWinner = 0

func (g *Game) playV1() int {
	for g.getGameWinner() == noWinner {
		g.playRound()
	}
	return g.getGameWinner()
}

func (g *Game) playV2() int {
	for g.getGameWinner() == noWinner {
		if g.isInMemory() {
			return 1 // loop in this game -> player 1 wins
		}
		g.saveToMemory()
		g.playRoundV2()
	}
	return g.getGameWinner()
}

func (g *Game) playRound() {
	cards := [2]int{g.decks[0].get(), g.decks[1].get()}
	winner := g.getRoundWinnerV1(cards)
	g.processRoundWinner(winner, cards)
}

func (g *Game) playRoundV2() {
	cards := [2]int{g.decks[0].get(), g.decks[1].get()}

	var winner int
	if cards[0] > len(g.decks[0].queue) || cards[1] > len(g.decks[1].queue) {
		// someone has not enough cards, play according to the v1 rules
		winner = g.getRoundWinnerV1(cards)
	} else {
		//play a sub-game
		subGame := g.getSubGame(cards)
		winner = subGame.playV2()
	}

	g.processRoundWinner(winner, cards)
}

func (g *Game) processRoundWinner(winner int, cards [2]int) {
	g.decks[winner-1].add(cards[winner-1]) // first the winners card
	g.decks[winner-1].add(cards[winner%2])
}

func (g *Game) getRoundWinnerV1(cards [2]int) int {
	if cards[0] > cards[1] {
		return 1
	} else if cards[1] > cards[0] {
		return 2
	}
	panic(fmt.Errorf("invalid deck - values [%d, %d] are the same", cards[0], cards[1])) // this is not specified in the rules
}

func (g *Game) getSubGame(cards [2]int) Game {
	queue1 := make([]int, cards[0])
	copy(queue1, g.decks[0].queue[0:cards[0]])
	queue2 := make([]int, cards[1])
	copy(queue2, g.decks[1].queue[0:cards[1]])
	return NewGame(Deck{queue: queue1}, Deck{queue: queue2})
}

func (g *Game) getGameWinner() int {
	if g.decks[0].isEmpty() {
		return 2
	}
	if g.decks[1].isEmpty() {
		return 1
	}
	return noWinner
}

func (g *Game) saveToMemory() {
	g.memory[0][g.decks[0].getFingerprint()] = struct{}{}
	g.memory[1][g.decks[1].getFingerprint()] = struct{}{}
}

func (g *Game) isInMemory() bool {
	_, ok1 := g.memory[0][g.decks[0].getFingerprint()]
	_, ok2 := g.memory[0][g.decks[0].getFingerprint()]
	return ok1 || ok2
}

type DeckMemory map[string]struct{}

type Deck struct {
	queue []int
}

func NewDeck(lines []string) Deck {
	return Deck{
		queue: StringsToInts(lines),
	}
}

func (d *Deck) add(card int) {
	d.queue = append(d.queue, card)
	if cap(d.queue) > 3*len(d.queue) {
		d.optimize()
	}
}

func (d *Deck) get() int {
	if len(d.queue) == 0 {
		panic("A")
	}
	value := d.queue[0]
	d.queue = d.queue[1:]
	return value
}

func (d *Deck) isEmpty() bool {
	return len(d.queue) == 0
}

func (d *Deck) getScore() int {
	score := 0
	for i := 0; i < len(d.queue); i++ {
		score += d.queue[i] * (len(d.queue) - i)
	}

	return score
}

func (d *Deck) optimize() {
	// memory optimization (slice-based queues ten to grow indefinitely in Go)
	newQueue := make([]int, len(d.queue))
	copy(newQueue, d.queue)
	d.queue = newQueue
}

func (d *Deck) getFingerprint() string {
	fingerprint := ""
	for i := range d.queue {
		fingerprint += strconv.Itoa(d.queue[i]) + "|"
	}
	return fingerprint
}
