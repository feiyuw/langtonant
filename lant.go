package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	defaultInterval = 100
	defaultMove     = 10000

	STEP = 5 // increase 5 cells once
	UP   = iota
	RIGHT
	DOWN
	LEFT
)

var (
	sb strings.Builder // string builder for drawing

	clearScreen = "\u001b[2J"
	cellAlive   = []byte("\u001b[48;5;28m  \u001b[0m")  // green
	cellDead    = []byte("\u001b[48;5;252m  \u001b[0m") // white

	antUp    = []byte("\u001b[48;5;252m\u001b[38;5;196m\u25b2=\u001b[0m")
	antDown  = []byte("\u001b[48;5;252m\u001b[38;5;196m\u25bc=\u001b[0m")
	antLeft  = []byte("\u001b[48;5;252m\u001b[38;5;196m\u25c0=\u001b[0m")
	antRight = []byte("\u001b[48;5;252m\u001b[38;5;196m=\u25b6\u001b[0m")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newGame() *ant {
	game := &ant{
		B: board{
			Rows:  STEP,
			Cols:  STEP,
			Cells: make([][]bool, STEP),
		},
		X:       STEP / 2,
		Y:       STEP / 2,
		D:       UP,
		StepCnt: 0,
	}

	for x := 0; x < STEP; x++ {
		game.B.Cells[x] = make([]bool, STEP)
	}

	return game
}

type board struct {
	Rows  int
	Cols  int
	Cells [][]bool
}

type ant struct {
	B board
	X int // ant's x position [0, B.Rows)
	Y int // ant's y position [0, B.Cols)
	D int // direction

	StepCnt int
}

func (a *ant) draw() {
	sb.Reset()
	sb.WriteString(fmt.Sprintf("\u001b[%dD", (a.B.Rows+1)*(a.B.Cols+1)))
	sb.WriteString(fmt.Sprintf("\u001b[%dA", a.B.Rows+1))
	sb.WriteString(fmt.Sprintf("\u001b[38;5;85mSteps: %d\u001b[0m\n", a.StepCnt))

	for x := 0; x < a.B.Rows; x++ {
		for y := 0; y < a.B.Cols; y++ {
			if x == a.X && y == a.Y { // draw ant
				switch a.D {
				case UP:
					sb.Write(antUp)
				case DOWN:
					sb.Write(antDown)
				case LEFT:
					sb.Write(antLeft)
				case RIGHT:
					sb.Write(antRight)
				default:
					log.Fatalf("invalid ant direction: %d", a.D)
				}
			} else { // draw cells
				if a.B.Cells[x][y] { // alive
					sb.Write(cellAlive)
				} else { // dead
					sb.Write(cellDead)
				}
			}
		}
		sb.WriteByte('\n')
	}
	fmt.Print(sb.String())
}

// left 90°
func (a *ant) turnLeft() {
	switch a.D {
	case UP:
		a.D = LEFT
	case RIGHT:
		a.D = UP
	case DOWN:
		a.D = RIGHT
	case LEFT:
		a.D = DOWN
	default:
		log.Fatalf("invalid direction: %d", a.D)
	}
}

// right 90°
func (a *ant) turnRight() {
	switch a.D {
	case UP:
		a.D = RIGHT
	case RIGHT:
		a.D = DOWN
	case DOWN:
		a.D = LEFT
	case LEFT:
		a.D = UP
	default:
		log.Fatalf("invalid direction: %d", a.D)
	}
}

func (a *ant) oneStep() {
	switch a.D {
	case UP:
		a.X--
	case RIGHT:
		a.Y++
	case DOWN:
		a.X++
	case LEFT:
		a.Y--
	default:
		log.Fatalf("invalid direction: %d", a.D)
	}

	// increase board size if needed
	if a.X < 1 { // increase top margin
		a.B.Rows += STEP
		a.X += STEP
		extendedCells := make([][]bool, STEP, a.B.Rows)
		for x := 0; x < STEP; x++ {
			extendedCells[x] = make([]bool, a.B.Cols)
		}
		a.B.Cells = append(extendedCells, a.B.Cells...)
	}

	if a.X >= (a.B.Rows - 1) { // increase bottom margin
		a.B.Rows += STEP
		for x := 0; x < STEP; x++ {
			a.B.Cells = append(a.B.Cells, make([]bool, a.B.Cols))
		}
	}

	if a.Y <= 1 { // increase left margin
		a.B.Cols += STEP
		a.Y += STEP
		for x := 0; x < a.B.Rows; x++ {
			extendedCells := make([]bool, STEP, a.B.Cols)
			a.B.Cells[x] = append(extendedCells, a.B.Cells[x]...)
		}
	}

	if a.Y >= (a.B.Cols - 1) { // increase right margin
		a.B.Cols += STEP
		for x := 0; x < a.B.Rows; x++ {
			a.B.Cells[x] = append(a.B.Cells[x], make([]bool, STEP)...)
		}
	}
}

// move one step to the direction
func (a *ant) move() {
	if a.B.Cells[a.X][a.Y] { // in alive cell, turn left 90 degree, and change cell to dead
		a.B.Cells[a.X][a.Y] = false
		a.turnLeft()
	} else { // in dead cell, turn right 90 degree, and change cell to alive
		a.B.Cells[a.X][a.Y] = true
		a.turnRight()
	}
	a.oneStep()
	a.StepCnt++
}

func main() {
	var interval, move int

	flag.IntVar(&interval, "i", defaultInterval, "interval between iterations (ms)")
	flag.IntVar(&move, "m", defaultMove, "limit of moves of ant")
	flag.Parse()

	if interval <= 0 {
		log.Printf("Invalid interval %d, use default value", interval)
		interval = defaultInterval
	}
	if move <= 0 {
		log.Printf("Invalid move %d, use default value", move)
		move = defaultMove
	}

	game := newGame()
	sleepInterval := time.Duration(interval) * time.Millisecond
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	fmt.Print(clearScreen)
	for {
		select {
		case <-c:
			log.Println("Bye.")
			return
		default:
			game.draw()
			game.move()
			if game.StepCnt > move {
				log.Println("reach moves, stop it!")
				return
			}
			time.Sleep(sleepInterval)
		}
	}
}
