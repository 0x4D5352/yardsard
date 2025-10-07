package main

import (
	"bufio"
	"fmt"
	"image/color"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"

	"github.com/vdobler/chart"
	"github.com/vdobler/chart/txtg"
)

type Agent struct {
	Wealth        []int
	CurrentWealth int
}

type SimulationState struct {
	Agents          []*Agent
	PopIndices      []int // an array of the indexes of Population, for Fisher-Yates
	Now             time.Time
	Population      int // the number of agents in the simulation
	Plays           int // the number of rounds within an iteration step
	Gain            int // how much A gives to B, in percentage of B's wealth
	Loss            int // how much A takes from B, in percentage of B's wealth
	TotalWealth     int // how much money each agent in the simulation begins with
	OligarchLimit   int // how much money an agent needs to be considered an oligarch
	TermWidth       int
	TermHeight      int
	CurrentOligarch int
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
		var wg sync.WaitGroup
		for i := range halfpeople {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				i1 := s.PopIndices[i]
				i2 := s.PopIndices[i+halfpeople]
				// fmt.Printf("Pairing off agent %d against agent %d\n", i1, i2)
				win := rand.IntN(2) == 1
				v1 := s.Agents[i1].CurrentWealth
				v2 := s.Agents[i2].CurrentWealth
				j1 := i1
				j2 := i2
				if v2 > v1 {
					j1 = i2
					j2 = i1
				}
				// fmt.Println("Starting income...")
				// fmt.Printf("Agent 1: $%d\n", s.Agents[j1])
				// fmt.Printf("Agent 2: $%d\n", s.Agents[j2])
				if win {
					// fmt.Println("Poorer agent wins!")
					delta := s.Agents[j2].CurrentWealth * s.Gain / 100
					s.Agents[j1].CurrentWealth -= delta
					s.Agents[j2].CurrentWealth += delta
				} else {
					// fmt.Println("Wealthier agent wins!")
					delta := s.Agents[j2].CurrentWealth * s.Loss / 100
					s.Agents[j1].CurrentWealth += delta
					s.Agents[j2].CurrentWealth -= delta
				}
				s.Agents[j1].Wealth = append(s.Agents[j1].Wealth, s.Agents[j1].CurrentWealth)
				s.Agents[j2].Wealth = append(s.Agents[j2].Wealth, s.Agents[j2].CurrentWealth)
				// fmt.Println("Resulting income...")
				// fmt.Printf("Agent 1: $%d\n", s.Agents[j1])
				// fmt.Printf("Agent 2: $%d\n", s.Agents[j2])
				// fmt.Println()
			}(i)
		}
		// fmt.Println()
		wg.Wait()
		// fmt.Printf("Loop %d of %d done...\n", j+1, s.Plays)
	}
}

func (s *SimulationState) printWealth() string {
	width := s.TermWidth
	height := (s.TermHeight / 3) + (s.TermHeight / 6) - 4
	floatWealth := make([]float64, len(s.Agents))
	xAxis := make([]float64, len(s.Agents))
	maxWealth := 0
	for i, v := range s.Agents {
		if v.CurrentWealth > maxWealth {
			maxWealth = v.CurrentWealth
		}
		xAxis[i] = float64(i + 1)
		floatWealth[i] = float64(v.CurrentWealth)
	}
	var yMax int
	if maxWealth < s.TotalWealth/8 {
		yMax = s.TotalWealth / 8
	} else if maxWealth < s.TotalWealth/4 {
		yMax = s.TotalWealth / 4
	} else if maxWealth < s.TotalWealth/2 {
		yMax = s.TotalWealth / 2
	} else {
		yMax = s.TotalWealth
	}
	wealth := chart.Style{
		Symbol:    '$',
		LineColor: color.NRGBA{0xcc, 0x00, 0x00, 0xff},
		FillColor: color.NRGBA{0xff, 0x80, 0x80, 0xff},
		LineStyle: chart.SolidLine,
		LineWidth: 1,
	}
	barc := chart.BarChart{
		Title:        "Wealth Distribution",
		SameBarWidth: true,
		Stacked:      true,
	}
	barc.Reset()
	barc.Key.Hide = true
	barc.XRange.ShowZero = true
	barc.XRange.MaxMode = chart.RangeMode{Fixed: true, Value: float64(s.Population)}
	barc.XRange.MinMode = chart.RangeMode{Fixed: true, Value: 0}
	barc.YRange.ShowZero = true
	barc.YRange.MaxMode = chart.RangeMode{Fixed: true, Value: float64(yMax)}
	barc.YRange.MinMode = chart.RangeMode{Fixed: true, Value: 0}
	barc.AddDataPair("Wealth", xAxis, floatWealth, wealth)
	barc.XRange.TicSetting.Delta = 0
	barc.YRange.TicSetting.Delta = 0
	tgr := txtg.New(width, height)
	barc.Plot(tgr)
	return tgr.String()
}

func (s *SimulationState) printWealthHistory() string {

	w := s.TermWidth
	h := (s.TermHeight / 3) + (s.TermHeight / 6) - 4
	m := 0
	mi := -1
	for i, v := range s.Agents {
		if v.CurrentWealth > m {
			m = v.CurrentWealth
			mi = i
		}
	}
	if mi == -1 {
		fmt.Printf("uh oh we don't have an oligarch (which is otherwise a good thing)")
	}
	scatter := &chart.ScatterChart{
		Title: fmt.Sprintf("Oligarch: Agent %d", mi),
	}
	scatter.Key.Hide = true
	scatter.Reset()

	timePoints := make([]float64, len(s.Agents[mi].Wealth))
	wealthPoints := make([]float64, len(s.Agents[mi].Wealth))

	for i, wealth := range s.Agents[mi].Wealth {
		timePoints[i] = float64(i)
		wealthPoints[i] = float64(wealth)
	}

	scatter.AddDataPair("Wealth", timePoints, wealthPoints,
		chart.PlotStyleLines, chart.Style{
			Symbol:    '+',
			LineColor: color.NRGBA{0xcc, 0x00, 0x00, 0xff},
			LineStyle: chart.SolidLine,
			LineWidth: 1,
		})

	scatter.XRange.ShowZero = true
	scatter.YRange.ShowZero = true
	scatter.XRange.MinMode = chart.RangeMode{Fixed: true, Value: 0}
	scatter.YRange.MinMode = chart.RangeMode{Fixed: true, Value: 0}
	scatter.XRange.TicSetting.Delta = 0
	scatter.YRange.TicSetting.Delta = 0
	tgr := txtg.New(w, h)
	scatter.Plot(tgr)
	return tgr.String()
}

func initializeSimulation(people, plays, gain, loss, init int, equalWealth bool) *SimulationState {
	fd := int(os.Stdin.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		fmt.Println("Error getting terminal size: ", err)
		width = 80
		height = 24
	}
	s := SimulationState{
		Agents:     make([]*Agent, 0),
		PopIndices: make([]int, 0),
		Population: people,
		Plays:      plays,
		Gain:       gain,
		Loss:       loss,
		TermWidth:  width,
		TermHeight: height,
	}
	for i := range people {
		s.PopIndices = append(s.PopIndices, i)
		var w int
		if equalWealth {
			w = init
		} else {
			w = rand.IntN(init) + 1
		}
		a := Agent{CurrentWealth: w, Wealth: make([]int, 0)}
		s.Agents = append(s.Agents, &a)
		s.TotalWealth += w
	}
	for s.TotalWealth < init*people {
		s.Agents[rand.IntN(people)].CurrentWealth += 1
		s.TotalWealth++
	}
	for i := range people {
		s.Agents[i].Wealth = append(s.Agents[i].Wealth, s.Agents[i].CurrentWealth)
	}
	s.OligarchLimit = int(float64(s.TotalWealth) * 0.95)
	return &s
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to yardsard, a Yard Sale Model Simulator!")
	fmt.Printf("Please enter the number of people to simulate (default = 100): ")
	people := 100
	var err error
	if scanner.Scan() {
		input := scanner.Text()
		if input != "" {
			people, err = strconv.Atoi(input)
			if err != nil {
				// this is all very ugly, forgive me
				people = 100
			}
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
	fmt.Printf("Please enter the amount of money each person should receive (default = $100): ")
	init := 100
	if scanner.Scan() {
		input := scanner.Text()
		if input != "" {
			init, err = strconv.Atoi(input)
			if err != nil {
				// this is all very ugly, forgive me
				people = 100
			}
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
	fps := (1000 / 60) * time.Millisecond
	if people > 3_000 {
		fps = 0 * time.Millisecond
	}
	equalWealth := true
	fmt.Printf("Should wealth be distributed equally for a level playing field? [Y]es or [N]o? (default = y): ")
	if scanner.Scan() {
		input := scanner.Text()
		if input != "" && strings.ToLower(strings.Split(input, "")[0]) == "n" {
			equalWealth = false
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
	// equalWealth := false
	p := message.NewPrinter(language.AmericanEnglish)
	s := initializeSimulation(people, 100, 20, 17, init, equalWealth)
	fmt.Println("Intitial Conditions:")
	p.Printf("Population Size: %d\n", s.Population)
	if equalWealth {
		p.Printf("Max Starting Wealth: $%d\n", number.Decimal(init))
	} else {

		p.Printf("Possible Max Starting Wealth: $%d\n", number.Decimal(init))
	}
	p.Printf("Total Available Wealth: $%d\n", s.TotalWealth)
	p.Printf("Plays per round: %d\n", s.Plays)
	fmt.Printf("TermInfo: Width = %d, Height = %d\n", s.TermWidth, s.TermHeight)
	fmt.Printf("Press enter to start simulaton!")
	startingWealth := make([]int, 0)
	for _, v := range s.Agents {
		startingWealth = append(startingWealth, v.CurrentWealth)
	}
	_, err = fmt.Scanln()
	if err != nil {
		fmt.Println(err)
	}
	s.Now = time.Now()
	j := 0
	for {
		var sb strings.Builder
		s.yardSaleIteration()
		sb.WriteString("\033[?25l")
		sb.WriteString("\033[2J\033[H")
		sb.WriteString(p.Sprintf("Population Size: %d | ", s.Population))
		if equalWealth {
			sb.WriteString(p.Sprintf("Max Starting Wealth: $%d | ", number.Decimal(init)))
		} else {
			sb.WriteString(p.Sprintf("Possible Max Starting Wealth: $%d | ", number.Decimal(init)))
		}
		sb.WriteString(p.Sprintf("Total Available Wealth: $%d | ", s.TotalWealth))
		sb.WriteString(p.Sprintf("Plays per round: %d | ", s.Plays))
		sb.WriteString(p.Sprintf("Round %d | ", number.Decimal(j)))
		now := time.Now()
		dur := now.Sub(s.Now)
		s.Now = now
		sb.WriteString(fmt.Sprintf("Elapsed time between iterations: %d milliseconds", dur.Milliseconds()))
		sb.WriteString("\n\n")
		sb.WriteString(s.printWealth())
		sb.WriteString("\n")
		sb.WriteString(s.printWealthHistory())
		sb.WriteString("\033[?25h")
		fmt.Println(sb.String())
		time.Sleep(fps)
		j++
		for i := range s.Agents {
			if s.Agents[i].CurrentWealth >= s.OligarchLimit {
				fmt.Printf("\033[2J\033[H")
				p.Printf("After %d rounds, Agent %d (started with $%d) has become an oligarch with $%d out of the available $%d!\n\n",
					number.Decimal((j+1)*s.Plays),
					i,
					number.Decimal(startingWealth[i]),
					number.Decimal(s.Agents[i].CurrentWealth),
					number.Decimal(s.TotalWealth))
				fmt.Println(s.printWealth())
				fmt.Println()
				fmt.Println(s.printWealthHistory())
				fmt.Printf("\033[?25h")
				os.Exit(0)
			}
		}
	}
}
