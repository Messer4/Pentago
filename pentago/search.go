package pentago

import (
	"fmt"
	"math"
	"math/rand"
)

// RandomMove returns a random valid move.
func (b Board) RandomMove() Move {
	moves := b.ValidMoves()
	return moves[rand.Intn(len(moves))]
}

type cellProbs struct {
	white float32
	black float32
}

// Get the "probability" of each cell being black or white in the future.
// This is an approximation of the actual probability, and simply distributes the value of a cell
// to all the positions it could be in in the future.
// All values are in [0, 1].
func (b Board) getProbs() [][]cellProbs {
	probs := make([][]cellProbs, 6)
	for r := range probs {
		probs[r] = make([]cellProbs, 6)
	}

	add := func(r, c int, addWhite, addBlack float32) {
		probs[r][c].white += addWhite
		probs[r][c].black += addBlack
	}

	for r := range b {
		for c := range b[r] {
			var addWhite, addBlack float32
			switch b[r][c] {
			case White:
				addWhite = .25
			case Black:
				addBlack = .25
			case Empty:
				addWhite = .25 * .33
				addBlack = .25 * .33
			}

			// Calculate all positions that this position could rotate to
			subStartR := (r / 3) * 3
			subStartC := (c / 3) * 3
			var m, n, o, p int
			m, n = r-subStartR, c-subStartC
			if m == 1 || n == 1 {
				o, p = 2-m, 2-n
				add(subStartR+m, subStartC+n, addWhite, addBlack)
				add(subStartR+n, subStartC+m, addWhite, addBlack)
				add(subStartR+o, subStartC+p, addWhite, addBlack)
				add(subStartR+p, subStartC+o, addWhite, addBlack)
			} else {
				o, p = 2-m, n
				add(subStartR+m, subStartC+n, addWhite, addBlack)
				add(subStartR+(2-m), subStartC+(2-n), addWhite, addBlack)
				add(subStartR+o, subStartC+p, addWhite, addBlack)
				add(subStartR+(2-o), subStartC+(2-p), addWhite, addBlack)
			}
		}
	}

	return probs
}

type span struct {
	r int
	c int
	deltaR int
	deltaC int
	pieces int
}

func inBounds(r, c int) bool {
	return r > 0 && c > 0 && r < 6 && c < 6;
}

func (b Board) getSpan(r, c, deltaR, deltaC int) span {
	s := span{r, c, deltaR, deltaC, 0}
	color := b[r][c]
	for b[r][c] == color {
		s.pieces++
		r += deltaR
		c += deltaC
		if !inBounds(r, c) || r % 3 == 0 || c % 3 == 0 {
			break
		}
	}
	return s
}

func (b Board) getSpans(r, c int) []span {
	s := []span{
		b.getSpan(r, c, 0, 1),
		b.getSpan(r, c, 1, 0),
	}
	// Add diagonals if this cell is in one that can win
	if -1 <= (r-c) && (r-c) <= +1 {
		s = append(s, b.getSpan(r, c, 1, 1))
	}
	if 4 <= (r+c) && (r+c) <= 6 {
		s = append(s, b.getSpan(r, c, -1, 1))
	}
	return s
}

func (b Board) getSpanWinProb(probs [][]cellProbs, s span) (float32, Piece) {
	if s.pieces == 5 {
		return 1.0, b[s.r][s.c]
	}

	// Assume this span remains stationary (simplification)
	// Go past beginning, then past end
	prob := float32(1.0)
	color := b[s.r][s.c]
	for r, c, l := s.r, s.c, s.pieces; inBounds(r, c) && l <= 5; l++ {
		switch color {
		case White: prob *= probs[r][c].white
		case Black: prob *= probs[r][c].black
		}
		r -= s.deltaR; c -= s.deltaC
	}
	for r, c, l := s.r + s.deltaR*s.pieces, s.c + s.deltaC*s.pieces, s.pieces; inBounds(r, c) && l <= 5; l++ {
		switch color {
		case White: prob *= probs[r][c].white
		case Black: prob *= probs[r][c].black
		}
		r += s.deltaR; c += s.deltaC
	}
	return prob, Empty
}

const whiteWin = -math.MaxFloat32/2
const blackWin = +math.MaxFloat32/2

// Evaluate returns a score for the board. The more positive the score is, the better a position it
// is for black and worse for white, and vice versa.
func (b Board) Evaluate() float32 {
	var score float32

	probs := b.getProbs()
	for r := range b {
		for c := range b[r] {
			if b[r][c] != Empty {
				spans := b.getSpans(r, c)
				for i := range spans {
					prob, winner := b.getSpanWinProb(probs, spans[i])
					switch winner {
						case White: return -whiteWin
						case Black: return +blackWin
					}
					switch b[r][c] {
						case White: score -= prob
						case Black: score += prob
					}
				}
			}
		}
	}
	return score
}

// BestMove calculates and returns the best available more for the given color.
func (b Board) BestMove(color Piece) Move {
	const depth = 3
	m, s := b.getBestMove(0, depth, color, -math.MaxFloat32, +math.MaxFloat32)
	fmt.Printf("best score: %f (depth %d)\n", s, depth)
	return m
}

func (b Board) getBestMove(depth, maxDepth int, color Piece, alpha, beta float32) (Move, float32) {
	if depth == maxDepth {
		return Move{}, b.Evaluate()
	}
	switch winner := b.CheckWinner(); winner {
		case White: return Move{}, whiteWin
		case Black: return Move{}, blackWin
	}

	var bestMove Move
	moves := b.ValidMoves()
	switch color {
	case Black:
		// maxmimizing
		for _, move := range moves {
			bp := b.Clone()
			bp.ApplyMove(move, Black)

			_, score := bp.getBestMove(depth+1, maxDepth, White, alpha, beta)
			if score > alpha {
				alpha = score
				bestMove = move
			}
			if beta <= alpha {
				break
			}
		}
		return bestMove, alpha
	case White:
		// minimizing
		for _, move := range moves {
			bp := b.Clone()
			bp.ApplyMove(move, White)

			_, score := bp.getBestMove(depth+1, maxDepth, Black, alpha, beta)
			if score < beta {
				beta = score
				bestMove = move
			}
			if beta <= alpha {
				break
			}
		}
		return bestMove, beta
	}
	return Move{}, 0
}
