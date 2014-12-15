package ui

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

func (ui *UI) AppendScreenLine(format string, a ...interface{}) {
	ui.ScreenLines = append(ui.ScreenLines, fmt.Sprintf(format, a...))
	ui.DrawScreen()
}

func (ui *UI) DrawScreen() {
	// --------------------- draw top border
	termbox.SetCell(ui.Dimensions.Left, 0, '\u256D', def, def)
	for i := 1; i < ui.Dimensions.Right; i++ {
		termbox.SetCell(i, 0, '\u2500', def, def)
	}
	termbox.SetCell(ui.Dimensions.Right, 0, '\u256E', def, def)

	lines := len(ui.ScreenLines) - 1

	for line := ui.Dimensions.InputBoxTop - 1; line > 0; line-- {
		termbox.SetCell(ui.Dimensions.Left, line, '\u2502', def, def)
		if lines > -1 {
			txt := ui.ScreenLines[lines]
			i := 1
			for _, c := range txt {
				termbox.SetCell(i, line, c, def, def)
				i++
			}
			for ; i < ui.Dimensions.Right; i++ {
				termbox.SetCell(i, line, ' ', def, def)
			}
			lines--
		}
		termbox.SetCell(ui.Dimensions.Right, line, '\u2502', def, def)
	}

	if len(ui.ScreenLines) > ui.Dimensions.Bottom {
		ui.ScreenLines = ui.ScreenLines[len(ui.ScreenLines)-ui.Dimensions.Bottom:]
	}

	termbox.Flush()
}
