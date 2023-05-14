package utils_test

import (
	"testing"

	"github.com/yawn77/wortklauberei/utils"
)

func TestIsLettersOnly(t *testing.T) {
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
			b := utils.IsLettersOnly(tt.s)
			if b != tt.expected {
				t.Errorf("IsLettersOnly(%s) = %t", tt.s, b)
			}
		})
	}
}

func TestIsWordInWordList(t *testing.T) {
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
			b := utils.IsWordInWordList(tt.word, tt.wordList)
			if b != tt.expected {
				t.Errorf("WordIsInWordList(%s, %v) = %t", tt.word, tt.wordList, b)
			}
		})
	}
}

func TestIsLower(t *testing.T) {
	tests := []struct {
		testcase string
		s        string
		expected bool
	}{
		{
			"common word",
			"hey",
			true,
		},
		{
			"umlauts",
			"möp",
			true,
		},
		{
			"capital letter",
			"Hello",
			false,
		},
		{
			"space and special character",
			"hey you!",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			b := utils.IsLower(tt.s)
			if b != tt.expected {
				t.Errorf("IsLower(%s) = %t", tt.s, b)
			}
		})
	}
}
