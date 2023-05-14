package views

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/yawn77/wortklauberei/handlers"
)

type CmdlineView struct {
	gameHandler        handlers.GameHandler
	app                *tview.Application
	mainView           *tview.Flex
	inputFields        [][]*tview.InputField
	curRow             int
	curCol             int
	errorlabel         *tview.TextView
	ngWordLength       int
	ngNumberOfAttempts int
	version            string
}

func NewCmdlineView(gh handlers.GameHandler, version string) (clv CmdlineView) {
	clv.gameHandler = gh
	clv.version = version
	clv.app = tview.NewApplication()
	clv.app.EnableMouse(true)
	return clv
}

func (clv *CmdlineView) CreateNewGameBoard(wordLength int, maxAttempts int) {
	clv.ngWordLength = wordLength
	clv.ngNumberOfAttempts = maxAttempts
	clv.mainView = clv.buildMainView(wordLength, maxAttempts, clv.version)
	clv.app.SetInputCapture(clv.inputCapture)
	clv.app.SetRoot(clv.mainView, true)
	clv.resetError()
	clv.curRow = 0
	clv.curCol = 0
	clv.app.SetFocus(clv.inputFields[clv.curRow][clv.curCol])
}

func (clv *CmdlineView) buildMainView(width int, height int, version string) (hflex *tview.Flex) {
	hflex = tview.NewFlex().SetDirection(tview.FlexRow)
	hflex.SetBorder(true).SetTitle(fmt.Sprintf(" wortklauberei %s ", version))
	hflex.AddItem(tview.NewBox(), 0, 1, false)

	clv.inputFields = make([][]*tview.InputField, height)
	for i := range clv.inputFields {
		clv.inputFields[i] = make([]*tview.InputField, width)
	}
	hflex.AddItem(tview.NewBox(), 1, 0, false)
	for i := 0; i < height; i++ {
		row := tview.NewFlex().SetDirection(tview.FlexColumn)
		row.AddItem(tview.NewBox(), 0, 1, false)
		row.AddItem(tview.NewBox(), 1, 0, false)
		for j := 0; j < width; j++ {
			clv.inputFields[i][j] = clv.newInputField()
			row.AddItem(clv.inputFields[i][j], 1, 0, false)
			row.AddItem(tview.NewBox(), 1, 0, false)
		}
		row.AddItem(tview.NewBox(), 0, 1, false)
		hflex.AddItem(row, 1, 0, false)
		hflex.AddItem(tview.NewBox(), 1, 0, false)
	}

	hflex.AddItem(tview.NewBox(), 0, 1, false)
	clv.errorlabel = tview.NewTextView().SetText("").SetTextColor(tcell.ColorRed)
	hflex.AddItem(clv.errorlabel, 1, 0, false)
	hflex.AddItem(clv.createFooter(), 1, 0, false)

	return hflex
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func (clv *CmdlineView) newInputField() *tview.InputField {
	inp := tview.NewInputField().
		SetFieldWidth(1).
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return len([]rune(textToCheck)) <= 1 && (unicode.IsLetter(lastChar))
		}).
		SetChangedFunc(func(text string) {
			if text != "" && text != "ß" && IsLower(text) {
				clv.inputFields[clv.curCol][clv.curRow].SetText(strings.ToUpper(text))
			}
		})
	inp.SetBackgroundColor(tcell.ColorBlue)
	inp.SetDisabled(true)

	return inp
}

func (gc CmdlineView) Run() {
	if err := gc.app.Run(); err != nil {
		panic(err)
	}
}

func (clv *CmdlineView) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlN:
		clv.newGame()
	case tcell.KeyCtrlQ:
		clv.quit()
	}
	return event
}

func verifyNumberInput(textToCheck string, lastChar rune) bool {
	i, err := strconv.Atoi(textToCheck)
	return err == nil && i >= 2 && i <= 9
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

func createButton(label string, handler func()) (btn *tview.Button) {
	btn = tview.NewButton(label).SetSelectedFunc(handler)
	btn.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlueViolet))
	btn.SetLabelColor(tcell.ColorGhostWhite)
	return btn
}

func (clv *CmdlineView) resetError() {
	clv.errorlabel.SetText("")
}

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
