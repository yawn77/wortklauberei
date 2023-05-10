package handlers

import "github.com/yawn77/wortklauberei/models"

type GameHandler interface {
	CreateNewGame(wordLength int, maxAttempts int) error
	CheckSolution(solution string) (correct bool, gameOver bool, colorCode models.ColorCode, keyboardColors models.KeyboardColors, valid error)
}
