package utils

import "unicode"

func IsLettersOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsWordInWordList(w string, wl []string) bool {
	for _, cw := range wl {
		if cw == w {
			return true
		}
	}
	return false
}
