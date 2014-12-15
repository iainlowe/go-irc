package main

import (
	"flag"
	"log"
	"os"

	"github.com/ilowe/irc/chatter/ui"
)

func main() {
	ui.DefaultUI.Channel = *flag.String("c", "#linux", "Channel to join on connect")
	flag.Parse()

	log.SetOutput(os.Stderr)

	ui.Loop()
}
