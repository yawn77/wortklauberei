package controllers

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/yawn77/wortklauberei/models"
)

func TestCheckSolutionCorrectAfterThreeAttempts(t *testing.T) {
	// arrange
	word := "dizzy"
	gc, _ := NewGameController("0.1.0")
	gm, _ := models.NewGameModel(word, []string{"hello", "whizz", "pozzy", "dizzy"}, 6)
	gc.gameModel = gm

	// act & assert
	correct, gameOver, colorCode, keyboardColors, valid := gc.CheckSolution("hello")
	if correct ||
		gameOver ||
		!reflect.DeepEqual(colorCode, models.ColorCode{models.Gray, models.Gray, models.Gray, models.Gray, models.Gray}) ||
		!reflect.DeepEqual(keyboardColors, models.KeyboardColors{104: models.Gray, 101: models.Gray, 108: models.Gray, 111: models.Gray}) ||
		valid != nil {
		t.Errorf("error after hello")
	}
	correct, gameOver, colorCode, keyboardColors, valid = gc.CheckSolution("whizz")
	if correct ||
		gameOver ||
		!reflect.DeepEqual(colorCode, models.ColorCode{models.Gray, models.Gray, models.Yellow, models.Green, models.Yellow}) ||
		!reflect.DeepEqual(keyboardColors, models.KeyboardColors{104: models.Gray, 101: models.Gray, 108: models.Gray, 111: models.Gray, 119: models.Gray, 105: models.Yellow, 122: models.Green}) ||
		valid != nil {
		t.Errorf("error after whizz")
	}
	correct, gameOver, colorCode, keyboardColors, valid = gc.CheckSolution("dizzy")
	if !correct ||
		!gameOver ||
		!reflect.DeepEqual(colorCode, models.ColorCode{models.Green, models.Green, models.Green, models.Green, models.Green}) ||
		!reflect.DeepEqual(keyboardColors, models.KeyboardColors{104: models.Gray, 101: models.Gray, 108: models.Gray, 111: models.Gray, 119: models.Gray, 105: models.Green, 122: models.Green, 100: models.Green, 121: models.Green}) ||
		valid != nil {
		t.Errorf("error after dizzy")
	}
}

func TestCheckSolutionCorrectAtLastAttempt(t *testing.T) {
	// arrange
	word := "dizzy"
	gc, _ := NewGameController("0.1.0")
	gm, _ := models.NewGameModel(word, []string{"hello", "whizz", "pozzy", "dizzy"}, 3)
	gc.gameModel = gm

	// act & assert
	correct, gameOver, colorCode, keyboardColors, valid := gc.CheckSolution("hello")
	if correct ||
		gameOver ||
		!reflect.DeepEqual(colorCode, models.ColorCode{models.Gray, models.Gray, models.Gray, models.Gray, models.Gray}) ||
		!reflect.DeepEqual(keyboardColors, models.KeyboardColors{104: models.Gray, 101: models.Gray, 108: models.Gray, 111: models.Gray}) ||
		valid != nil {
		t.Errorf("error after hello")
	}
	correct, gameOver, colorCode, keyboardColors, valid = gc.CheckSolution("whizz")
	if correct ||
		gameOver ||
		!reflect.DeepEqual(colorCode, models.ColorCode{models.Gray, models.Gray, models.Yellow, models.Green, models.Yellow}) ||
		!reflect.DeepEqual(keyboardColors, models.KeyboardColors{104: models.Gray, 101: models.Gray, 108: models.Gray, 111: models.Gray, 119: models.Gray, 105: models.Yellow, 122: models.Green}) ||
		valid != nil {
		t.Errorf("error after whizz")
	}
	correct, gameOver, colorCode, keyboardColors, valid = gc.CheckSolution("dizzy")
	if !correct ||
		!gameOver ||
		!reflect.DeepEqual(colorCode, models.ColorCode{models.Green, models.Green, models.Green, models.Green, models.Green}) ||
		!reflect.DeepEqual(keyboardColors, models.KeyboardColors{104: models.Gray, 101: models.Gray, 108: models.Gray, 111: models.Gray, 119: models.Gray, 105: models.Green, 122: models.Green, 100: models.Green, 121: models.Green}) ||
		valid != nil {
		t.Errorf("error after dizzy")
	}
}

func TestCheckSolutionAfterGameOver(t *testing.T) {
	// arrange
	word := "dizzy"
	gc, _ := NewGameController("0.1.0")
	gm, _ := models.NewGameModel(word, []string{"hello", "whizz", "pozzy", "dizzy"}, 3)
	gc.gameModel = gm

	// act & assert
	correct, gameOver, colorCode, keyboardColors, valid := gc.CheckSolution("dizzy")
	if !correct ||
		!gameOver ||
		!reflect.DeepEqual(colorCode, models.ColorCode{models.Green, models.Green, models.Green, models.Green, models.Green}) ||
		!reflect.DeepEqual(keyboardColors, models.KeyboardColors{100: models.Green, 105: models.Green, 122: models.Green, 121: models.Green}) ||
		valid != nil {
		t.Errorf("error after dizzy")
	}
	correct, gameOver, colorCode, keyboardColors, valid = gc.CheckSolution("whizz")
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
			gc, _ := NewGameController("0.1.0")
			gm, _ := models.NewGameModel(tt.word, tt.validWords, 6)
			gc.gameModel = gm

			// act
			correct, gameOver, colorCode, keyboardColors, valid := gc.CheckSolution(tt.solution)

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
		expectedCC models.ColorCode
		expectedKC models.KeyboardColors
	}{
		{
			"test suggestion == solution",
			"hello",
			"hello",
			models.ColorCode{models.Green, models.Green, models.Green, models.Green, models.Green},
			models.KeyboardColors{104: models.Green, 101: models.Green, 108: models.Green, 111: models.Green},
		},
		{
			"test len(suggestion) != len(solution)",
			"hello",
			"hell",
			nil,
			models.KeyboardColors{},
		},
		{
			"test some models.Green, some models.Yellow, some models.Gray",
			"halle",
			"hello",
			models.ColorCode{models.Green, models.Gray, models.Green, models.Green, models.Yellow},
			models.KeyboardColors{104: models.Green, 97: models.Gray, 108: models.Green, 101: models.Yellow},
		},
		{
			"test one letter models.Green, models.Yellow and models.Gray",
			"laffllk",
			"laflenr",
			models.ColorCode{models.Green, models.Green, models.Green, models.Gray, models.Yellow, models.Gray, models.Gray},
			models.KeyboardColors{108: models.Green, 97: models.Green, 102: models.Green, 107: models.Gray},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			gm, _ := models.NewGameModel(tt.solution, []string{tt.solution, tt.suggestion}, 6)
			cc := calculateColorCodeAndUpdateKeyboardColors(gm, tt.suggestion)
			if !reflect.DeepEqual(cc, tt.expectedCC) {
				t.Errorf("calculateColorCodeAndUpdateKeyboardColors(%s, %s) = %v, expected %v", tt.suggestion, tt.solution, cc, tt.expectedCC)
			}
			if !reflect.DeepEqual(gm.KeyboardColors, tt.expectedKC) {
				t.Errorf("unexpected keyboard colors: %v != %v", gm.KeyboardColors, tt.expectedKC)
			}
		})
	}
}

func TestCalculateRuneHistogram(t *testing.T) {
	tests := []struct {
		testcase string
		s        string
		expected models.RuneHistogram
	}{
		{
			"test casual string",
			"hello",
			models.RuneHistogram{
				104: 1,
				101: 1,
				108: 2,
				111: 1,
			},
		},
		{
			"test umlauts",
			"ööÖßü",
			models.RuneHistogram{
				246: 2,
				214: 1,
				223: 1,
				252: 1,
			},
		},
		{
			"test empty string",
			"",
			models.RuneHistogram{},
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
