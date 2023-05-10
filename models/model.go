package models

import (
	"fmt"

	"github.com/yawn77/wortklauberei/utils"
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

type GameModel struct {
	Word           string
	WordLength     int
	ValidWords     []string
	GameOver       bool
	CurAttempt     int
	MaxAttempts    int
	KeyboardColors KeyboardColors
}

func NewGameModel(word string, validWords []string, maxAttempts int) (m GameModel, err error) {
	// TODO check for valid letters (those which are on keyboard)
	if !utils.IsLettersOnly(word) {
		return m, fmt.Errorf("game initialization failed: word %s must consist of letters only", word)
	}
	if !utils.IsWordInWordList(word, validWords) {
		return m, fmt.Errorf("game initialization failed: word %s not in list of valid words %v", word, validWords)
	}
	if maxAttempts <= 0 {
		return m, fmt.Errorf("game initialization failed: number of maximum attempts must be greater than 0 but is %d", maxAttempts)
	}
	m.Word = word
	m.WordLength = len([]rune(word))
	m.ValidWords = validWords
	m.GameOver = false
	m.CurAttempt = 0
	m.MaxAttempts = maxAttempts
	m.KeyboardColors = make(KeyboardColors)
	return m, nil
}
