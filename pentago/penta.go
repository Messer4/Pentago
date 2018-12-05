// pentago contains an implementation of the Pentago game logic and a Minimax search algorithm.
package pentago

import (
	"bytes"
)

// Piece represents the color of the piece. It can also be Empty, meaning no piece.
type Piece int

const (
	Empty Piece = iota
	Black
	White
)

// Board represents the state of the board in a Pentago game.
type Board [][]Piece

// NewBoard makes a new empty Pentago board.
func NewBoard() Board {
	b := make([][]Piece, 6)
	for r := 0; r < 6; r++ {
		b[r] = make([]Piece, 6)
	}
	return b
}

// Clone makes a copy of an existing board.
func (b Board) Clone() Board {
	b2 := NewBoard()
	for r := range b {
		copy(b2[r], b[r])
	}
	return b2
}

// String returns the board state as human-readable string.
func (b Board) String() string {
	var buffer bytes.Buffer
	for _, row := range b {
		for _, piece := range row {
			switch piece {
			case Empty:
				buffer.WriteString(" .")
			case Black:
				buffer.WriteString(" B")
			case White:
				buffer.WriteString(" W")
			}
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

// Coord represents a location on the board.
type Coord struct {
	Row int  // 0-based row, beginning at the top
	Col int  // 0-based column, beginning at the left
}

// Get piece returns the color of the piece at the given Coord.
func (b Board) GetPiece(c Coord) Piece {
	return b[c.Row][c.Col]
}

func (b Board) getTriple(start Coord, deltaR int, deltaC int) [3]Piece {
	return [3]Piece{
		b[start.Row][start.Col],
		b[start.Row+deltaR][start.Col+deltaC],
		b[start.Row+2*deltaR][start.Col+2*deltaC],
	}
}

func (b Board) setTriple(start Coord, deltaR int, deltaC int, triple [3]Piece) {
	b[start.Row][start.Col] = triple[0]
	b[start.Row+deltaR][start.Col+deltaC] = triple[1]
	b[start.Row+2*deltaR][start.Col+2*deltaC] = triple[2]
}

const (
	Clockwise = iota
	CounterClockwise
)

// Rotate rotates one of the sub-boards, using a 0-based row and column (subRow and subCol are each
// either 0 or 1.) direction is either Clockwise or CounterClockwise.
func (b Board) Rotate(subRow int, subCol int, direction int) {
	topLeft := Coord{3*subRow, 3*subCol}
	botLeft := Coord{topLeft.Row+2, topLeft.Col}
	botRight := Coord{botLeft.Row, botLeft.Col+2}
	topRight := Coord{topLeft.Row, topLeft.Col+2}

	top := b.getTriple(topLeft, 0, 1)
	switch direction {
	case Clockwise:
		b.setTriple(topLeft, 0, 1, b.getTriple(botLeft, -1, 0))
		b.setTriple(botLeft, -1, 0, b.getTriple(botRight, 0, -1))
		b.setTriple(botRight, 0, -1, b.getTriple(topRight, 1, 0))
		b.setTriple(topRight, 1, 0, top)
	case CounterClockwise:
		b.setTriple(topLeft, 0, 1, b.getTriple(topRight, 1, 0))
		b.setTriple(topRight, 1, 0, b.getTriple(botRight, 0, -1))
		b.setTriple(botRight, 0, -1, b.getTriple(botLeft, -1, 0))
		b.setTriple(botLeft, -1, 0, top)
	}
}

// Checks if a span of consecutive pieces of the same color exists on a given
// row or column. If it does, returns the color of the span. If not, returns
// Empty.
func (b Board) checkSpan(coord Coord, deltaR, deltaC, length, slotLength int) Piece {
	for i := 0; i < (slotLength-length); i++ {
		// If first two are the same, they are part of the span in this row if it
		// exists. If not, the only way is if the span starts at the next index.
		if b[coord.Row][coord.Col] != b[coord.Row+deltaR][coord.Col+deltaC] {
			coord.Row += deltaR
			coord.Col += deltaC
		}
	}

	color := b[coord.Row][coord.Col]
	for i := 1; i < length; i++ {
		coord.Row += deltaR
		coord.Col += deltaC
		if b[coord.Row][coord.Col] != color {
			return Empty
		}
	}
	return color
}

// CheckWinner returns the color of the winner if there is one, or Empty if none.
func (b Board) CheckWinner() Piece {
	// Check rows
	for r := 0; r < 6; r++ {
		if color := b.checkSpan(Coord{r, 0}, 0, 1, 5, 6); color != Empty {
			return color
		}
	}
	// Check cols
	for c := 0; c < 6; c++ {
		if color := b.checkSpan(Coord{0, c}, 1, 0, 5, 6); color != Empty {
			return color
		}
	}
	// Check diags
	if color := b.checkSpan(Coord{0, 1}, 1, 1, 5, 5); color != Empty {
		return color
	}
	if color := b.checkSpan(Coord{1, 0}, 1, 1, 5, 5); color != Empty {
		return color
	}
	if color := b.checkSpan(Coord{0, 0}, 1, 1, 5, 6); color != Empty {
		return color
	}
	if color := b.checkSpan(Coord{4, 0}, -1, 1, 5, 5); color != Empty {
		return color
	}
	if color := b.checkSpan(Coord{5, 1}, -1, 1, 5, 5); color != Empty {
		return color
	}
	if color := b.checkSpan(Coord{5, 0}, -1, 1, 5, 6); color != Empty {
		return color
	}
	return Empty
}

// Move describes a move that a player could make.
type Move struct {
	R int `json:"r"`
	C int `json:"c"`
	Sub int `json:"sub"`
	Dir int `json:"dir"`
}

// NewMove creates a move, given the row and column of the piece, the sub-board to be rotated (0-3),
// and the direction of the rotation.
func NewMove(r, c, sub, dir int) Move {
	return Move{r, c, sub, dir}
}

// String returns the move in a human-readable form.
//func (m Move) String() string {
//	var directions = map[int]string{Clockwise: "CW", CounterClockwise: "CCW"}
//	return fmt.Sprintf("Put piece (%d, %d), rotate subboard %d %s", m.r, m.c, m.sub, directions[m.dir])
//}

// IsValid returns whether the move is valid.
func (m Move) IsValid(b Board) bool {
	return b[m.R][m.C] == Empty
}

// ValidMoves returns a list of valid moves for a given board state.
func (b Board) ValidMoves() []Move {
	moves := make([]Move, 0, 50)
	for r := 0; r < 6; r++ {
		for c := 0; c < 6; c++ {
			if b[r][c] == Empty {
				for sub := 0; sub < 4; sub++ {
					moves = append(moves, Move{r, c, sub, 0})
					moves = append(moves, Move{r, c, sub, 1})
				}
			}
		}
	}
	return moves
}

// ApplyMove executes a move on the board state and returns whether it was successful.
func (b Board) ApplyMove(m Move, color Piece) bool {
	if b[m.R][m.C] != Empty {
		return false
	}
	b[m.C][m.C] = color
	b.Rotate(m.Sub / 2, m.Sub % 2, m.Dir)
	return true
}

// Game represents the game state: both a board state and whose turn it is.
type Game struct {
	Board
	Turn Piece
}

// NewGame creates a new Game.
func NewGame() Game {
	return Game{
		Board: NewBoard(),
		Turn: White,
	}
}

// Move executes the given Move on a game state.
func (g *Game) Move(m Move) bool {
	if !g.Board.ApplyMove(m, g.Turn) {
		return false
	}

	switch g.Turn {
	case White: g.Turn = Black
	case Black: g.Turn = White
	}
	return true
}
