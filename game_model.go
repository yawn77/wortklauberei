package wortklauberei

import (
	"fmt"
	"strings"
	"unicode"
)

type Color uint8
type ColorCode []Color
type KeyboardColors map[rune]Color
type RuneHistogram map[rune]uint8

const (
	Gray   Color = 0
	Yellow Color = 1
	Green  Color = 2
)

type gameModel struct {
	word           string
	wordLength     int
	validWords     []string
	gameOver       bool
	curAttempt     int
	maxAttempts    int
	keyboardColors KeyboardColors
}

func NewGame(word string, validWords []string, maxAttempts int) (m gameModel, err error) {
	// TODO check for valid letters (those which are on keyboard)
	if !IsAllLetters(word) {
		return m, fmt.Errorf("game initialization failed: word %s must consist of letters only", word)
	}
	if !WordInWordList(word, validWords) {
		return m, fmt.Errorf("game initialization failed: word %s not in list of valid words %v", word, validWords)
	}
	if maxAttempts <= 0 {
		return m, fmt.Errorf("game initialization failed: number of maximum attempts must be greater than 0 but is %d", maxAttempts)
	}
	m.word = word
	m.wordLength = len([]rune(word))
	m.validWords = validWords
	m.gameOver = false
	m.curAttempt = 0
	m.maxAttempts = maxAttempts
	m.keyboardColors = make(KeyboardColors)
	return m, nil
}

func (g *gameModel) CheckSolution(s string) (correct bool, gameOver bool, colorCode ColorCode, keyboardColors KeyboardColors, valid error) {
	// game is already over
	if g.gameOver {
		return false, true, nil, nil, fmt.Errorf("game is already over")
	}
	// input validation
	// check word length
	wl := len([]rune(s))
	if wl != g.wordLength {
		return false, false, nil, nil, fmt.Errorf("length of solution must be %d but is %d", g.wordLength, wl)
	}
	// all letters
	if !IsAllLetters(s) {
		return false, false, nil, nil, fmt.Errorf("solution must consists of letters  only")
	}
	// word in wordlist
	if !WordInWordList(s, g.validWords) {
		return false, false, nil, nil, fmt.Errorf("suggested solution is not in list of valid words")
	}

	// correct word or was last attempt
	correct = s == g.word
	if correct || g.curAttempt >= g.maxAttempts-1 {
		g.gameOver = true
	}

	colorCode = g.calculateColorCodeAndUpdateKeyboardColors(s)
	return correct, g.gameOver, colorCode, g.keyboardColors, nil
}

func (g *gameModel) calculateColorCodeAndUpdateKeyboardColors(suggestion string) (colorCode ColorCode) {
	rsug := []rune(suggestion)
	lsug := len(rsug)
	rsol := []rune(g.word)
	if lsug != len(rsol) {
		return nil
	}
	rh := cacluateRuneHistogram(g.word)

	// mark greens
	for i := 0; i < lsug; i++ {
		r := rsug[i]
		if r == rsol[i] {
			colorCode = append(colorCode, Green)
			g.updateKeyboardColor(r, Green)
			val, ok := rh[r]
			if ok {
				rh[r] = val - 1
			}
		} else {
			// may be overridden by yellow in next loop
			colorCode = append(colorCode, Gray)
			g.updateKeyboardColor(r, Gray)
		}
	}

	// mark yellows and grays
	for i := 0; i < lsug; i++ {
		if colorCode[i] == Green {
			continue
		}
		r := rsug[i]
		val, ok := rh[r]
		sInSol := strings.Contains(g.word, string(r))
		if sInSol && ok && val > 0 {
			rh[r] = val - 1
			colorCode[i] = Yellow
			g.updateKeyboardColor(r, Yellow)
		}
	}

	return colorCode
}

func (g *gameModel) updateKeyboardColor(r rune, color Color) {
	val, ok := g.keyboardColors[r]
	if !ok || val < color {
		g.keyboardColors[r] = color
	}
}

func cacluateRuneHistogram(s string) (rh RuneHistogram) {
	rh = make(RuneHistogram)
	for _, r := range s {
		val, ok := rh[r]
		if ok {
			rh[r] = val + 1
		} else {
			rh[r] = 1
		}
	}
	return rh
}

func IsAllLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func WordInWordList(w string, wl []string) bool {
	for _, cw := range wl {
		if cw == w {
			return true
		}
	}
	return false
}
