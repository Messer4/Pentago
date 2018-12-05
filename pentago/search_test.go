package pentago

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

func newProbs() [][]cellProbs {
	probs := make([][]cellProbs, 6)
	for r := range probs {
		probs[r] = make([]cellProbs, 6)
		for c := range probs[r] {
			probs[r][c] = cellProbs{.33, .33}
		}
	}
	return probs
}

func probsEqual(a, b [][]cellProbs) bool {
	for r := range a {
		for c := range a[r] {
			if math.Abs(float64(a[r][c].black - b[r][c].black)) > .01 ||
				math.Abs(float64(a[r][c].white - b[r][c].white)) > .01 {
				return false
			}
		}
	}
	return true
}

func probsString(p [][]cellProbs) string {
	var buffer bytes.Buffer
	for r := range p {
		for c := range p[r] {
			buffer.WriteString(fmt.Sprintf("{%.2f %.2f} ", p[r][c].black, p[r][c].white))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func TestGetProbs(t *testing.T) {
	b := NewBoard()
	want := newProbs()
	got := b.getProbs()
	if !probsEqual(want, got) {
		t.Errorf("want all .33's, got\n%v", probsString(got))
	}

	b = NewBoard()
	b[0][1] = Black
	got = b.getProbs()
	want = newProbs()
	black := float32(.25 + .33*.25*3);
	white := float32(.33*.25*3);
	want[0][1].black, want[0][1].white = black, white
	want[1][2].black, want[1][2].white = black, white
	want[2][1].black, want[2][1].white = black, white
	want[1][0].black, want[1][0].white = black, white
	if !probsEqual(want, got) {
		t.Errorf("side: want\n%v, got\n%v", probsString(want), probsString(got))
	}

	b = NewBoard()
	b[3][2] = Black
	got = b.getProbs()
	want = newProbs()
	want[3][0].black, want[3][0].white = black, white
	want[3][2].black, want[3][2].white = black, white
	want[5][0].black, want[5][0].white = black, white
	want[5][2].black, want[5][2].white = black, white
	if !probsEqual(want, got) {
		t.Errorf("corner: want\n%v, got\n%v", probsString(want), probsString(got))
	}

	b = NewBoard()
	b[4][1] = Black
	got = b.getProbs()
	want = newProbs()
	want[4][1].black, want[4][1].white = 1, 0
	if !probsEqual(want, got) {
		t.Errorf("middle: want\n%v, got\n%v", probsString(want), probsString(got))
	}
}
