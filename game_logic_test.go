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
			"test word with space",
			"he llo",
			[]string{"he llo"},
			6,
			fmt.Errorf("game initialization failed: word %s must consist of letters only", "he llo"),
		},
		{
			"test word with dot",
			"hello.",
			[]string{"hello."},
			6,
			fmt.Errorf("game initialization failed: word %s must consist of letters only", "hello."),
		},
		{
			"test word that is not in word list",
			"hello",
			[]string{"hallo"},
			6,
			fmt.Errorf("game initialization failed: word %s not in list of valid words %v", "hello", []string{"hallo"}),
		},
		{
			"test max attempts = 0",
			"hello",
			[]string{"hello"},
			0,
			fmt.Errorf("game initialization failed: number of maximum attempts must be greater than 0 but is %d", 0),
		},
		{
			"test max attempts is negative",
			"hello",
			[]string{"hello"},
			-1,
			fmt.Errorf("game initialization failed: number of maximum attempts must be greater than 0 but is %d", -1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			_, err := NewGame(tt.word, tt.validWords, tt.maxAttempts)
			if err == nil || err.Error() != tt.expected.Error() {
				t.Errorf("error expected: %v != %v", tt.expected, err)
			}
		})
	}
}

func TestCheckSolutionCorrectAfterThreeAttempts(t *testing.T) {
	// arrange
	word := "dizzy"
	m, _ := NewGame(word, []string{"hello", "whizz", "pozzy", "dizzy"}, 6)

	// act & assert
	correct, gameOver, colorCode, keyboardColors, valid := m.CheckSolution("hello")
	if correct ||
		gameOver ||
		!reflect.DeepEqual(colorCode, ColorCode{Gray, Gray, Gray, Gray, Gray}) ||
		!reflect.DeepEqual(keyboardColors, KeyboardColors{104: Gray, 101: Gray, 108: Gray, 111: Gray}) ||
		valid != nil {
		t.Errorf("error after hello")
	}
	correct, gameOver, colorCode, keyboardColors, valid = m.CheckSolution("whizz")
	if correct ||
		gameOver ||
		!reflect.DeepEqual(colorCode, ColorCode{Gray, Gray, Yellow, Green, Yellow}) ||
		!reflect.DeepEqual(keyboardColors, KeyboardColors{104: Gray, 101: Gray, 108: Gray, 111: Gray, 119: Gray, 105: Yellow, 122: Green}) ||
		valid != nil {
		t.Errorf("error after whizz")
	}
	correct, gameOver, colorCode, keyboardColors, valid = m.CheckSolution("dizzy")
	if !correct ||
		!gameOver ||
		!reflect.DeepEqual(colorCode, ColorCode{Green, Green, Green, Green, Green}) ||
		!reflect.DeepEqual(keyboardColors, KeyboardColors{104: Gray, 101: Gray, 108: Gray, 111: Gray, 119: Gray, 105: Green, 122: Green, 100: Green, 121: Green}) ||
		valid != nil {
		t.Errorf("error after dizzy")
	}
}

func TestCheckSolutionCorrectAtLastAttempt(t *testing.T) {
	// arrange
	word := "dizzy"
	m, _ := NewGame(word, []string{"hello", "whizz", "pozzy", "dizzy"}, 3)

	// act & assert
	correct, gameOver, colorCode, keyboardColors, valid := m.CheckSolution("hello")
	if correct ||
		gameOver ||
		!reflect.DeepEqual(colorCode, ColorCode{Gray, Gray, Gray, Gray, Gray}) ||
		!reflect.DeepEqual(keyboardColors, KeyboardColors{104: Gray, 101: Gray, 108: Gray, 111: Gray}) ||
		valid != nil {
		t.Errorf("error after hello")
	}
	correct, gameOver, colorCode, keyboardColors, valid = m.CheckSolution("whizz")
	if correct ||
		gameOver ||
		!reflect.DeepEqual(colorCode, ColorCode{Gray, Gray, Yellow, Green, Yellow}) ||
		!reflect.DeepEqual(keyboardColors, KeyboardColors{104: Gray, 101: Gray, 108: Gray, 111: Gray, 119: Gray, 105: Yellow, 122: Green}) ||
		valid != nil {
		t.Errorf("error after whizz")
	}
	correct, gameOver, colorCode, keyboardColors, valid = m.CheckSolution("dizzy")
	if !correct ||
		!gameOver ||
		!reflect.DeepEqual(colorCode, ColorCode{Green, Green, Green, Green, Green}) ||
		!reflect.DeepEqual(keyboardColors, KeyboardColors{104: Gray, 101: Gray, 108: Gray, 111: Gray, 119: Gray, 105: Green, 122: Green, 100: Green, 121: Green}) ||
		valid != nil {
		t.Errorf("error after dizzy")
	}
}

func TestCheckSolutionAfterGameOver(t *testing.T) {
	// arrange
	word := "dizzy"
	m, _ := NewGame(word, []string{"hello", "whizz", "pozzy", "dizzy"}, 3)

	// act & assert
	correct, gameOver, colorCode, keyboardColors, valid := m.CheckSolution("dizzy")
	if !correct ||
		!gameOver ||
		!reflect.DeepEqual(colorCode, ColorCode{Green, Green, Green, Green, Green}) ||
		!reflect.DeepEqual(keyboardColors, KeyboardColors{100: Green, 105: Green, 122: Green, 121: Green}) ||
		valid != nil {
		t.Errorf("error after dizzy")
	}
	correct, gameOver, colorCode, keyboardColors, valid = m.CheckSolution("whizz")
	if correct ||
		!gameOver ||
		colorCode != nil ||
		keyboardColors != nil ||
		valid.Error() != fmt.Errorf("game is already over").Error() {
		t.Errorf("game was not over yet")
	}
}

func TestCheckSolutionInputValidation(t *testing.T) {
	tests := []struct {
		testcase   string
		word       string
		validWords []string
		solution   string
		expected   error
	}{
		{
			"test invalid word length",
			"hello",
			[]string{"hello", "hell"},
			"hell",
			fmt.Errorf("length of solution must be %d but is %d", 5, 4),
		},
		{
			"test solution contains special characters",
			"hello",
			[]string{"hello", "hell."},
			"hell.",
			fmt.Errorf("solution must consists of letters  only"),
		},
		{
			"test solution contains numbers",
			"hello",
			[]string{"hello", "hell6"},
			"hell6",
			fmt.Errorf("solution must consists of letters  only"),
		},
		{
			"test solution not in word list",
			"hello",
			[]string{"hello"},
			"hallo",
			fmt.Errorf("suggested solution is not in list of valid words"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			// arrange
			g, _ := NewGame(tt.word, tt.validWords, 6)

			// act
			correct, gameOver, colorCode, keyboardColors, valid := g.CheckSolution(tt.solution)

			// assert
			if valid == nil || valid.Error() != tt.expected.Error() {
				t.Errorf("error expected: %v != %v", tt.expected, valid)
			}
			if correct || gameOver || colorCode != nil || keyboardColors != nil {
				t.Error("expected other return values")
			}
		})
	}
}

func TestCalculateColorCodeAndUpdateKeyboardColors(t *testing.T) {
	tests := []struct {
		testcase   string
		suggestion string
		solution   string
		expectedCC ColorCode
		expectedKC KeyboardColors
	}{
		{
			"test suggestion == solution",
			"hello",
			"hello",
			ColorCode{Green, Green, Green, Green, Green},
			KeyboardColors{104: Green, 101: Green, 108: Green, 111: Green},
		},
		{
			"test len(suggestion) != len(solution)",
			"hello",
			"hell",
			nil,
			KeyboardColors{},
		},
		{
			"test some green, some yellow, some gray",
			"halle",
			"hello",
			ColorCode{Green, Gray, Green, Green, Yellow},
			KeyboardColors{104: Green, 97: Gray, 108: Green, 101: Yellow},
		},
		{
			"test one letter green, yellow and gray",
			"laffllk",
			"laflenr",
			ColorCode{Green, Green, Green, Gray, Yellow, Gray, Gray},
			KeyboardColors{108: Green, 97: Green, 102: Green, 107: Gray},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			g, _ := NewGame(tt.solution, []string{tt.solution, tt.suggestion}, 6)
			cc := g.calculateColorCodeAndUpdateKeyboardColors(tt.suggestion)
			if !reflect.DeepEqual(cc, tt.expectedCC) {
				t.Errorf("calculateColorCodeAndUpdateKeyboardColors(%s, %s) = %v, expected %v", tt.suggestion, tt.solution, cc, tt.expectedCC)
			}
			if !reflect.DeepEqual(g.keyboardColors, tt.expectedKC) {
				t.Errorf("unexpected keyboard colors: %v != %v", g.keyboardColors, tt.expectedKC)
			}
		})
	}
}

func TestCalculateRuneHistogram(t *testing.T) {
	tests := []struct {
		testcase string
		s        string
		expected RuneHistogram
	}{
		{
			"test casual string",
			"hello",
			RuneHistogram{
				104: 1,
				101: 1,
				108: 2,
				111: 1,
			},
		},
		{
			"test umlauts",
			"ööÖßü",
			RuneHistogram{
				246: 2,
				214: 1,
				223: 1,
				252: 1,
			},
		},
		{
			"test empty string",
			"",
			RuneHistogram{},
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

func TestIsAllLetters(t *testing.T) {
	tests := []struct {
		testcase string
		s        string
		expected bool
	}{
		{
			"test word with letters only",
			"hallo",
			true,
		},
		{
			"test umlauts",
			"öäüßÖÄÜ",
			true,
		},
		{
			"test empty word",
			"",
			true,
		},
		{
			"test two words with space",
			"hallo du",
			false,
		},
		{
			"test word with special character",
			"hallo.du",
			false,
		},
		{
			"test word with number",
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

func TestWordInWordList(t *testing.T) {
	tests := []struct {
		testcase string
		word     string
		wordList []string
		expected bool
	}{
		{
			"test word is in word list",
			"hey",
			[]string{"you", "hey", "there"},
			true,
		},
		{
			"test word is not in word list",
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
