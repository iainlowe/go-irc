package ui

import (
	"time"

	"github.com/nsf/termbox-go"
)

func pause(length int64) {
	select {
	case <-time.After(time.Second * 5):
		break
	}
}

func println(y int, msg string) {
	cd := termbox.ColorDefault

	for i, c := range msg {
		termbox.SetCell(i, y, c, cd, cd)
	}
}

func drawLines(start int, lines []string) {
	lineidx := len(lines) - 1
	idx := start
	for idx > -1 && lineidx > -1 {
		println(idx, lines[lineidx])
		idx--
		lineidx--
	}
}
