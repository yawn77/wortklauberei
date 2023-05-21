package cmdlineview

import (
	"testing"
)

func TestVerifyNumberInput(t *testing.T) {
	tests := []struct {
		testcase    string
		textToCheck string
		lastChar    rune
		expected    bool
	}{
		{
			"test digit",
			"6",
			[]rune("6")[0],
			true,
		},
		{
			"test empty string",
			"",
			[]rune("a")[0],
			false,
		},
		{
			"test t > 9",
			"10",
			[]rune("0")[0],
			false,
		},
		{
			"test t < 2",
			"1",
			[]rune("1")[0],
			false,
		},
		{
			"test char",
			"a",
			[]rune("a")[0],
			false,
		},
		{
			"test special char",
			".",
			[]rune(".")[0],
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			b := verifyNumberInput(tt.textToCheck, tt.lastChar)
			if b != tt.expected {
				t.Errorf("%s: %t != %t", tt.testcase, b, tt.expected)
			}
		})
	}
}
