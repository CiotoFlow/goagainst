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
	config ServerConfig
	
	conn bool
	sock net.Conn
	r *textproto.Reader
	w *textproto.Writer
	chanSubscribers [](chan *Message)
	callbackSubscribers [](IrcCallback)
}

type IrcCallback func(msg *Message);

type Message struct {
	Type string
	From *string
	To *string
	Content *string
}

func NewIRC(config ServerConfig) *IRC {
	return &IRC { config: config }
}

func (irc *IRC) NotifyChan(c chan *Message) {
	irc.chanSubscribers = append(irc.chanSubscribers, c)
}

func (irc *IRC) NotifyCallback(cb IrcCallback) {
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
		c(cmd)
	}
}

func (irc *IRC) Loop() error {
	var cmd string

	if (!irc.conn) {
		return errors.New("not connected")
	}

	irc.Send("NICK %s", irc.config.Nickname)
	irc.Send("USER %s 0 * :Stocazzo", irc.config.Nickname)

	for {
		line, err := irc.r.ReadLine()
		if err != nil {
			return err
		}
		
		fmt.Println(line)

		resp := strings.Split(line, " ")

		var from *string
		var to *string
		var content *string
		
		/* FIXME: */
		if (resp[0][0] == ':' && len(resp) > 1) {
			cmd = resp[1]
			from = &resp[0]
			to = &resp[2]
			if len(resp) >= 4 {
				content = &resp[3]
			} else {
				content = nil
			}
		} else {
			cmd = resp[0]
			from = nil
			to = nil
			content = &resp[1]
		}

		if (cmd == "PING") {
			if irc.config.AutoPing {
				irc.Send("PONG %s", *content)
			}
		} else if (cmd == "001") {
			for _, name := range irc.config.AutoJoin {
				irc.Send("JOIN %s", name)
			}
		} else if (cmd == "443") {
			/* Duplicated NICK */
		} else {
			cmd := Message { cmd, from, to, content }
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
