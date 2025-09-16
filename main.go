package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

type SimulationState struct {
	PopWealth     []int // an array of the current wealth for each agent
	PopIndices    []int // an array of the indexes of Population, for Fisher-Yates
	Population    int   // the number of agents in the simulation
	Plays         int   // the number of rounds within an iteration step
	Gain          int   // how much A gives to B, in percentage of B's wealth
	Loss          int   // how much A takes from B, in percentage of B's wealth
	TotalWealth   int   // how much money each agent in the simulation begins with
	OligarchLimit int   // how much money an agent needs to be considered an oligarch
}

// Fisher-Yates is an algorithm that produces an unbiased random
// permutation of a given finite sequence of values in O(n) time.
// This uses the Durstenfeld version of the algorithm, as popularized
// by Donald Knuth in The Art of Computer Programming.
func (s *SimulationState) fisherYates() {
	for i := len(s.PopIndices) - 1; i > 0; i-- {
		r := rand.IntN(i + 1)
		s.PopIndices[i], s.PopIndices[r] = s.PopIndices[r], s.PopIndices[i]
	}
}

func (s *SimulationState) yardSaleIteration() {
	for j := 0; j < s.Plays; j++ {

		halfpeople := s.Population / 2
		// fmt.Println("Randomizing pairs...")
		// fmt.Println()
		s.fisherYates()
		for i := range halfpeople {
			i1 := s.PopIndices[i]
			i2 := s.PopIndices[i+halfpeople]
			// fmt.Printf("Pairing off agent %d against agent %d\n", i1, i2)
			win := rand.IntN(2) == 1
			v1 := s.PopWealth[i1]
			v2 := s.PopWealth[i2]
			j1 := i1
			j2 := i2
			if v2 > v1 {
				j1 = i2
				j2 = i1
			}
			// fmt.Println("Starting income...")
			// fmt.Printf("Agent 1: $%d\n", s.PopWealth[j1])
			// fmt.Printf("Agent 2: $%d\n", s.PopWealth[j2])
			if win {
				// fmt.Println("Poorer agent wins!")
				delta := s.PopWealth[j2] * s.Gain / 100
				s.PopWealth[j1] -= delta
				s.PopWealth[j2] += delta
			} else {
				// fmt.Println("Wealthier agent wins!")
				delta := s.PopWealth[j2] * s.Loss / 100
				s.PopWealth[j1] += delta
				s.PopWealth[j2] -= delta
			}
			// fmt.Println("Resulting income...")
			// fmt.Printf("Agent 1: $%d\n", s.PopWealth[j1])
			// fmt.Printf("Agent 2: $%d\n", s.PopWealth[j2])
			// fmt.Println()
		}
		// fmt.Println()
		// fmt.Printf("Loop %d of %d done...\n", j+1, s.Plays)
	}
}

func initializeSimulation(people, plays, gain, loss, init int) *SimulationState {
	s := SimulationState{
		PopWealth:   make([]int, people),
		PopIndices:  make([]int, people),
		Population:  people,
		Plays:       plays,
		Gain:        gain,
		Loss:        loss,
		TotalWealth: init,
	}
	for i := range people {
		s.PopIndices[i] = i
		s.PopWealth[i] = rand.IntN(init) + 1
		s.TotalWealth += s.PopWealth[i]
	}
	s.OligarchLimit = int(float64(s.TotalWealth) * 0.75)
	return &s
}

func main() {
	init := 1000
	p := message.NewPrinter(language.AmericanEnglish)
	s := initializeSimulation(100_000, 100, 20, 17, init)
	fmt.Println("Intitial Conditions:")
	p.Printf("Population Size: %d\n", s.Population)
	p.Printf("Possible Max Starting Wealth: $%d\n", number.Decimal(init))
	p.Printf("Total Available Wealth: $%d\n", s.TotalWealth)
	p.Printf("Plays per round: %d\n", s.Plays)
	p.Println("Note: one period (.) represents one round.")
	fmt.Printf("Starting simulation!")
	startingWealth := make([]int, s.Population)
	copy(startingWealth, s.PopWealth)
	time.Sleep(1 * time.Second)
	j := 0
	for {
		// fmt.Printf("Starting round %d...\n", j)
		if j%100 == 0 {
			fmt.Println()
		}
		fmt.Printf(".")
		// time.Sleep(1 * time.Millisecond)
		s.yardSaleIteration()
		j++
		// fmt.Println()
		for i := range s.PopWealth {
			if s.PopWealth[i] >= s.OligarchLimit {
				fmt.Printf(".")
				p.Printf("\nAfter %d rounds, Agent %d (started with $%d) has become an oligarch with $%d out of the available $%d!\n", number.Decimal((j+1)*s.Plays), i, number.Decimal(startingWealth[i]), number.Decimal(s.PopWealth[i]), number.Decimal(s.TotalWealth))
				os.Exit(0)
			}
		}
	}
}
