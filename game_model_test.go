package wortklauberei

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewGameSuccess(t *testing.T) {
	// arrange
	var (
		word           = "hellö"
		wordLength     = 5
		validWords     = []string{"hellö", "funny"}
		gameOver       = false
		curAttempt     = 0
		maxAttempts    = 6
		keyboardColors = KeyboardColors{}
	)

	// act
	m, err := NewGame(word, validWords, maxAttempts)

	// assert
	if err != nil {
		t.Error(err)
	}
	if m.word != word {
		t.Errorf("model initialization failed: word %s != %s", m.word, word)
	}
	if m.wordLength != wordLength {
		t.Errorf("model initialization failed: word length %d != %d", m.wordLength, wordLength)
	}
	if !reflect.DeepEqual(m.validWords, validWords) {
		t.Errorf("model initialization failed: valid words %v != %v", m.validWords, validWords)
	}
	if m.gameOver != gameOver {
		t.Errorf("model initialization failed: game over %t != %t", m.gameOver, gameOver)
	}
	if m.curAttempt != curAttempt {
		t.Errorf("model initialization failed: current attempt %d != %d", m.curAttempt, curAttempt)
	}
	if m.maxAttempts != maxAttempts {
		t.Errorf("model initialization failed: max attempts %d != %d", m.maxAttempts, maxAttempts)
	}
	if !reflect.DeepEqual(m.keyboardColors, keyboardColors) {
		t.Errorf("model initialization failed: keyboard colors %v != %v", m.keyboardColors, keyboardColors)
	}
}

func TestNewGameFail(t *testing.T) {
	tests := []struct {
		testcase    string
		word        string
		validWords  []string
		maxAttempts int
		expected    error
	}{
		{
			"word contains space",
			"he llo",
			[]string{"he llo"},
			6,
			fmt.Errorf("game initialization failed: word %s must consist of letters only", "he llo"),
		},
		{
			"word contains dot",
			"hello.",
			[]string{"hello."},
			6,
			fmt.Errorf("game initialization failed: word %s must consist of letters only", "hello."),
		},
		{
			"word not in word list",
			"hello",
			[]string{"hallo"},
			6,
			fmt.Errorf("game initialization failed: word %s not in list of valid words %v", "hello", []string{"hallo"}),
		},
		{
			"max attempts = 0",
			"hello",
			[]string{"hello"},
			0,
			fmt.Errorf("game initialization failed: number of maximum attempts must be greater than 0 but is %d", 0),
		},
		{
			"max attempts is negative",
			"hello",
			[]string{"hello"},
			-1,
			fmt.Errorf("game initialization failed: number of maximum attempts must be greater than 0 but is %d", -1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			_, err := NewGame(tt.word, tt.validWords, tt.maxAttempts)
			if err.Error() != tt.expected.Error() {
				t.Errorf("error expected: %v != %v", tt.expected, err)
			}
		})
	}
}

// TODO
func TestCheckSolutionCorrect(t *testing.T) {
	// arrange
	word := "hello"
	m, _ := NewGame(word, []string{}, 6)

	// act
	c, o, cc, _, v := m.CheckSolution(word)

	// assert
	if !c && !o && cc != nil && v != nil {
		t.Errorf("m.checkSolution(%s) != %t", word, true)
	}
}

// TODO
func TestCheckSolutionWrongWord(t *testing.T) {
	// arrange
	m, _ := NewGame("hello", []string{}, 6)

	// act
	s := "hallo"
	c, o, cc, _, v := m.CheckSolution(s)

	// assert
	if c && !o && cc != nil && v != nil {
		t.Errorf("m.checkSolution(%s) == %t", s, true)
	}
}

// TODO
func TestCountRunes(t *testing.T) {
	tests := []struct {
		testcase string
		s        string
		expected map[rune]uint8
	}{
		{
			"test casual string: hello",
			"hello",
			map[rune]uint8{
				104: 1,
				101: 1,
				108: 2,
				111: 1,
			},
		},
		{
			"umlauts: ööÖßü",
			"ööÖßü",
			map[rune]uint8{
				246: 2,
				214: 1,
				223: 1,
				252: 1,
			},
		},
		{
			"empty string",
			"",
			map[rune]uint8{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			cnt := cacluateRuneHistogram(tt.s)
			if !reflect.DeepEqual(cnt, tt.expected) {
				t.Errorf("countRunes(%s) = %v, expected %v", tt.s, cnt, tt.expected)
			}
		})
	}
}

// TODO
func TestCalculateColorCode(t *testing.T) {
	tests := []struct {
		testcase   string
		suggestion string
		solution   string
		expected   ColorCode
	}{
		{
			"suggestion == solution",
			"hello",
			"hello",
			ColorCode{Green, Green, Green, Green, Green},
		},
		{
			"len(suggestion) != len(solution)",
			"hello",
			"hell",
			nil,
		},
		{
			"some green, some yellow, some gray",
			"halle",
			"hello",
			ColorCode{Green, Gray, Green, Green, Yellow},
		},
		{
			"one letter green, yellow and gray",
			"laffllk",
			"laflenr",
			ColorCode{Green, Green, Green, Gray, Yellow, Gray, Gray},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			g, _ := NewGame(tt.solution, []string{tt.solution, tt.suggestion}, 6)
			cc := g.calculateColorCodeAndUpdateKeyboardColors(tt.suggestion)
			if !reflect.DeepEqual(cc, tt.expected) {
				t.Errorf("calculateColorCode(%s, %s) = %v, expected %v", tt.suggestion, tt.solution, cc, tt.expected)
			}
			// TODO check g.usedLetters
		})
	}
}

// TODO
func TestIsAllLetters(t *testing.T) {
	tests := []struct {
		testcase string
		s        string
		expected bool
	}{
		{
			"Testcase IsAllLetters(\"hallo\")",
			"hallo",
			true,
		},
		{
			"Testcase IsAllLetters(\"öäüßÖÄÜ\")",
			"öäüßÖÄÜ",
			true,
		},
		{
			"Testcase IsAllLetters(\"\")",
			"",
			true,
		},
		{
			"Testcase IsAllLetters(\"hallo du\")",
			"hallo du",
			false,
		},
		{
			"Testcase IsAllLetters(\"hallo.du\")",
			"hallo.du",
			false,
		},
		{
			"Testcase IsAllLetters(\"hallo7du\")",
			"hallo7du",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			b := IsAllLetters(tt.s)
			if b != tt.expected {
				t.Errorf("IsAllLetters(%s) = %t", tt.s, b)
			}
		})
	}
}

// TODO
func TestWordInWordList(t *testing.T) {
	tests := []struct {
		testcase string
		word     string
		wordList []string
		expected bool
	}{
		{
			"Testcase WordInWordList('hey', ['you', 'hey', 'there'])",
			"hey",
			[]string{"you", "hey", "there"},
			true,
		},
		{
			"Testcase WordInWordList('hey', ['you', 'there'])",
			"hey",
			[]string{"you", "there"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			b := WordInWordList(tt.word, tt.wordList)
			if b != tt.expected {
				t.Errorf("WordInWordList(%s, %v) = %t", tt.word, tt.wordList, b)
			}
		})
	}
}
