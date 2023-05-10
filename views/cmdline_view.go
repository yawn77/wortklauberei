package views

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/yawn77/wortklauberei/handlers"
)

type CmdlineView struct {
	app                *tview.Application
	mv                 *tview.Flex
	gh                 handlers.GameHandler
	ngWordLength       int
	ngNumberOfAttempts int
}

func NewCmdlineView(gh handlers.GameHandler, version string) (clv CmdlineView) {
	clv.app = tview.NewApplication()
	clv.mv = clv.buildMainView(version)
	clv.app.SetRoot(clv.mv, true)
	clv.app.SetInputCapture(clv.inputCapture)
	clv.app.EnableMouse(true)
	clv.gh = gh
	clv.ngWordLength = 5
	clv.ngNumberOfAttempts = 6
	return clv
}

func (gc CmdlineView) Run() {
	if err := gc.app.Run(); err != nil {
		panic(err)
	}
}

func (clv CmdlineView) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlN:
		clv.newGame()
	case tcell.KeyCtrlQ:
		clv.quit()
	}
	return event
}

func (clv *CmdlineView) buildMainView(version string) (hflex *tview.Flex) {
	hflex = tview.NewFlex().SetDirection(tview.FlexRow)
	hflex.SetBorder(true).SetTitle(fmt.Sprintf(" wortklauberei %s ", version))
	hflex.AddItem(tview.NewBox(), 0, 1, false)
	ft := clv.createFooter()
	hflex.AddItem(ft, 1, 0, false)
	return hflex
}

func (clv *CmdlineView) BuildNewBoard(wordLength int, maxAttempts int) {
	clv.mv.Clear()
}

func (clv *CmdlineView) createFooter() (ft *tview.Flex) {
	ft = tview.NewFlex()
	quitBtn := createButton("Quit (Ctrl+Q)", func() {
		clv.quit()
	})
	ngBtn := createButton("New Game (Ctrl+N)", func() {
		clv.newGame()
	})
	ft.AddItem(tview.NewBox(), 1, 0, false)
	ft.AddItem(ngBtn, 25, 0, false)
	ft.AddItem(tview.NewBox(), 1, 0, false)
	ft.AddItem(quitBtn, 25, 0, false)
	ft.AddItem(tview.NewBox(), 0, 1, false)
	return ft
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

func (clv *CmdlineView) cancelDialog() {
	clv.app.SetRoot(clv.mv, true)
}

func createButton(label string, handler func()) (btn *tview.Button) {
	btn = tview.NewButton(label).SetSelectedFunc(handler)
	btn.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlueViolet))
	btn.SetLabelColor(tcell.ColorGhostWhite)
	return btn
}

func verifyNumberInput(textToCheck string, lastChar rune) bool {
	i, err := strconv.Atoi(textToCheck)
	return err == nil && i >= 2 && i <= 9
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
			strconv.Itoa(6),
			1,
			verifyNumberInput,
			setWordLength).
		AddInputField("Number of Attempts (2-9)",
			strconv.Itoa(5),
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
	err := clv.gh.CreateNewGame(clv.ngWordLength, clv.ngNumberOfAttempts)
	if err != nil {
		// TODO
	}
}
