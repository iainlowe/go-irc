package irc

import (
	"fmt"
	"testing"
)

func TestParseMessage1(t *testing.T) {
	m := ParseMessage(":source PRIVMSG target :This is the message")

	switch {
	case m.Source != "source":
		fmt.Println("source not parsed properly:", m.Source)
		t.Fail()
	case m.Command != "PRIVMSG":
		fmt.Println("command not parsed properly:", m.Args)
		t.Fail()
	case m.Args[0] != "target" || m.Args[1] != "This is the message":
		fmt.Printf("args not parsed properly: (target '%v') (trailing: '%v') %v\n", m.Args[0], m.Args[1], m.Args)
		t.Fail()
	}
}

func TestParseMessage2(t *testing.T) {
	m := ParseMessage(":source PRIVMSG target abc :This is the message")

	switch {
	case m.Source != "source":
		fmt.Println("source not parsed properly:", m.Source)
		t.Fail()
	case m.Command != "PRIVMSG":
		fmt.Println("command not parsed properly:", m.Args)
		t.Fail()
	case m.Args[0] != "target" || m.Args[2] != "This is the message":
		fmt.Printf("args not parsed properly: (target '%v') (trailing: '%v') %v\n", m.Args[0], m.Args[2], m.Args)
		t.Fail()
	}
}
