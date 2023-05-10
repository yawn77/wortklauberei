package models_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/yawn77/wortklauberei/models"
)

func TestNewGameModelSuccess(t *testing.T) {
	// arrange
	var (
		word           = "hellö"
		wordLength     = 5
		validWords     = []string{"hellö", "funny"}
		gameOver       = false
		curAttempt     = 0
		maxAttempts    = 6
		keyboardColors = models.KeyboardColors{}
	)

	// act
	m, err := models.NewGameModel(word, validWords, maxAttempts)

	// assert
	if err != nil {
		t.Error(err)
	}
	if m.Word != word {
		t.Errorf("model initialization failed: word %s != %s", m.Word, word)
	}
	if m.WordLength != wordLength {
		t.Errorf("model initialization failed: word length %d != %d", m.WordLength, wordLength)
	}
	if !reflect.DeepEqual(m.ValidWords, validWords) {
		t.Errorf("model initialization failed: valid words %v != %v", m.ValidWords, validWords)
	}
	if m.GameOver != gameOver {
		t.Errorf("model initialization failed: game over %t != %t", m.GameOver, gameOver)
	}
	if m.CurAttempt != curAttempt {
		t.Errorf("model initialization failed: current attempt %d != %d", m.CurAttempt, curAttempt)
	}
	if m.MaxAttempts != maxAttempts {
		t.Errorf("model initialization failed: max attempts %d != %d", m.MaxAttempts, maxAttempts)
	}
	if !reflect.DeepEqual(m.KeyboardColors, keyboardColors) {
		t.Errorf("model initialization failed: keyboard colors %v != %v", m.KeyboardColors, keyboardColors)
	}
}

func TestNewGameModelFail(t *testing.T) {
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
			_, err := models.NewGameModel(tt.word, tt.validWords, tt.maxAttempts)
			if err == nil || err.Error() != tt.expected.Error() {
				t.Errorf("error expected: %v != %v", tt.expected, err)
			}
		})
	}
}
