package ui

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/nsf/termbox-go"
)

const FG = 0
const BG = 1

const def = termbox.ColorDefault

//TODO implement Screens
type UI struct {
	Dimensions     *UIDimensions
	InputBox       []rune
	InputBoxColors [2]termbox.Attribute
	ScreenLines    []string
	ScreenColors   [2]termbox.Attribute
	w              *bufio.Writer
	Channel        string
}

type UIDimensions struct {
	Right       int
	Left        int
	Top         int
	Bottom      int
	InputBoxTop int
	InputBox    int
}

var DefaultUI *UI // Default UI

func init() {
	DefaultUI = NewUI()
}

func NewUI() *UI {
	ui := &UI{
		ScreenColors:   [2]termbox.Attribute{termbox.ColorDefault, termbox.ColorDefault},
		InputBoxColors: [2]termbox.Attribute{termbox.ColorDefault, termbox.ColorDefault},
	}

	return ui
}

// Close the UI
func (u *UI) Close() {
	termbox.Flush()
	termbox.Close()
}

// Draw all child Drawables. Use this method to re-draw the whole UI.
// Instead of using this method, you should use the appropriate methods
// on child Drawables. These will re-draw themselves automatically.
//
// It is appropriate to use this method to draw the entire UI on first load.
//
// NOTE: You must call termbox.Flush() after calling this method;
// you should prefer the Loop() method.
func (ui *UI) Draw() {
	x, y := termbox.Size()

	ui.Dimensions = &UIDimensions{
		Right:       x - 1,
		Left:        0,
		Top:         0,
		Bottom:      y - 1,
		InputBoxTop: y - 3,
		InputBox:    y - 2,
	}

	termbox.Clear(def, def)

	ui.drawTopLine()
	ui.DrawScreen()
	ui.DrawInputBox()
	ui.drawBottomLine()

	termbox.Flush()
}

func (ui *UI) drawBottomLine() {
	line := ui.Dimensions.Bottom
	termbox.SetCell(ui.Dimensions.Left, line, '\u2570', def, def)
	for i := 1; i < ui.Dimensions.Right; i++ {
		termbox.SetCell(i, line, '\u2500', def, def)
	}
	termbox.SetCell(ui.Dimensions.Right, line, '\u256F', def, def)
}

func (ui *UI) drawTopLine() {
	line := ui.Dimensions.Top
	termbox.SetCell(ui.Dimensions.Left, line, '\u256D', def, def)
	for i := 1; i < ui.Dimensions.Right; i++ {
		termbox.SetCell(i, line, '\u2500', def, def)
	}
	termbox.SetCell(ui.Dimensions.Right, line, '\u256E', def, def)
}

// Enter the UI loop for the default UI
func Loop() {
	DefaultUI.Loop()
}

func (ui *UI) drawHorizontalLine(y int) {
	for i := 1; i < ui.Dimensions.Right; i++ {
		termbox.SetCell(i, y, '\u2500', termbox.ColorDefault, termbox.ColorDefault)
	}
}

func (u *UI) handleOneEvent(ev termbox.Event) bool {
	if ev.Type == termbox.EventKey {
		switch ev.Key {
		case termbox.KeyEsc:
			return true
		case termbox.KeyCtrlC:
			return true
		case termbox.KeyCtrlD:
			return true
		case termbox.KeyBackspace2:
			u.Backspace()
		case termbox.KeyEnter:
			defer u.Blank()

			s := u.GetInputString()

			switch {
			case s[0] != '/' && u.Channel == "":
				u.AppendScreenLine("WARNING: you can't talk here!")
			case u.Channel == "" || s[0] == '/':
				u.AppendScreenLine("DEBUG: got here")
				if s[0] == '/' {
					s = s[1:]
					u.AppendScreenLine(s)
				}

				cmdargs := strings.Split(fmt.Sprintf("%s", s), string([]byte{'\x00'}))

				u.AppendScreenLine(cmdargs[0])

				if cmdargs[0] == "join" {
					u.Channel = cmdargs[1]
				}

				u.AppendScreenLine("joining channel %s", u.Channel)

				switch cmdargs[0] {
				case "join":
					u.w.WriteString("JOIN " + u.Channel + "\r\n")
				default:
					u.AppendScreenLine("unknown command sequence: %v", cmdargs)
				}
			default:
				s := "PRIVMSG " + u.Channel + " :" + s + "\r\n"
				u.AppendScreenLine(strings.Trim(s, "\r\n"))
				u.w.WriteString(s)
			}

			u.w.Flush()
		default:
			switch ev.Ch {
			case 3:
				return true
			case 4:
				return true
			default:
				u.InputBox = append(u.InputBox, ev.Ch)
				u.DrawInputBox()
			}
		}
	}

	return false
}

func (u *UI) Loop() {
	termbox.Init()
	termbox.SetInputMode(termbox.InputEsc)

	defer func() {
		// recover()
		u.Close()
	}()

	var done bool

	u.Draw()

	u.AppendScreenLine("Welcome to Chatter v%d", 0)

	conn, err := net.Dial("tcp", "irc.freenode.net:6667")

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		r := bufio.NewReader(conn)
		for {
			rawline, _, err := r.ReadLine()
			line := string(rawline)

			lineparts := strings.SplitN(line, " ", 3)

			source := strings.Split(lineparts[0], "!")[0][1:]
			cmd := lineparts[1]
			args := lineparts[2]

			switch cmd {
			case "PRIVMSG":
				u.AppendScreenLine(source + ": " + strings.SplitN(args, ":", 2)[1])
			case "JOIN":
				u.AppendScreenLine("%s joined %s", source, args[1:])
			case "PART":
				u.AppendScreenLine("%s left %s", source, args[1:])
			default:
				if strings.Count(line, ":") > 1 {
					u.AppendScreenLine(strings.SplitN(line, ":", 3)[2])
				} else {
					u.AppendScreenLine(line)
				}
			}

			if err == io.EOF {
				return
			}
		}
	}()

	w := bufio.NewWriter(conn)

	w.WriteString("USER abc23234 0 * abc342342\r\n")
	w.WriteString("NICK abc2342342342\r\n")

	log.Println("logged in; joining", u.Channel)

	if u.Channel != "" {
		w.WriteString("JOIN " + u.Channel + "\r\n")
	}

	u.w = w

	w.Flush()

	for !done {
		done = u.handleOneEvent(termbox.PollEvent())
		termbox.Flush()
	}
}
