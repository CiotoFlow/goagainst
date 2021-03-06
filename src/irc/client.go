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
	"strings"
)

const (
	RPL_TOPIC = "332";
	RPL_NAMREPLY = "353"
)

type IRC struct {
	config ServerConfig
	Nickname string
	Mode string

	Channels []Channel
	Users []User
	
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
		irc.Send("USER %s 0 * :%s", nick, irc.config.Realname)
	}

	c := irc.StartMachine()
	defer irc.StopMachine(c)
	
	reg()
	for {
		msg := <- c
		if msg.Command == "443" {
			nick = nick+"_"
			reg()
		} else if msg.Command == "MODE" && msg.Param(0) == nick {
			server, ok := msg.Entity.(*Server)
			if ok && server.Name == nick {
				irc.Nickname = nick
				irc.Mode = msg.Trailing
				fmt.Println("Registered as", irc.Nickname, "with mode", irc.Mode)
				break
			}
		}
	}
}

func (irc *IRC) autoJoin() {
	c := irc.StartMachine()
	defer irc.StopMachine(c)

	for {
		msg := <- c
		if msg.Command == "001" {
			// welcome
			for _, name := range irc.config.AutoJoin {
				irc.Send("JOIN %s", name)
			}
			break
		}
	}
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

func (irc *IRC) GetChannel(name string) *Channel {
	for _, c := range(irc.Channels) {
		if c.Name == name {
			return &c
		}
	}
	c := Channel{name, []*User{}}
	irc.Channels = append(irc.Channels, c)
	return &c
}

func (irc *IRC) GetUser(nick string) *User {
	for _, u := range(irc.Users) {
		if u.Nick == nick {
			return &u
		}
	}
	u := User{nick, "", "", false}
	irc.Users = append(irc.Users, u)
	return &u
}

func (irc *IRC) handleState(msg *Message) {
	if (msg.Command == RPL_NAMREPLY && IsServer(msg.Entity) && msg.Param(0) == irc.Nickname) {
		name := msg.Param(2)
		channel := irc.GetChannel(name)
		nicks := strings.Split(msg.Trailing, " ")
		channel.Users = make([]*User, len(nicks))
		for i, nick := range(nicks) {
			user := irc.GetUser(nick)
			user.Valid = true
			channel.Users[i] = user
		}
		fmt.Println ("Channel names:", channel.Users)
	}
}

func (irc *IRC) Loop() error {

	if (!irc.conn) {
		return errors.New("not connected")
	}

	go irc.registerNick()
	if irc.config.AutoPong {
		go irc.autoPong()
	}
	go irc.autoJoin()
	f := func(msg *Message) { irc.handleState(msg) }
	irc.NotifyCallback (&f)

	for {
		line, err := irc.r.ReadLine()
		if err != nil {
			return err
		}

		msg, err := ParseMessage(line)
		if (err != nil) {
			fmt.Println(err)
			continue
		}

		irc.notify(msg)
	}

	return nil
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
