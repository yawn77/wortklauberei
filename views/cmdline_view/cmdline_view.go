package cmdlineview

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/yawn77/wortklauberei/handlers"
	"github.com/yawn77/wortklauberei/utils"
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
	clv.enableRow(clv.curRow)
	clv.setFocus()
}

func (clv *CmdlineView) setFocus() {
	clv.app.SetFocus(clv.inputFields[clv.curRow][clv.curCol])
}

func (clv *CmdlineView) curInputField() *tview.InputField {
	return clv.inputFields[clv.curRow][clv.curCol]
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

func (clv *CmdlineView) newInputField() *tview.InputField {
	inp := tview.NewInputField().
		SetFieldWidth(1).
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return len([]rune(textToCheck)) <= 1 && (unicode.IsLetter(lastChar))
		}).
		SetChangedFunc(func(text string) {
			if text != "" && text != "ß" && utils.IsLower(text) {
				clv.curInputField().SetText(strings.ToUpper(text))
				clv.goToNextInputField(false)
			}
		})
	inp.SetBackgroundColor(tcell.ColorBlue)
	inp.SetDisabled(true)

	return inp
}

func (clv *CmdlineView) goToNextInputField(overflow bool) {
	if overflow && clv.curCol == len(clv.inputFields[clv.curRow])-1 {
		clv.curCol = -1
	}
	if clv.curCol >= len(clv.inputFields[clv.curRow])-1 {
		return
	}
	clv.curCol++
	clv.setFocus()
}

func (clv *CmdlineView) goToPreviousInputField() {
	curText := clv.curInputField().GetText()
	if clv.curCol == 0 || curText != "" {
		return
	}
	clv.curCol--
	clv.setFocus()
}

func (clv *CmdlineView) enableRow(row int) {
	if row < 0 || row >= len(clv.inputFields) {
		return
	}
	if row > 0 {
		prev := row - 1
		for col := 0; col < len(clv.inputFields[prev]); col++ {
			clv.inputFields[prev][col].SetDisabled(true)
		}
	}
	for col := 0; col < len(clv.inputFields[row]); col++ {
		clv.inputFields[row][col].SetDisabled(false)
	}
	clv.curCol = 0
	clv.setFocus()
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
	case tcell.KeyTab:
		clv.goToNextInputField(true)
	case tcell.KeyBackspace2:
		clv.goToPreviousInputField()
	case tcell.KeyEnter:
		clv.checkSolution()
	}
	return event
}

func (clv *CmdlineView) checkSolution() {
	s := ""
	fields := clv.inputFields[clv.curRow]
	for i := 0; i < len(fields); i++ {
		s += fields[i].GetText()
	}
	correct, gameOver, _, _, valid := clv.gameHandler.CheckSolution(s)
	if valid != nil {
		clv.errorlabel.SetText(valid.Error())
	} else if correct {
		clv.errorlabel.SetText("CONGRATULATIONS! YOU WON :)")
	} else if gameOver {
		clv.errorlabel.SetText("GAME OVER. YOU LOST :(")
	} else {
		clv.curRow++
		clv.enableRow(clv.curRow)
	}
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
