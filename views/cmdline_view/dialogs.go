package cmdlineview

import (
	"fmt"
	"strconv"

	"github.com/rivo/tview"
)

func (clv *CmdlineView) cancelDialog() {
	clv.app.SetRoot(clv.mainView, true)
}

func (clv *CmdlineView) quit() {
	qd := tview.NewModal().
		SetText("Do you want to quit wortklauberei?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				clv.app.Stop()
			} else {
				clv.cancelDialog()
			}
		})
	clv.app.SetRoot(qd, true)
}

func (clv CmdlineView) newGame() {
	setWordLength := func(text string) {
		i, nil := strconv.Atoi(text)
		if nil != nil {
			i = 5
		}
		clv.ngWordLength = i
	}
	setNumberOfAttempts := func(text string) {
		i, nil := strconv.Atoi(text)
		if nil != nil {
			i = 6
		}
		clv.ngNumberOfAttempts = i
	}
	form := tview.NewForm().
		AddInputField("Word Length (2-9)",
			strconv.Itoa(clv.ngWordLength),
			1,
			verifyNumberInput,
			setWordLength).
		AddInputField("Number of Attempts (2-9)",
			strconv.Itoa(clv.ngNumberOfAttempts),
			1,
			verifyNumberInput,
			setNumberOfAttempts).
		AddButton("New Game", clv.setupGame).
		AddButton("Cancel", clv.cancelDialog)
	form.SetBorder(true).SetTitle("New Game").SetTitleAlign(tview.AlignCenter)

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
}

func (clv *CmdlineView) setupGame() {
	err := clv.gameHandler.CreateNewGame(clv.ngWordLength, clv.ngNumberOfAttempts)
	if err != nil {
		// TODO
		fmt.Println("ahoi")
	}
}
