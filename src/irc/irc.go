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
	IRC_EVENT_NICK = 0
	IRC_EVENT_PRIV_MSG = 1
	/* TODO: define more... */
)

type IRC struct {
	conn bool
	nick string
	server string
	channel string
	useTls bool
	sock net.Conn
	r *textproto.Reader
	w *textproto.Writer
	chanSubscribers [](chan *IrcMsg)
	callbackSubscribers [](IrcCallback)
}

type IrcCallback func(msg *IrcMsg);

type IrcMsg struct {
	Type string
	From string
	To *string
	Content string
}

func NewIRC(nick, server, channel string, useTls bool) *IRC {
	return &IRC {
		nick: nick,
		server: server,
		channel: channel,
		useTls: useTls,
	}
}

func (irc *IRC) NotifyChan(c chan *IrcMsg) {
	irc.chanSubscribers = append(irc.chanSubscribers, c)
}

func (irc *IRC) NotifyCallback(cb IrcCallback) {
	irc.callbackSubscribers = append(irc.callbackSubscribers, cb)
}

func (irc *IRC) Send(format string, a ...interface{}) error {
	return irc.w.PrintfLine(format, a...)
}

func (irc *IRC) notify(cmd *IrcMsg) {
	for _, c := range irc.chanSubscribers {
		c <- cmd
	}
	for _, c := range irc.callbackSubscribers {
		c(cmd)
	}
}

func (irc *IRC) Loop() error {
	var cmd string

	if (!irc.conn) {
		return errors.New("not connected")
	}

	irc.Send("NICK %s", irc.nick)
	irc.Send("USER %s 0 * :Stocazzo", irc.nick)

	for {
		line, err := irc.r.ReadLine()
		if err != nil {
			return err
		}
		
		fmt.Println(line)

		resp := strings.Split(line, " ")

		/* XXX: handle better */
		if (resp[0][0] == ':' && len(resp) > 1) {
			cmd = resp[1]
		} else {
			cmd = resp[0]
		}

		if (cmd == "PING") {
			irc.Send("PONG %s", resp[1])
		} else if (cmd == "001") {
			irc.Send("JOIN %s", irc.channel)
		} else if (cmd == "443") {
			/* Duplicated NICK */
		} else {
			cmd := IrcMsg { cmd, "", nil, "" }
			irc.notify(&cmd)
		}

		/* TODO: add more commands */
	}

	return nil
}

func (irc *IRC) Disconnect() {
	irc.Send("QUIT against")
	irc.sock.Close()
}

func (irc *IRC) Connect() error {
	var err error

	if (irc.conn) {
		return errors.New("already connected")
	}

	if (irc.useTls) {
		irc.sock, err = tls.Dial("tcp", irc.server, &tls.Config{InsecureSkipVerify: true})
	} else {
		irc.sock, err = net.Dial("tcp", irc.server)
	}

	if (err == nil) {
		irc.conn = true
		irc.r = textproto.NewReader(bufio.NewReader(irc.sock))
		irc.w = textproto.NewWriter(bufio.NewWriter(irc.sock))		
	}

	return err
}
