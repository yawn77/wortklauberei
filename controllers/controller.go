package controllers

import (
	"fmt"
	"strings"

	"github.com/yawn77/wortklauberei/models"
	"github.com/yawn77/wortklauberei/utils"
	cmdlineview "github.com/yawn77/wortklauberei/views/cmdline_view"
)

type GameController struct {
	gameModel models.GameModel
	view      cmdlineview.CmdlineView
	version   string
}

func NewGameController(test bool, version string) (gc GameController, err error) {
	gc.version = version
	gc.view = cmdlineview.NewCmdlineView(&gc, test, gc.version)
	err = gc.CreateNewGame(4, 3)
	return gc, err
}

func (gc GameController) Run() {
	gc.view.Run()
}

func (gc *GameController) CreateNewGame(wordLength int, maxAttempts int) error {
	// TODO provide propper inputs
	gm, err := models.NewGameModel("JUPP", []string{"JUPP", "ATHI", "ZAHL", "COKE", "JACK", "DAUP"}, maxAttempts)
	if err != nil {
		return err
	}
	gc.gameModel = gm
	// TODO provide propper word lenth
	gc.view.NewGame(4, maxAttempts)
	return err
}

func (gc *GameController) CheckSolution(solution string) (correct bool, gameOver bool, colorCode models.ColorCode, keyboardColors models.KeyboardColors, valid error) {
	gm := &gc.gameModel
	if gm.GameOver {
		return false, true, nil, nil, fmt.Errorf("game is already over")
	}

	wl := len([]rune(solution))
	if wl != gm.WordLength {
		return false, false, nil, nil, fmt.Errorf("length of solution must be %d but is %d", gm.WordLength, wl)
	}
	if !utils.IsLettersOnly(solution) {
		return false, false, nil, nil, fmt.Errorf("solution must consists of letters  only")
	}
	if !utils.IsWordInWordList(solution, gm.ValidWords) {
		return false, false, nil, nil, fmt.Errorf("suggested solution is not in list of valid words")
	}

	correct = solution == gm.Word
	if correct || gm.CurAttempt >= gm.MaxAttempts-1 {
		gm.GameOver = true
	}

	colorCode = calculateColorCodeAndUpdateKeyboardColors(*gm, solution)
	gm.CurAttempt++
	return correct, gm.GameOver, colorCode, gm.KeyboardColors, nil
}

func calculateColorCodeAndUpdateKeyboardColors(gm models.GameModel, suggestion string) (colorCode models.ColorCode) {
	rsug := []rune(suggestion)
	lsug := len(rsug)
	rsol := []rune(gm.Word)
	if lsug != len(rsol) {
		return nil
	}
	rh := cacluateRuneHistogram(gm.Word)

	// mark greens
	for i := 0; i < lsug; i++ {
		r := rsug[i]
		if r == rsol[i] {
			colorCode = append(colorCode, models.Green)
			updateKeyboardColor(gm, r, models.Green)
			val, ok := rh[r]
			if ok {
				rh[r] = val - 1
			}
		} else {
			// may be overridden by yellow in next loop
			colorCode = append(colorCode, models.Gray)
			updateKeyboardColor(gm, r, models.Gray)
		}
	}

	// mark yellows and grays
	for i := 0; i < lsug; i++ {
		if colorCode[i] == models.Green {
			continue
		}
		r := rsug[i]
		val, ok := rh[r]
		sInSol := strings.Contains(gm.Word, string(r))
		if sInSol && ok && val > 0 {
			rh[r] = val - 1
			colorCode[i] = models.Yellow
			updateKeyboardColor(gm, r, models.Yellow)
		}
	}

	return colorCode
}

func updateKeyboardColor(gm models.GameModel, r rune, color models.Color) {
	val, ok := gm.KeyboardColors[r]
	if !ok || val < color {
		gm.KeyboardColors[r] = color
	}
}

func cacluateRuneHistogram(s string) (rh models.RuneHistogram) {
	rh = make(models.RuneHistogram)
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
