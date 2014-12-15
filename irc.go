//Package irc provides an IRC client
package irc

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type Client struct {
	sync.WaitGroup
	Server   string
	Port     int32
	Nick     string
	User     string
	Fullname string
	conn     *net.Conn
	Inbound  chan Message
	Outbound chan Message
}

type Message struct {
	Source  string
	Command string
	Args    []string
}

func (c *Client) Join(channel string) *Channel {
	return NewChannel(c)
}

func (c *Client) Quit() {
	c.send("QUIT :Leaving")
}

func (c *Client) send(cmd string, args ...string) {
	w := bufio.NewWriter(*c.conn)
	w.WriteString(":" + c.User + " " + cmd + " " + strings.Join(args, " ") + "\n")
	w.Flush()
}

func ParseMessage(s string) *Message {
	parts := strings.SplitN(s, " ", 3)

	sp := strings.SplitN(parts[2], ":", 2)

	params := strings.Split(strings.Trim(sp[0], " "), " ")
	fmt.Println(len(params))
	trailing := sp[1]

	m := &Message{
		Source:  parts[0][1:],
		Command: parts[1],
		Args:    append(params, trailing),
	}

	return m
}

func (c *Client) handleIncoming() {
	r := bufio.NewReader(*c.conn)

	for {
		rawline, _, err := r.ReadLine()
		line := string(rawline)

		if err != nil {
			if err == io.EOF {
				fmt.Println(line)
				log.Fatal("connection closed")
			}
			break
		}

		fmt.Println(line)
	}

	c.Done()
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Server, c.Port))

	if err != nil {
		return err
	}

	defer conn.Close()

	c.conn = &conn

	c.Add(1)

	go c.handleIncoming()

	c.send("USER", c.User, "0", "*", c.User)
	c.send("NICK", c.Nick)

	r := bufio.NewReader(os.Stdin)

	var cmd string
	var args []string

	for {
		fmt.Print("> ")

		rawline, _, _ := r.ReadLine()

		args = strings.Split(string(rawline), " ")
		cmd = args[0]
		args = args[1:]

		c.send(strings.ToUpper(cmd), args...)
		fmt.Println(cmd, args)
	}

	c.Wait()

	return nil
}

func NewClient(server, nick, user string, initialChannels ...string) *Client {
	c := &Client{
		Server: server,
		Nick:   nick,
		User:   user,
	}

	// for i, channel := range initialChannels {
	// 	c.Join(channel)
	// }

	return c
}
