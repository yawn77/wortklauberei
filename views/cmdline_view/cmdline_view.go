package cmdlineview

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/yawn77/wortklauberei/handlers"
	"github.com/yawn77/wortklauberei/models"
	"github.com/yawn77/wortklauberei/utils"
)

var (
	// backgroundColor         = tcell.NewRGBColor(6, 0, 71)
	backgroundColor         = tcell.ColorBlack
	textColor               = tcell.NewRGBColor(255, 255, 255)
	activeBackgroundColor   = tcell.NewRGBColor(255, 95, 158)
	disabledBackgroundColor = tcell.NewRGBColor(179, 0, 94)
	enabledBackgroundColor  = tcell.NewRGBColor(233, 0, 100)
	solutionGray            = tcell.NewRGBColor(133, 133, 133)
	solutionGreen           = tcell.NewRGBColor(0, 150, 0)
	solutionYellow          = tcell.NewRGBColor(0, 150, 150)
)

type CmdlineView struct {
	gameHandler   handlers.GameHandler
	app           *tview.Application
	mainView      *tview.Flex
	fields        [][]*tview.InputField
	height        int
	width         int
	activeRow     int
	activeCol     int
	label         *tview.Flex
	footer        *tview.Flex
	gameOver      bool
	newGameHeight int
	newGameWidth  int
	version       string
}

func NewCmdlineView(gh handlers.GameHandler, version string) (clv CmdlineView) {
	clv.gameHandler = gh
	clv.version = version
	// TODO handle error
	s, _ := tcell.NewScreen()
	// TODO enable and fix tests
	// s.SetCursorStyle(tcell.CursorStyleSteadyUnderline)
	clv.app = tview.NewApplication()
	clv.app.SetScreen(s)
	return clv
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
		if !clv.gameOver {
			clv.activateNextField(true)
		}
	case tcell.KeyBackspace2:
		if !clv.gameOver {
			clv.activatePreviousField()
		}
	case tcell.KeyEnter:
		if !clv.gameOver {
			clv.checkSolution()
		}
	}
	return event
}

func (clv *CmdlineView) checkSolution() {
	if clv.gameOver {
		return
	}
	s := ""
	fields := clv.fields[clv.activeRow]
	for i := 0; i < len(fields); i++ {
		s += fields[i].GetText()
	}
	correct, gameOver, colorCode, _, valid := clv.gameHandler.CheckSolution(s)
	clv.gameOver = gameOver

	if valid != nil {
		clv.SetLabelText(valid.Error())
		return
	}

	clv.applyColorCode(colorCode)

	if correct || gameOver {
		if correct {
			clv.SetLabelText("CONGRATULATIONS! YOU WON :)")
		} else if gameOver {
			clv.SetLabelText("GAME OVER. YOU LOST :(")
		}
		return
	}

	clv.activateNextRow()
}

func (clv *CmdlineView) activeField() *tview.InputField {
	return clv.fields[clv.activeRow][clv.activeCol]
}

func setColor(field *tview.InputField, color tcell.Color) {
	field.SetBackgroundColor(color)
	field.SetFieldBackgroundColor(color)
}

func (clv *CmdlineView) applyColorCode(code models.ColorCode) {
	row := clv.fields[clv.activeRow]
	for i := 0; i < clv.width; i++ {
		field := row[i]
		switch code[i] {
		case models.Gray:
			setColor(field, solutionGray)
		case models.Green:
			setColor(field, solutionGreen)
		case models.Yellow:
			setColor(field, solutionYellow)
		}
	}
}

/*
 * create board
 */
func (clv *CmdlineView) NewGame(wordLength int, maxAttempts int) {
	clv.gameOver = false
	clv.newGameHeight = maxAttempts
	clv.newGameWidth = wordLength
	clv.height = maxAttempts
	clv.width = wordLength
	clv.mainView = clv.buildMainView(wordLength, maxAttempts, clv.version)
	clv.app.SetInputCapture(clv.inputCapture)
	clv.app.SetRoot(clv.mainView, true)
	clv.enableRow(0)
	clv.activateField(0, 0)
}

func (clv *CmdlineView) buildMainView(width int, height int, version string) (hflex *tview.Flex) {
	hflex = tview.NewFlex().SetDirection(tview.FlexRow)
	hflex.SetBorder(true).SetTitle(fmt.Sprintf(" wortklauberei %s ", version))
	hflex.SetTitleColor(textColor)
	hflex.SetBorderColor(activeBackgroundColor)
	hflex.SetBackgroundColor(backgroundColor)
	hflex.AddItem(createBox(), 0, 1, false)

	clv.fields = make([][]*tview.InputField, height)
	for i := range clv.fields {
		clv.fields[i] = make([]*tview.InputField, width)
	}
	hflex.AddItem(createBox(), 1, 0, false)
	for i := 0; i < height; i++ {
		row := tview.NewFlex().SetDirection(tview.FlexColumn)
		row.AddItem(createBox(), 0, 1, false)
		row.AddItem(createBox(), 1, 0, false)
		for j := 0; j < width; j++ {
			clv.fields[i][j] = clv.createInputField()
			row.AddItem(clv.fields[i][j], 1, 0, false)
			row.AddItem(createBox(), 1, 0, false)
		}
		row.AddItem(createBox(), 0, 1, false)
		hflex.AddItem(row, 1, 0, false)
		hflex.AddItem(createBox(), 1, 0, false)
	}

	hflex.AddItem(createBox(), 0, 1, false)
	clv.label = tview.NewFlex().SetDirection(tview.FlexColumn)
	clv.SetLabelText("")
	hflex.AddItem(clv.label, 1, 0, false)
	clv.footer = clv.createFooter()
	hflex.AddItem(clv.footer, 1, 0, false)

	return hflex
}

func createBox() *tview.Box {
	return tview.NewBox().SetBackgroundColor(backgroundColor)
}

func (clv *CmdlineView) SetLabelText(text string) {
	clv.label.Clear()
	clv.label.AddItem(createBox(), 0, 1, false)
	label := tview.NewTextView().SetText(text).SetTextColor(textColor)
	label.SetBackgroundColor(backgroundColor)
	clv.label.AddItem(label, utf8.RuneCountInString(text), 0, false)
	clv.label.AddItem(createBox(), 0, 1, false)
}

func (clv *CmdlineView) createInputField() *tview.InputField {
	field := tview.NewInputField().
		SetFieldWidth(1).
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return len([]rune(textToCheck)) <= 1 && (unicode.IsLetter(lastChar))
		}).
		SetChangedFunc(func(text string) {
			clv.SetLabelText("")
			if text == "" {
				return
			}
			if text != "" && text != "ß" {
				if utils.IsLower(text) {
					clv.activeField().SetText(strings.ToUpper(text))
					clv.activateNextField(false)
				}
			}
		})
	field.SetFieldTextColor(textColor)
	setColor(field, disabledBackgroundColor)

	return field
}

func (clv *CmdlineView) createFooter() (ft *tview.Flex) {
	ft = tview.NewFlex()
	quitBtn := createButton("Quit (Ctrl+Q)", func() {
		clv.quit()
	})
	ngBtn := createButton("New Game (Ctrl+N)", func() {
		clv.newGame()
	})
	ft.AddItem(createBox(), 1, 0, false)
	ft.AddItem(ngBtn, 25, 0, false)
	ft.AddItem(createBox(), 1, 0, false)
	ft.AddItem(quitBtn, 25, 0, false)
	ft.AddItem(createBox(), 0, 1, false)
	return ft
}

func createButton(label string, handler func()) (btn *tview.Button) {
	btn = tview.NewButton(label).SetSelectedFunc(handler)
	btn.SetStyle(tcell.StyleDefault.Background(activeBackgroundColor))
	btn.SetLabelColor(textColor)
	return btn
}

/*
 * navigation
 */
func (clv *CmdlineView) activateField(row int, col int) {
	clv.activeRow = row
	clv.activeCol = col
	field := clv.fields[row][col]
	setColor(field, activeBackgroundColor)
	clv.app.SetFocus(field)
}

func (clv *CmdlineView) deactivateField(row int, col int) {
	field := clv.fields[row][col]
	setColor(field, enabledBackgroundColor)
}

func (clv *CmdlineView) activateNextField(tab bool) {
	if !tab && clv.activeCol >= clv.width-1 {
		return
	}
	clv.deactivateField(clv.activeRow, clv.activeCol)
	if tab && clv.activeCol == clv.width-1 {
		clv.activeCol = -1
	}
	clv.activateField(clv.activeRow, clv.activeCol+1)
}

func (clv *CmdlineView) activatePreviousField() {
	if clv.activeField().GetText() != "" || clv.activeCol <= 0 {
		return
	}
	clv.deactivateField(clv.activeRow, clv.activeCol)
	clv.activateField(clv.activeRow, clv.activeCol-1)
}

func (clv *CmdlineView) activateNextRow() {
	clv.enableRow(clv.activeRow + 1)
	clv.activateField(clv.activeRow, 0)
}

func (clv *CmdlineView) enableRow(row int) {
	if row < 0 || row >= clv.height {
		return
	}
	clv.activeRow = row
	for _, field := range clv.fields[row] {
		setColor(field, enabledBackgroundColor)
	}
}
