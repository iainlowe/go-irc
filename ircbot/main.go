package main

import "github.com/ilowe/irc"

func main() {
	c := irc.NewClient("irc.freenode.net:6667", "abc123", "bob123")
	c.Connect()
	defer c.Quit()
	c.Join("#kjhkjhasd")
}
