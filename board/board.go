package board

import (
	"math/rand"
	"time"
)

// Board size
const (
	X = 4
	Y = 4
)

// Direction
type Direction int32

// Direction
const (
	LEFT = iota
	UP
	RIGHT
	DOWN
)

var (
	DIRECTIONS = map[int]int{
		LEFT:  -1,
		UP:    -1,
		RIGHT: 1,
		DOWN:  1,
	}
)

type Board struct {
	Cells [Y][X]int

	goal   int
	points int
}

func New() Board {
	board := Board{
		/*
		   Cells: [Y][X]int {
		       {0, 3, 0, 0},
		       {1, 0, 2, 0},
		       {2, 1, 1, 0},
		       {0, 6, 5, 0},
		   },
		*/
		Cells: [Y][X]int{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		goal:   2048,
		points: 0,
	}

	// Seed rng
	rand.Seed(time.Now().Unix())

	// Add two random tiles
	board.AddTile()
	board.AddTile()

	return board
}

func (b *Board) emptyRow(n int) []int {
	row := []int{0, 0, 0, 0}
	return row[:n]
}

func (b *Board) moveLine(row [4]int, direction int) [4]int {
	var empty []int
	var nonEmpty []int
	var result [4]int

	for i := 0; i < len(row); i++ {
		if row[i] == 0 {
			continue
		}

		nonEmpty = append(nonEmpty, row[i])
	}

	empty = b.emptyRow(X - len(nonEmpty))

	// Copy merges to result array
	if direction == -1 {
		copy(result[:], append(nonEmpty, empty...)[0:4])
	} else {
		copy(result[:], append(empty, nonEmpty...)[0:4])
	}

	return result
}

// Is a given line mergeable
func canMergeLine(row [4]int) bool {
	for i := 0; i < len(row); i++ {
		// Previous
		if i > 0 && row[i] == row[i-1] {
			return true
		}

		// Next
		if i+1 < len(row) && row[i] == row[i+1] {
			return true
		}
	}

	return false
}

func (b *Board) mergeLine(row [4]int, direction int) [4]int {
	var newRow [4]int
	var start, end, pos, nextpos int

	if direction == -1 {
		end = 0
		start = len(row) - 1
	} else {
		start = 0
		end = len(row) - 1
	}

	pos = start
	for i := 0; i < len(row); i++ {
		nextpos = pos + direction

		// Don't merge empty cells
		// or already merged cells
		if row[pos] == 0 || newRow[pos] != 0 {
			pos = nextpos
			continue
		}

		// Next cell is identical
		if pos != end && row[pos] == row[nextpos] {
			var value = row[pos] + 1
			newRow[pos] = 0
			newRow[nextpos] = value
		} else {
			newRow[pos] = row[pos]
		}

		// Update position
		pos = nextpos
	}

	return newRow
}

func (b *Board) setRow(y int, row [4]int) {
	for x := 0; x < X; x++ {
		b.Cells[y][x] = row[x]
	}
}

func (b *Board) getRow(y int) [4]int {
	return b.Cells[y]
}

func (b *Board) setCol(x int, row [4]int) {
	for y := 0; y < Y; y++ {
		b.Cells[y][x] = row[y]
	}
}

func (b *Board) getCol(x int) [4]int {
	var a [4]int

	for y := 0; y < Y; y++ {
		a[y] = b.Cells[y][x]
	}

	return a
}

func (b *Board) moveRows(d int) {
	for y := 0; y < Y; y++ {
		// Get new row by moving and merging previous row
		var newRow = b.moveLine(
			b.mergeLine(
				b.moveLine(
					b.getRow(y),
					d,
				),
				d,
			),
			d,
		)

		// Set new row
		b.setRow(y, newRow)
	}
}

func (b *Board) moveCols(d int) {
	for x := 0; x < X; x++ {
		// Get new col by moving and merging previous col
		var newCol = b.moveLine(
			b.mergeLine(
				b.moveLine(
					b.getCol(x),
					d,
				),
				d,
			),
			d,
		)

		// Set new col
		b.setCol(x, newCol)
	}
}

func (b *Board) emptyCells() [][2]int {
	var arr [][2]int

	for y := 0; y < Y; y++ {
		for x := 0; x < X; x++ {
			if b.Cells[y][x] == 0 {
				var cell = [2]int{x, y}
				arr = append(arr, cell)
			}
		}
	}

	return arr
}

func (b *Board) IsFull() bool {
	return len(b.emptyCells()) == 0
}

func (b *Board) AddTile() {
	cells := b.emptyCells()
	cell := cells[rand.Int()%len(cells)]
	x := cell[0]
	y := cell[1]

	// Set cell randomly to 1 or 2
	b.Cells[y][x] = (rand.Int() % 2) + 1
}

func (b *Board) Playable() bool {
	if !b.IsFull() {
		return true
	}

	for y := 0; y < Y; y++ {
		if canMergeLine(b.getRow(y)) {
			return true
		}
	}

	for x := 0; x < X; x++ {
		if canMergeLine(b.getCol(x)) {
			return true
		}
	}

	return false
}

func (b *Board) Values() []int {
	var arr []int

	for y := 0; y < Y; y++ {
		for x := 0; x < X; x++ {
			if b.Cells[y][x] != 0 {
				arr = append(arr, b.Cells[y][x])
			}
		}
	}
	return arr
}

// Move board in a given direction
func (b *Board) Move(d Direction) {
	switch d {
	case UP:
		b.moveCols(DIRECTIONS[UP])
	case DOWN:
		b.moveCols(DIRECTIONS[DOWN])

	case LEFT:
		b.moveRows(DIRECTIONS[LEFT])
	case RIGHT:
		b.moveRows(DIRECTIONS[RIGHT])
	}

	// Add new tile if not empty
	if !b.IsFull() {
		b.AddTile()
	}
}
