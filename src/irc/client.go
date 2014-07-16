/*
 * RFCs:
 *   https://tools.ietf.org/html/rfc2810
 *   https://tools.ietf.org/html/rfc2811
 *   https://tools.ietf.org/html/rfc2812
 *   https://tools.ietf.org/html/rfc2813
 *
 * https://www.alien.net.au/irc/irc2numerics.html
 */
package irc

import (
	"bufio"
	"net/textproto"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
)

const (
	IRC_EVENT_NICK = 0
	IRC_EVENT_PRIV_MSG = 1
	/* TODO: define more... */
)

type IRC struct {
	config ServerConfig

	conn bool
	sock net.Conn
	r *textproto.Reader
	w *textproto.Writer
	chanSubscribers [](chan *Message)
	callbackSubscribers [](*func(*Message))
}

func NewIRC(config ServerConfig) *IRC {
	return &IRC { config: config }
}

func (irc *IRC) RemoveCallback(c *func(*Message)) {
	// FIXME: lock
	for i, cc := range irc.callbackSubscribers {
		if cc == c {
			irc.callbackSubscribers[i] = irc.callbackSubscribers[len(irc.callbackSubscribers)-1]
			irc.callbackSubscribers = irc.callbackSubscribers[0:len(irc.callbackSubscribers)-1]
			break
		}
	}
}

func (irc *IRC) NotifyCallback(cb *func(*Message)) {
	irc.callbackSubscribers = append(irc.callbackSubscribers, cb)
}

func (irc *IRC) Send(format string, a ...interface{}) error {
	return irc.w.PrintfLine(format, a...)
}

func (irc *IRC) notify(cmd *Message) {
	for _, c := range irc.chanSubscribers {
		c <- cmd
	}
	for _, c := range irc.callbackSubscribers {
		(*c)(cmd)
	}
}

func (irc *IRC) StartMachine() chan *Message {
	c := make(chan *Message, 1000)
	// FIXME: lock
	irc.chanSubscribers = append(irc.chanSubscribers, c)
	return c
}

func (irc *IRC) StopMachine(c chan *Message) {
	// FIXME: lock
	for i, cc := range irc.chanSubscribers {
		if cc == c {
			irc.chanSubscribers[i] = irc.chanSubscribers[len(irc.chanSubscribers)-1]
			irc.chanSubscribers = irc.chanSubscribers[0:len(irc.chanSubscribers)-1]
			break
		}
	}
}

func (irc *IRC) registerNick() {
	// register nick until it's not duplicated
	nick := irc.config.Nickname
	reg := func() {
		irc.Send("NICK %s", nick)
		irc.Send("USER %s 0 * :Stocazzo", nick)
	}

	c := irc.StartMachine()
	reg()
	for {
		msg := <- c
		if msg.Command == "443" {
			nick = nick+"_"
			reg()
		} else if msg.Command == "MODE" && msg.Prefix == nick && msg.Params[0] == nick {
			fmt.Println("Registered as", nick)
			break
		}
	}
	irc.StopMachine(c)
}

func (irc *IRC) autoPong() {
	// register nick until it's not duplicated
	f := func(msg *Message) {
		if msg.Command == "PING" {
			irc.Send("PONG %s", msg.Trailing)
		}
	}
	
	irc.NotifyCallback (&f)
}

func (irc *IRC) Loop() error {

	if (!irc.conn) {
		return errors.New("not connected")
	}

	go irc.registerNick()
	if irc.config.AutoPong {
		go irc.autoPong()
	}
	

	for {
		line, err := irc.r.ReadLine()
		if err != nil {
			return err
		}

		fmt.Println(line)

		msg, err := ParseMessage(line)
		if (err != nil) {
			fmt.Println(err)
			continue
		}

		irc.handleMessage(msg)
		irc.notify(msg)
	}

	return nil
}

func (irc *IRC) handleMessage(msg *Message) {
	switch msg.Command {
	case "001":
		/* Welcome */
		for _, name := range irc.config.AutoJoin {
			irc.Send("JOIN %s", name)
		}
	}
}

func (irc *IRC) Disconnect() {
	irc.Send("QUIT :against")
	irc.sock.Close()
}

func (irc *IRC) Connect() error {
	var err error

	if (irc.conn) {
		return errors.New("already connected")
	}

	if (irc.config.UseTls) {
		irc.sock, err = tls.Dial("tcp", irc.config.Address, &tls.Config{InsecureSkipVerify: true})
	} else {
		irc.sock, err = net.Dial("tcp", irc.config.Address)
	}

	if (err == nil) {
		irc.conn = true
		irc.r = textproto.NewReader(bufio.NewReader(irc.sock))
		irc.w = textproto.NewWriter(bufio.NewWriter(irc.sock))
	}

	return err
}
