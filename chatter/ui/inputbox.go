package ui

import (
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

func (ui *UI) GetInputString() string {
	buf := make([]byte, 0)
	for _, r := range ui.InputBox {
		b := make([]byte, 1)
		utf8.EncodeRune(b, r)
		buf = append(buf, b...)
	}
	return string(buf)
}

func (ui *UI) drawInputBoxTop() {
	line := ui.Dimensions.InputBoxTop
	termbox.SetCell(ui.Dimensions.Left, line, '\u251C', def, def)
	ui.drawHorizontalLine(line)
	termbox.SetCell(ui.Dimensions.Right, line, '\u2524', def, def)
}

func (ui *UI) drawInputBox() {
	line := ui.Dimensions.InputBox

	termbox.SetCell(ui.Dimensions.Left, line, '\u2502', def, def)

	i := 1
	tidx := 0
	for i < ui.Dimensions.Right && tidx < len(ui.InputBox) {
		termbox.SetCell(i, line, ui.InputBox[tidx], def, def)
		i++
		tidx++
	}

	// fill with blanks
	for i < ui.Dimensions.Right {
		termbox.SetCell(i, line, ' ', def, def)
		i++
	}

	termbox.SetCell(ui.Dimensions.Right, line, '\u2502', def, def)
	termbox.SetCursor(len(ui.InputBox)+1, ui.Dimensions.InputBox)
}

func (ui *UI) DrawInputBox() {
	termbox.SetCursor(2, ui.Dimensions.InputBox)

	ui.drawInputBoxTop()
	ui.drawInputBox()

	termbox.Flush()
}

func (ui *UI) Backspace() {
	if len(ui.InputBox) == 0 {
		return
	}

	ui.InputBox = ui.InputBox[:len(ui.InputBox)-1]

	ui.DrawInputBox()
}

func (ui *UI) Blank() {
	ui.InputBox = make([]rune, 0)
	ui.DrawInputBox()
}
