package cmdlineview

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
)

func (clv *CmdlineView) cancelDialog() {
	clv.app.SetRoot(clv.mainView, true)
	clv.app.SetInputCapture(clv.inputCapture)
}

func (clv *CmdlineView) quit() {
	qd := tview.NewModal().
		SetText(" Do you want to quit wortklauberei? ").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				clv.app.Stop()
			} else {
				clv.cancelDialog()
			}
		})
	qd.SetBackgroundColor(backgroundColor)
	qd.SetTextColor(textColor)
	qd.SetButtonBackgroundColor(activeBackgroundColor)
	qd.SetButtonTextColor(textColor)
	clv.app.SetRoot(qd, true)
	clv.app.SetInputCapture(nil)
}

func (clv CmdlineView) newGame() {
	setWordLength := func(text string) {
		i, nil := strconv.Atoi(text)
		if nil != nil {
			i = 5
		}
		clv.newGameWidth = i
	}
	setNumberOfAttempts := func(text string) {
		i, nil := strconv.Atoi(text)
		if nil != nil {
			i = 6
		}
		clv.newGameHeight = i
	}
	form := tview.NewForm().
		AddInputField("Word Length (2-9)",
			strconv.Itoa(clv.newGameWidth),
			1,
			verifyNumberInput,
			setWordLength).
		AddInputField("Number of Attempts (2-9)",
			strconv.Itoa(clv.newGameHeight),
			1,
			verifyNumberInput,
			setNumberOfAttempts).
		AddButton("New Game", clv.setupGame).
		AddButton("Cancel", clv.cancelDialog)
	form.SetBorder(true).SetTitle(" New Game ").SetTitleAlign(tview.AlignCenter)
	form.SetBackgroundColor(backgroundColor)
	form.SetBorderColor(activeBackgroundColor)
	form.SetTitleColor(textColor)
	form.SetLabelColor(textColor)
	form.SetFieldTextColor(textColor)
	form.SetFieldBackgroundColor(activeBackgroundColor)
	form.SetButtonBackgroundColor(activeBackgroundColor)
	form.SetButtonTextColor(textColor)

	modal := func(p tview.Primitive, width, height int) *tview.Flex {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}(form, 30, 9)
	clv.app.SetRoot(modal, true)
	clv.app.SetInputCapture(nil)
}

func (clv *CmdlineView) setupGame() {
	err := clv.gameHandler.CreateNewGame(clv.newGameWidth, clv.newGameHeight)
	if err != nil {
		// TODO
		fmt.Println("ahoi")
	}
}

func verifyNumberInput(textToCheck string, lastChar rune) bool {
	i, err := strconv.Atoi(textToCheck)
	return err == nil && i >= 2 && i <= 9
}
